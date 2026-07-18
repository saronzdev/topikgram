package store

import (
	"context"
	"time"
	"topikgram/api/internal/domain"

	"github.com/jackc/pgx/v5"
)

type PostStore struct {
	db *Pool
}

func NewPostStore(db *Pool) *PostStore {
	return &PostStore{db: db}
}

const postCols = "p.id, p.body, p.topics_id, p.created_at, p.user_id, u.id, u.name, u.username"

func scanPost(row pgx.Row) (domain.Post, error) {
	p := domain.Post{}
	err := row.Scan(&p.ID, &p.Body, &p.TopicsID, &p.CreatedAt, &p.UserID, &p.User.ID, &p.User.Name, &p.User.Username, &p.Likes, &p.Saves, &p.Comments, &p.Liked, &p.Saved)
	return p, err
}

func (s *PostStore) Create(ctx context.Context, userID int, input *domain.CreatePostInput) (*domain.Post, error) {
	p := &domain.Post{}
	err := s.db.QueryRow(ctx,
		`INSERT INTO posts (body, user_id, topics_id) VALUES ($1, $2, $3)
		 RETURNING id, body, topics_id, created_at, user_id`, input.Body, userID, input.Topics,
	).Scan(&p.ID, &p.Body, &p.TopicsID, &p.CreatedAt, &p.UserID)
	if err != nil {
		return nil, err
	}
	p.User = domain.UserPublic{ID: userID}
	return p, nil
}

func (s *PostStore) GetByID(ctx context.Context, id int) (*domain.Post, error) {
	p, err := scanPost(s.db.QueryRow(ctx,
		`SELECT `+postCols+`,
		COALESCE((SELECT COUNT(*) FROM likes l WHERE l.post_id=p.id), 0) AS likes,
		COALESCE((SELECT COUNT(*) FROM saves s WHERE s.post_id=p.id), 0) AS saves,
		COALESCE((SELECT COUNT(*) FROM comments c WHERE c.post_id=p.id), 0) AS comments,
		EXISTS(SELECT 1 FROM likes WHERE post_id=p.id AND user_id=$1) AS liked,
		EXISTS(SELECT 1 FROM saves WHERE post_id=p.id AND user_id=$1) AS saved
		FROM posts p JOIN users u ON u.id=p.user_id WHERE p.id=$1`, id))
	if err == pgx.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	return &p, err
}

func (s *PostStore) GetTopicsByID(ctx context.Context, id int) ([]domain.TopicID, error) {
	var topics []domain.TopicID
	err := s.db.QueryRow(ctx, `SELECT topics_id FROM posts WHERE id=$1`, id).Scan(&topics)
	if err == pgx.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	return topics, err
}

func (s *PostStore) Update(ctx context.Context, postID, userID int, input *domain.UpdatePostInput) (*domain.Post, error) {
	p := &domain.Post{}
	var modifiedAt *time.Time
	err := s.db.QueryRow(ctx,
		`UPDATE posts SET body=$1, modified_at=NOW()
		 WHERE id=$2 AND user_id=$3
		 RETURNING id, body, topics_id, created_at, user_id, modified_at`,
		input.Body, postID, userID,
	).Scan(&p.ID, &p.Body, &p.TopicsID, &p.CreatedAt, &p.UserID, &modifiedAt)
	if err == pgx.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (s *PostStore) List(ctx context.Context, userID int, cursor string, limit int) ([]domain.Post, string, bool, error) {
	var cursorTime time.Time
	useCursor := false
	if cursor != "" {
		if t, err := time.Parse("2006-01-02T15:04:05.000Z", cursor); err == nil {
			cursorTime = t
			useCursor = true
		}
	}

	var rows pgx.Rows
	var err error

	if useCursor {
		rows, err = s.db.Query(ctx,
			`SELECT `+postCols+`,
				COALESCE((SELECT COUNT(*) FROM likes l WHERE l.post_id=p.id), 0) AS likes,
				COALESCE((SELECT COUNT(*) FROM saves s WHERE s.post_id=p.id), 0) AS saves,
				COALESCE((SELECT COUNT(*) FROM comments c WHERE c.post_id=p.id), 0) AS comments,
				EXISTS(SELECT 1 FROM likes WHERE post_id=p.id AND user_id=$1) AS liked,
				EXISTS(SELECT 1 FROM saves WHERE post_id=p.id AND user_id=$1) AS saved
			 FROM posts p JOIN users u ON u.id=p.user_id
			 WHERE p.created_at < $2
			 ORDER BY p.created_at DESC
			 LIMIT $3`, userID, cursorTime, limit)
	} else {
		rows, err = s.db.Query(ctx,
			`SELECT `+postCols+`,
				COALESCE((SELECT COUNT(*) FROM likes l WHERE l.post_id=p.id), 0) AS likes,
				COALESCE((SELECT COUNT(*) FROM saves s WHERE s.post_id=p.id), 0) AS saves,
				COALESCE((SELECT COUNT(*) FROM comments c WHERE c.post_id=p.id), 0) AS comments,
				EXISTS(SELECT 1 FROM likes WHERE post_id=p.id AND user_id=$1) AS liked,
				EXISTS(SELECT 1 FROM saves WHERE post_id=p.id AND user_id=$1) AS saved
			 FROM posts p JOIN users u ON u.id=p.user_id
			 ORDER BY p.created_at DESC
			 LIMIT $2`, userID, limit)
	}

	if err != nil {
		return nil, "", false, err
	}
	defer rows.Close()

	type postRow struct {
		domain.Post
		Likes    int  `db:"likes"`
		Saves    int  `db:"saves"`
		Comments int  `db:"comments"`
		Liked    bool `db:"liked"`
		Saved    bool `db:"saved"`
	}

	posts := []domain.Post{}
	for rows.Next() {
		var r postRow
		err := rows.Scan(
			&r.Post.ID, &r.Post.Body, &r.Post.TopicsID, &r.Post.CreatedAt,
			&r.Post.UserID, &r.Post.User.ID, &r.Post.User.Name, &r.Post.User.Username,
			&r.Likes, &r.Saves, &r.Comments, &r.Liked, &r.Saved,
		)
		if err != nil {
			return nil, "", false, err
		}
		r.Post.Likes = r.Likes
		r.Post.Saves = r.Saves
		r.Post.Comments = r.Comments
		r.Post.Liked = r.Liked
		r.Post.Saved = r.Saved
		posts = append(posts, r.Post)
	}
	if err := rows.Err(); err != nil {
		return nil, "", false, err
	}

	hasMore := len(posts) > 0 && len(posts) == limit
	return posts, cursor, hasMore, nil
}

func (s *PostStore) Delete(ctx context.Context, postID, userID int) error {
	tag, err := s.db.Exec(ctx, `DELETE FROM posts WHERE id=$1 AND user_id=$2`, postID, userID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (s *PostStore) Like(ctx context.Context, userID, postID int) error {
	_, err := s.db.Exec(ctx,
		`INSERT INTO likes (user_id, post_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`,
		userID, postID)
	return err
}

func (s *PostStore) Unlike(ctx context.Context, userID, postID int) error {
	_, err := s.db.Exec(ctx, `DELETE FROM likes WHERE user_id=$1 AND post_id=$2`, userID, postID)
	return err
}

func (s *PostStore) Save(ctx context.Context, userID, postID int) error {
	_, err := s.db.Exec(ctx,
		`INSERT INTO saves (user_id, post_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`,
		userID, postID)
	return err
}

func (s *PostStore) Unsave(ctx context.Context, userID, postID int) error {
	_, err := s.db.Exec(ctx, `DELETE FROM saves WHERE user_id=$1 AND post_id=$2`, userID, postID)
	return err
}

func (s *PostStore) GetLikes(ctx context.Context, postID, page, limit int) ([]domain.UserPublic, int, error) {
	var total int
	err := s.db.QueryRow(ctx, `SELECT COUNT(*) FROM likes WHERE post_id=$1`, postID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	rows, err := s.db.Query(ctx,
		`SELECT u.id, u.name, u.username FROM likes l
		 JOIN users u ON u.id=l.user_id
		 WHERE l.post_id=$1
		 ORDER BY l.user_id
		 LIMIT $2 OFFSET $3`, postID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	users := []domain.UserPublic{}
	for rows.Next() {
		var u domain.UserPublic
		if err := rows.Scan(&u.ID, &u.Name, &u.Username); err != nil {
			return nil, 0, err
		}
		users = append(users, u)
	}
	return users, total, rows.Err()
}

func (s *PostStore) GetSaves(ctx context.Context, postID, page, limit int) ([]domain.UserPublic, int, error) {
	var total int
	err := s.db.QueryRow(ctx, `SELECT COUNT(*) FROM saves WHERE post_id=$1`, postID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	rows, err := s.db.Query(ctx,
		`SELECT u.id, u.name, u.username FROM saves s
		 JOIN users u ON u.id=s.user_id
		 WHERE s.post_id=$1
		 ORDER BY s.user_id
		 LIMIT $2 OFFSET $3`, postID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	users := []domain.UserPublic{}
	for rows.Next() {
		var u domain.UserPublic
		if err := rows.Scan(&u.ID, &u.Name, &u.Username); err != nil {
			return nil, 0, err
		}
		users = append(users, u)
	}
	return users, total, rows.Err()
}
