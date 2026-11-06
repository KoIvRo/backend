package domain

import "testing"

func TestSignupValidate(t *testing.T) {
	tests := []struct {
		name      string
		signup    Signup
		wantField string
	}{
		{
			name:      "missing client",
			signup:    Signup{Phone: "+79990001122"},
			wantField: "client",
		},
		{
			name:      "missing phone",
			signup:    Signup{Client: "Иванов"},
			wantField: "phone",
		},
		{
			name:   "valid signup",
			signup: Signup{Client: "Иванов", Phone: "+79990001122"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.signup.Validate()

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("Validate() error = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("Validate() error = nil, want field %q", tt.wantField)
			}
			if err.Field != tt.wantField {
				t.Errorf("Validate() field = %q, want %q", err.Field, tt.wantField)
			}
		})
	}
}
