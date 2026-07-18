package transport

import (
	"encoding/json"
	"net/http"
	"strconv"
	"topikgram/api/internal/domain"
	"topikgram/api/internal/service"
	"topikgram/api/internal/store"
	"topikgram/api/internal/validator"
)

type CommentHandler struct {
	mux             *http.ServeMux
	commentStore    *store.CommentStore
	interestService *service.InterestService
	postStore       *store.PostStore
}

func NewCommentHandler(mux *http.ServeMux, db *store.Pool) *CommentHandler {
	ps := store.NewPostStore(db)
	return &CommentHandler{
		mux:             mux,
		commentStore:    store.NewCommentStore(db),
		interestService: service.NewInterestService(ps, store.NewInterestStore(db)),
		postStore:       ps,
	}
}

func (h *CommentHandler) RegisterRoutes(prefix string) {
	h.mux.Handle("GET /"+prefix+"/comments/{id}", optionalAuthMiddleware(http.HandlerFunc(h.GetByPostID)))
	h.mux.Handle("GET /"+prefix+"/comments/user/{id}", optionalAuthMiddleware(http.HandlerFunc(h.GetByUserID)))
	h.mux.Handle("POST /"+prefix+"/comments", authMiddleware(http.HandlerFunc(h.Create)))
}

func (h *CommentHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input domain.CreateCommentInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		badRequest(w, "invalid body")
		return
	}
	if input.Content == "" || input.PostID == 0 {
		badRequest(w, "content and post_id are required")
		return
	}
	if !validator.MaxLength(input.Content, 2000) {
		badRequest(w, "content must be at most 2000 characters")
		return
	}
	userID := UserIDFromCtx(r.Context())
	comment, err := h.commentStore.Create(r.Context(), userID, &input)
	if err != nil {
		internalError(w, mapStoreError(err))
		return
	}

	h.interestService.UpdatePostInterest(r.Context(), userID, input.PostID, 0.15)

	jsonResp(w, http.StatusCreated, comment)
}

func (h *CommentHandler) GetByPostID(w http.ResponseWriter, r *http.Request) {
	postID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		badRequest(w, "invalid post id")
		return
	}
	comments, err := h.commentStore.GetByPostID(r.Context(), postID)
	if err != nil {
		internalError(w, mapStoreError(err))
		return
	}
	if comments == nil {
		comments = []domain.Comment{}
	}
	jsonResp(w, http.StatusOK, comments)
}

func (h *CommentHandler) GetByUserID(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		badRequest(w, "invalid user id")
		return
	}
	comments, err := h.commentStore.GetByUserID(r.Context(), userID)
	if err != nil {
		internalError(w, mapStoreError(err))
		return
	}
	if comments == nil {
		comments = []domain.Comment{}
	}
	jsonResp(w, http.StatusOK, comments)
}
