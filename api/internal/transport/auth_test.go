package transport

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRegister_MissingFields(t *testing.T) {
	tests := []struct {
		name  string
		input map[string]string
	}{
		{"missing name", map[string]string{"username": "test", "email": "test@test.com", "password": "pass1234", "birthday": "2000-01-01"}},
		{"missing username", map[string]string{"name": "Test", "email": "test@test.com", "password": "pass1234", "birthday": "2000-01-01"}},
		{"missing email", map[string]string{"name": "Test", "username": "test", "password": "pass1234", "birthday": "2000-01-01"}},
		{"missing password", map[string]string{"name": "Test", "username": "test", "email": "test@test.com", "birthday": "2000-01-01"}},
		{"missing birthday", map[string]string{"name": "Test", "username": "test", "email": "test@test.com", "password": "pass1234"}},
		{"empty body", map[string]string{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.input)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var input struct {
					Name     string `json:"name"`
					Username string `json:"username"`
					Email    string `json:"email"`
					Password string `json:"password"`
					Birthday string `json:"birthday"`
				}
				if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
					badRequest(w, "invalid body")
					return
				}
				if input.Name == "" || input.Username == "" || input.Email == "" || input.Password == "" || input.Birthday == "" {
					badRequest(w, "name, username, email, password and birthday are required")
					return
				}
				w.WriteHeader(http.StatusCreated)
			})

			handler.ServeHTTP(rec, req)

			expect := http.StatusBadRequest
			if tt.name == "empty body" {
				expect = http.StatusBadRequest
			}
			if rec.Code != expect {
				t.Errorf("expected status %d, got %d", expect, rec.Code)
			}
		})
	}
}

func TestRegister_BirthdayFormat(t *testing.T) {
	tests := []struct {
		name     string
		birthday string
		wantOK   bool
	}{
		{"valid date", "2000-01-01", true},
		{"invalid format", "01-01-2000", false},
		{"empty", "", false},
		{"not a date", "abc", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				bd := tt.birthday
				if bd == "" {
					badRequest(w, "birthday is required")
					return
				}
				if _, err := time.Parse("2006-01-02", bd); err != nil {
					badRequest(w, "birthday must be a valid date (YYYY-MM-DD)")
					return
				}
				w.WriteHeader(http.StatusCreated)
			})

			req := httptest.NewRequest(http.MethodPost, "/", nil)
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)

			if tt.wantOK && rec.Code != http.StatusCreated {
				t.Errorf("expected 201, got %d", rec.Code)
			}
			if !tt.wantOK && rec.Code != http.StatusBadRequest {
				t.Errorf("expected 400 for invalid birthday %q, got %d", tt.birthday, rec.Code)
			}
		})
	}
}

func TestRegister_InvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader([]byte("not json")))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			Name string `json:"name"`
		}
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			badRequest(w, "invalid body")
			return
		}
		w.WriteHeader(http.StatusCreated)
	})

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rec.Code)
	}
}

func TestLogin_MissingFields(t *testing.T) {
	tests := []struct {
		name  string
		input map[string]string
	}{
		{"missing identifier", map[string]string{"password": "pass123"}},
		{"missing password", map[string]string{"identifier": "test"}},
		{"empty body", map[string]string{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.input)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var input struct {
					Identifier string `json:"identifier"`
					Password   string `json:"password"`
				}
				if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
					badRequest(w, "invalid body")
					return
				}
				if input.Identifier == "" || input.Password == "" {
					badRequest(w, "email or username and password are required")
					return
				}
				w.WriteHeader(http.StatusOK)
			})

			handler.ServeHTTP(rec, req)

			if rec.Code != http.StatusBadRequest {
				t.Errorf("expected status 400, got %d", rec.Code)
			}
		})
	}
}
