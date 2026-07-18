package transport

import (
	"net/http"
	"strconv"
	"topikgram/api/internal/domain"
	"topikgram/api/internal/store"
	"topikgram/api/internal/validator"
)

type UserHandler struct {
	mux       *http.ServeMux
	userStore *store.UserStore
}

func NewUserHandler(mux *http.ServeMux, db *store.Pool) *UserHandler {
	return &UserHandler{mux: mux, userStore: store.NewUserStore(db)}
}

func (h *UserHandler) RegisterRoutes(prefix string) {
	h.mux.Handle("GET /"+prefix+"/users", optionalAuthMiddleware(http.HandlerFunc(h.GetAll)))
	h.mux.Handle("GET /"+prefix+"/users/{id}", optionalAuthMiddleware(http.HandlerFunc(h.GetByID)))
	h.mux.Handle("POST /"+prefix+"/users/{id}/follow", authMiddleware(http.HandlerFunc(h.Follow)))
	h.mux.Handle("DELETE /"+prefix+"/users/{id}/follow", authMiddleware(http.HandlerFunc(h.Unfollow)))
	h.mux.HandleFunc("GET /"+prefix+"/users/{id}/followers", h.GetFollowers)
	h.mux.HandleFunc("GET /"+prefix+"/users/{id}/following", h.GetFollowing)
}

func (h *UserHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username != "" {
		if !validator.Username(username) {
			badRequest(w, "invalid username format")
			return
		}
		user, err := h.userStore.GetByUsernamePublic(r.Context(), username)
		if err == domain.ErrNotFound {
			notFound(w, "user not found")
			return
		}
		if err != nil {
			internalError(w, mapStoreError(err))
			return
		}
		jsonResp(w, http.StatusOK, user)
		return
	}
	users, err := h.userStore.List(r.Context())
	if err != nil {
		internalError(w, mapStoreError(err))
		return
	}
	jsonResp(w, http.StatusOK, users)
}

func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		badRequest(w, "invalid id")
		return
	}
	user, err := h.userStore.GetByIDPublic(r.Context(), id)
	if err == domain.ErrNotFound {
		notFound(w, "user not found")
		return
	}
	if err != nil {
		internalError(w, mapStoreError(err))
		return
	}
	jsonResp(w, http.StatusOK, user)
}

func (h *UserHandler) Follow(w http.ResponseWriter, r *http.Request) {
	followeeID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		badRequest(w, "invalid id")
		return
	}
	me := UserIDFromCtx(r.Context())
	if me == followeeID {
		badRequest(w, "cannot follow yourself")
		return
	}
	if err := h.userStore.Follow(r.Context(), me, followeeID); err != nil {
		internalError(w, mapStoreError(err))
		return
	}
	jsonResp(w, http.StatusOK, map[string]string{"message": "followed"})
}

func (h *UserHandler) Unfollow(w http.ResponseWriter, r *http.Request) {
	followeeID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		badRequest(w, "invalid id")
		return
	}
	meID := UserIDFromCtx(r.Context())
	if err := h.userStore.Unfollow(r.Context(), meID, followeeID); err != nil {
		internalError(w, mapStoreError(err))
		return
	}
	jsonResp(w, http.StatusOK, map[string]string{"message": "unfollowed"})
}

func (h *UserHandler) GetFollowers(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		badRequest(w, "invalid id")
		return
	}
	users, err := h.userStore.Followers(r.Context(), id)
	if err != nil {
		internalError(w, mapStoreError(err))
		return
	}
	if users == nil {
		users = []domain.UserPublic{}
	}
	jsonResp(w, http.StatusOK, users)
}

func (h *UserHandler) GetFollowing(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		badRequest(w, "invalid id")
		return
	}
	users, err := h.userStore.Following(r.Context(), id)
	if err != nil {
		internalError(w, mapStoreError(err))
		return
	}
	if users == nil {
		users = []domain.UserPublic{}
	}
	jsonResp(w, http.StatusOK, users)
}
