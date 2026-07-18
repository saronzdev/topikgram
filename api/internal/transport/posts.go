package transport

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"topikgram/api/internal/domain"
	"topikgram/api/internal/service"
	"topikgram/api/internal/store"
	"topikgram/api/internal/validator"
)

type PostHandler struct {
	mux             *http.ServeMux
	postStore       *store.PostStore
	userStore       *store.UserStore
	interestService *service.InterestService
}

func NewPostHandler(mux *http.ServeMux, db *store.Pool) *PostHandler {
	ps := store.NewPostStore(db)
	return &PostHandler{
		mux:             mux,
		postStore:       ps,
		userStore:       store.NewUserStore(db),
		interestService: service.NewInterestService(ps, store.NewInterestStore(db)),
	}
}

func (h *PostHandler) RegisterRoutes(prefix string) {
	h.mux.Handle("GET /"+prefix+"/posts", optionalAuthMiddleware(http.HandlerFunc(h.List)))
	h.mux.Handle("GET /"+prefix+"/posts/{id}", optionalAuthMiddleware(http.HandlerFunc(h.GetByID)))
	h.mux.Handle("POST /"+prefix+"/posts", authMiddleware(http.HandlerFunc(h.Create)))
	h.mux.Handle("PUT /"+prefix+"/posts/{id}", authMiddleware(http.HandlerFunc(h.Update)))
	h.mux.Handle("POST /"+prefix+"/posts/{id}/like", authMiddleware(http.HandlerFunc(h.Like)))
	h.mux.Handle("POST /"+prefix+"/posts/{id}/save", authMiddleware(http.HandlerFunc(h.Save)))
	h.mux.Handle("DELETE /"+prefix+"/posts/{id}", authMiddleware(http.HandlerFunc(h.Delete)))
	h.mux.Handle("DELETE /"+prefix+"/posts/{id}/like", authMiddleware(http.HandlerFunc(h.Unlike)))
	h.mux.Handle("DELETE /"+prefix+"/posts/{id}/save", authMiddleware(http.HandlerFunc(h.Unsave)))
	h.mux.Handle("GET /"+prefix+"/posts/{id}/likes", authMiddleware(http.HandlerFunc(h.GetLikes)))
	h.mux.Handle("GET /"+prefix+"/posts/{id}/saves", authMiddleware(http.HandlerFunc(h.GetSaves)))
}

func (h *PostHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input domain.CreatePostInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		badRequest(w, "invalid body")
		return
	}
	if input.Body == "" {
		badRequest(w, "body is required")
		return
	}
	if !validator.MaxLength(input.Body, 5000) {
		badRequest(w, "body must be at most 5000 characters")
		return
	}
	if len(input.Topics) == 0 || len(input.Topics) > 3 {
		badRequest(w, "between 1 and 3 topics are required")
		return
	}
	for _, t := range input.Topics {
		if t < 0 || t > 20 {
			badRequest(w, "invalid topic id")
			return
		}
	}
	userID := UserIDFromCtx(r.Context())
	post, err := h.postStore.Create(r.Context(), userID, &input)
	if err != nil {
		log.Printf("[ERROR] Create post: userID=%d body=%q topics=%v err=%v", userID, input.Body[:min(50, len(input.Body))], input.Topics, err)
		internalError(w, mapStoreError(err))
		return
	}
	jsonResp(w, http.StatusCreated, post)
}

func (h *PostHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		badRequest(w, "invalid id")
		return
	}
	var input domain.UpdatePostInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		badRequest(w, "invalid body")
		return
	}
	if input.Body == "" {
		badRequest(w, "body is required")
		return
	}
	if !validator.MaxLength(input.Body, 5000) {
		badRequest(w, "body must be at most 5000 characters")
		return
	}
	userID := UserIDFromCtx(r.Context())
	post, err := h.postStore.Update(r.Context(), id, userID, &input)
	if err == domain.ErrNotFound {
		notFound(w, "post not found or not yours")
		return
	}
	if err != nil {
		internalError(w, mapStoreError(err))
		return
	}
	jsonResp(w, http.StatusOK, post)
}

func (h *PostHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := UserIDFromCtx(r.Context())
	cursor := r.URL.Query().Get("cursor")
	limit := 50
	if l, err := strconv.Atoi(r.URL.Query().Get("limit")); err == nil && l > 0 && l <= 100 {
		limit = l
	}

	posts, _, hasMore, err := h.postStore.List(r.Context(), userID, cursor, limit+1)
	if err != nil {
		log.Printf("[ERROR] List posts: userID=%d cursor=%q limit=%d err=%v", userID, cursor, limit, err)
		internalError(w, mapStoreError(err))
		return
	}

	if hasMore && len(posts) > limit {
		posts = posts[:limit]
	}

	resp := domain.PostListResponse{
		Posts:   posts,
		HasMore: hasMore,
	}
	if hasMore && len(posts) > 0 {
		resp.NextCursor = posts[len(posts)-1].CreatedAt.Format("2006-01-02T15:04:05.000Z")
	}
	jsonResp(w, http.StatusOK, resp)
}

func (h *PostHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		badRequest(w, "invalid id")
		return
	}
	post, err := h.postStore.GetByID(r.Context(), id)
	if err == domain.ErrNotFound {
		notFound(w, "post not found")
		return
	}
	if err != nil {
		internalError(w, mapStoreError(err))
		return
	}
	jsonResp(w, http.StatusOK, post)
}

func (h *PostHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		badRequest(w, "invalid id")
		return
	}
	userID := UserIDFromCtx(r.Context())
	if err := h.postStore.Delete(r.Context(), id, userID); err == domain.ErrNotFound {
		notFound(w, "post not found or not yours")
		return
	} else if err != nil {
		internalError(w, mapStoreError(err))
		return
	}
	jsonResp(w, http.StatusOK, map[string]string{"message": "deleted"})
}

func (h *PostHandler) Like(w http.ResponseWriter, r *http.Request) {
	h.toggleLike(w, r, true)
}

func (h *PostHandler) Unlike(w http.ResponseWriter, r *http.Request) {
	h.toggleLike(w, r, false)
}

func (h *PostHandler) toggleLike(w http.ResponseWriter, r *http.Request, like bool) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		badRequest(w, "invalid id")
		return
	}
	userID := UserIDFromCtx(r.Context())
	if like {
		err = h.postStore.Like(r.Context(), userID, id)
	} else {
		err = h.postStore.Unlike(r.Context(), userID, id)
	}
	if err != nil {
		internalError(w, mapStoreError(err))
		return
	}

		delta := -0.1
		if like {
			delta = 0.1
		}
		h.interestService.UpdatePostInterest(r.Context(), userID, id, delta)

	jsonResp(w, http.StatusOK, map[string]bool{"liked": like})
}

func (h *PostHandler) Save(w http.ResponseWriter, r *http.Request) {
	h.toggleSave(w, r, true)
}

func (h *PostHandler) Unsave(w http.ResponseWriter, r *http.Request) {
	h.toggleSave(w, r, false)
}

func (h *PostHandler) toggleSave(w http.ResponseWriter, r *http.Request, save bool) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		badRequest(w, "invalid id")
		return
	}
	userID := UserIDFromCtx(r.Context())
	if save {
		err = h.postStore.Save(r.Context(), userID, id)
	} else {
		err = h.postStore.Unsave(r.Context(), userID, id)
	}
	if err != nil {
		internalError(w, mapStoreError(err))
		return
	}

		delta := -0.2
		if save {
			delta = 0.2
		}
		h.interestService.UpdatePostInterest(r.Context(), userID, id, delta)

	jsonResp(w, http.StatusOK, map[string]bool{"saved": save})
}

func (h *PostHandler) GetLikes(w http.ResponseWriter, r *http.Request) {
	postID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		badRequest(w, "invalid id")
		return
	}
	page, limit := parsePagination(r)
	users, total, err := h.postStore.GetLikes(r.Context(), postID, page, limit)
	if err != nil {
		internalError(w, mapStoreError(err))
		return
	}
	jsonResp(w, http.StatusOK, domain.PaginatedUsers{Users: users, Total: total, Page: page, Limit: limit})
}

func (h *PostHandler) GetSaves(w http.ResponseWriter, r *http.Request) {
	postID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		badRequest(w, "invalid id")
		return
	}
	page, limit := parsePagination(r)
	users, total, err := h.postStore.GetSaves(r.Context(), postID, page, limit)
	if err != nil {
		internalError(w, mapStoreError(err))
		return
	}
	jsonResp(w, http.StatusOK, domain.PaginatedUsers{Users: users, Total: total, Page: page, Limit: limit})
}

func parsePagination(r *http.Request) (int, int) {
	page := 1
	limit := 20
	if p, err := strconv.Atoi(r.URL.Query().Get("page")); err == nil && p > 0 {
		page = p
	}
	if l, err := strconv.Atoi(r.URL.Query().Get("limit")); err == nil && l > 0 && l <= 100 {
		limit = l
	}
	return page, limit
}
