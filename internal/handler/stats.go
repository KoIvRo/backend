package handler

import "net/http"

func (h *Handler) ClassStats(w http.ResponseWriter, r *http.Request) {
	statistics, err := h.stats.MostAttended(r.Context())
	if err != nil {
		writeInternalServerError(w)
		return
	}

	writeJSON(w, http.StatusOK, statistics)
}
