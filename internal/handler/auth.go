package handler

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/KoIvRo/backend/internal/repository"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type userKey struct{}
type attemptStore struct {
	sync.Mutex
	failures map[string][]time.Time
}
type Auth struct {
	Users    *repository.UserRepository
	Secret   []byte
	TTL      time.Duration
	attempts attemptStore
}

func NewAuth(users *repository.UserRepository, secret string, ttl time.Duration) *Auth {
	return &Auth{Users: users, Secret: []byte(secret), TTL: ttl, attempts: attemptStore{failures: make(map[string][]time.Time)}}
}
func CurrentUser(r *http.Request) repository.User {
	return r.Context().Value(userKey{}).(repository.User)
}
func (a *Auth) credentials(r *http.Request) (repository.User, error) {
	header := r.Header.Get("Authorization")
	if strings.HasPrefix(header, "Basic ") {
		username, password, ok := r.BasicAuth()
		if !ok {
			return repository.User{}, errors.New("invalid credentials")
		}
		u, err := a.Users.ByUsername(r.Context(), username)
		if err != nil || bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)) != nil {
			return repository.User{}, errors.New("invalid credentials")
		}
		return u, nil
	}
	if strings.HasPrefix(header, "Bearer ") {
		return a.verify(r.Context(), strings.TrimPrefix(header, "Bearer "))
	}
	return repository.User{}, errors.New("invalid credentials")
}
func (a *Auth) Required(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := a.attemptKey(r)
		if key != "" && a.limited(key) {
			writeError(w, 429, "too many login attempts", "")
			return
		}
		u, err := a.credentials(r)
		if err != nil {
			if key != "" {
				a.failed(key)
			}
			w.Header().Set("WWW-Authenticate", "Basic realm=fitness")
			writeError(w, 401, "unauthorized", "")
			return
		}
		if key != "" {
			a.clear(key)
		}
		next(w, r.WithContext(context.WithValue(r.Context(), userKey{}, u)))
	}
}
func (a *Auth) Owned(classes repository.ClassRepository, next http.HandlerFunc) http.HandlerFunc {
	return a.Required(func(w http.ResponseWriter, r *http.Request) {
		id, err := getIDFromRequest(r)
		if err != nil {
			writeError(w, 400, "invalid class ID", "id")
			return
		}
		c, err := classes.GetByID(r.Context(), id)
		if errors.Is(err, repository.ErrNotFound) {
			writeError(w, 404, "class not found", "")
			return
		}
		if err != nil {
			writeInternalServerError(w)
			return
		}
		u := CurrentUser(r)
		if u.Role != "admin" && c.OwnerID != u.ID {
			writeError(w, 403, "forbidden", "")
			return
		}
		next(w, r)
	})
}
func (a *Auth) Register(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if json.NewDecoder(r.Body).Decode(&body) != nil || strings.TrimSpace(body.Username) == "" || body.Password == "" {
		writeError(w, 400, "username and password required", "")
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		writeError(w, 400, "invalid password", "")
		return
	}
	u, err := a.Users.Create(r.Context(), strings.TrimSpace(body.Username), string(hash))
	if errors.Is(err, repository.ErrUserExists) {
		writeError(w, 409, "username already exists", "")
		return
	}
	if err != nil {
		writeInternalServerError(w)
		return
	}
	writeJSON(w, 201, u)
}
func (a *Auth) Me(w http.ResponseWriter, r *http.Request) { writeJSON(w, 200, CurrentUser(r)) }
func (a *Auth) Login(w http.ResponseWriter, r *http.Request) {
	key := a.attemptKey(r)
	if a.limited(key) {
		writeError(w, 429, "too many login attempts", "")
		return
	}
	var u repository.User
	var err error
	if r.Header.Get("Authorization") != "" {
		u, err = a.credentials(r)
	} else {
		var body struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		if json.NewDecoder(r.Body).Decode(&body) != nil {
			err = errors.New("invalid credentials")
		} else {
			key = r.RemoteAddr + ":" + body.Username
			u, err = a.Users.ByUsername(r.Context(), body.Username)
			if err == nil && bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(body.Password)) != nil {
				err = errors.New("invalid credentials")
			}
		}
	}
	if err != nil {
		a.failed(key)
		w.Header().Set("WWW-Authenticate", "Basic realm=fitness")
		writeError(w, 401, "unauthorized", "")
		return
	}
	a.clear(key)
	writeJSON(w, 200, map[string]string{"access_token": a.token(u)})
}
func (a *Auth) attemptKey(r *http.Request) string {
	username, _, ok := r.BasicAuth()
	if !ok {
		return ""
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		host = r.RemoteAddr
	}
	return host + ":" + username
}
func (a *Auth) limited(key string) bool {
	a.attempts.Lock()
	defer a.attempts.Unlock()
	now := time.Now()
	recent := a.attempts.failures[key][:0]
	for _, t := range a.attempts.failures[key] {
		if now.Sub(t) < time.Minute {
			recent = append(recent, t)
		}
	}
	a.attempts.failures[key] = recent
	return len(recent) >= 5
}
func (a *Auth) failed(key string) {
	a.attempts.Lock()
	defer a.attempts.Unlock()
	a.attempts.failures[key] = append(a.attempts.failures[key], time.Now())
}
func (a *Auth) clear(key string) {
	a.attempts.Lock()
	defer a.attempts.Unlock()
	delete(a.attempts.failures, key)
}

var b64 = base64.RawURLEncoding

func (a *Auth) token(u repository.User) string {
	head := b64.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`))
	payload, _ := json.Marshal(map[string]any{"sub": strconv.Itoa(u.ID), "exp": time.Now().Add(a.TTL).Unix()})
	body := head + "." + b64.EncodeToString(payload)
	mac := hmac.New(sha256.New, a.Secret)
	mac.Write([]byte(body))
	return body + "." + b64.EncodeToString(mac.Sum(nil))
}
func (a *Auth) verify(ctx context.Context, token string) (repository.User, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return repository.User{}, errors.New("invalid token")
	}
	header, err := b64.DecodeString(parts[0])
	if err != nil || string(header) != `{"alg":"HS256","typ":"JWT"}` {
		return repository.User{}, errors.New("invalid token")
	}
	sig, err := b64.DecodeString(parts[2])
	if err != nil {
		return repository.User{}, err
	}
	mac := hmac.New(sha256.New, a.Secret)
	mac.Write([]byte(parts[0] + "." + parts[1]))
	if !hmac.Equal(sig, mac.Sum(nil)) {
		return repository.User{}, errors.New("invalid token")
	}
	raw, err := b64.DecodeString(parts[1])
	if err != nil {
		return repository.User{}, err
	}
	var claims struct {
		Sub string `json:"sub"`
		Exp int64  `json:"exp"`
	}
	if json.Unmarshal(raw, &claims) != nil || claims.Exp <= time.Now().Unix() || claims.Sub == "" {
		return repository.User{}, errors.New("invalid token")
	}
	id, err := strconv.Atoi(claims.Sub)
	if err != nil {
		return repository.User{}, err
	}
	u, err := a.Users.ByID(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return repository.User{}, errors.New("invalid token")
	}
	return u, err
}
