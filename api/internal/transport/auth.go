package transport

import (
	"encoding/json"
	"net/http"
	"time"
	"topikgram/api/internal/domain"
	"topikgram/api/internal/store"
	"topikgram/api/internal/validator"
)

type AuthHandler struct {
	mux       *http.ServeMux
	authStore *store.AuthStore
	userStore *store.UserStore
}

func NewAuthHandler(mux *http.ServeMux, db *store.Pool) *AuthHandler {
	return &AuthHandler{mux: mux, authStore: store.NewAuthStore(db), userStore: store.NewUserStore(db)}
}

func (h *AuthHandler) RegisterRoutes(prefix string) {
	h.mux.HandleFunc("POST /"+prefix+"/auth/register", h.Register)
	h.mux.HandleFunc("POST /"+prefix+"/auth/login", h.Login)
	h.mux.Handle("GET /"+prefix+"/auth/me", authMiddleware(http.HandlerFunc(h.Me)))
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var input domain.RegisterInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		badRequest(w, "invalid body")
		return
	}
	if input.Name == "" || input.Username == "" || input.Email == "" || input.Password == "" || input.Birthday == "" {
		badRequest(w, "name, username, email, password and birthday are required")
		return
	}
	if !validator.Email(input.Email) {
		badRequest(w, "invalid email format")
		return
	}
	if !validator.Username(input.Username) {
		badRequest(w, "username must be 3-20 chars, alphanumeric and underscores only")
		return
	}
	if !validator.Password(input.Password) {
		badRequest(w, "password must be between 8 and 24 characters")
		return
	}
	if _, err := time.Parse("2006-01-02", input.Birthday); err != nil {
		badRequest(w, "birthday must be a valid date (YYYY-MM-DD)")
		return
	}

	user, err := h.authStore.Create(r.Context(), &input)
	if err != nil {
		msg := mapStoreError(err)
		if msg == "username already exists" || msg == "email already exists" {
			conflict(w, msg)
			return
		}
		internalError(w, msg)
		return
	}

	token, err := GenerateToken(user.ID)
	if err != nil {
		internalError(w, "internal server error")
		return
	}

	setCookie(w, token)
	jsonResp(w, http.StatusCreated, domain.AuthResponse{User: *user})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var input domain.LoginInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		badRequest(w, "invalid body")
		return
	}
	if input.Identifier == "" || input.Password == "" {
		badRequest(w, "email or username and password are required")
		return
	}

	user, err := h.authStore.Login(r.Context(), &input)
	if err == domain.ErrNotFound || err == domain.ErrIncorrectPassword {
		unauthorized(w)
		return
	}
	if err != nil {
		internalError(w, mapStoreError(err))
		return
	}

	token, err := GenerateToken(user.ID)
	if err != nil {
		internalError(w, "internal server error")
		return
	}
	setCookie(w, token)
	jsonResp(w, http.StatusOK, domain.AuthResponse{User: *user})
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID := UserIDFromCtx(r.Context())
	user, err := h.userStore.GetByID(r.Context(), userID)
	if err != nil {
		unauthorized(w)
		return
	}
	jsonResp(w, http.StatusOK, user)
}
