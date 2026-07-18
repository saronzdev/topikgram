package domain

import "time"

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Birthday  time.Time `json:"birthday"`
	CreatedAt time.Time `json:"created_at"`
}

type UserPublic struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
}

type RegisterInput struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Birthday string `json:"birthday"`
}

type LoginInput struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}

type Post struct {
	ID        int        `json:"id"`
	UserID    int        `json:"user_id"`
	Body      string     `json:"body"`
	TopicsID  []TopicID  `json:"topics_id"`
	CreatedAt time.Time  `json:"created_at"`
	User      UserPublic `json:"user"`
	Likes     int        `json:"likes"`
	Saves     int        `json:"saves"`
	Comments  int        `json:"comments"`
	Saved     bool       `json:"saved"`
	Liked     bool       `json:"liked"`
}

type CreatePostInput struct {
	Body   string    `json:"body"`
	Topics []TopicID `json:"topics"`
}

type UpdatePostInput struct {
	Body string `json:"body"`
}

type PostListResponse struct {
	Posts      []Post `json:"posts"`
	NextCursor string `json:"next_cursor,omitempty"`
	HasMore    bool   `json:"has_more"`
}

type PaginatedUsers struct {
	Users []UserPublic `json:"users"`
	Total int          `json:"total"`
	Page  int          `json:"page"`
	Limit int          `json:"limit"`
}

type Comment struct {
	ID        int        `json:"id"`
	UserID    int        `json:"user_id"`
	PostID    int        `json:"post_id"`
	User      UserPublic `json:"user"`
	Content   string     `json:"content"`
	CreatedAt time.Time  `json:"created_at"`
}

type CreateCommentInput struct {
	PostID  int    `json:"post_id"`
	UserID  int    `json:"user_id"`
	Content string `json:"content"`
}

type AuthResponse struct {
	User User `json:"user"`
}

type TopicID int

const (
	General       TopicID = 0
	Programming   TopicID = 1
	Cybersecurity TopicID = 2
	Entertainment TopicID = 3
	Funny         TopicID = 4
	Art           TopicID = 5
	Sports        TopicID = 6
	Politics      TopicID = 7
	Science       TopicID = 8
	News          TopicID = 9
	Cinema        TopicID = 10
	Games         TopicID = 11
	Literature    TopicID = 12
	Travel        TopicID = 13
	Cuisine       TopicID = 14
	Tech          TopicID = 15
	Economy       TopicID = 16
	Health        TopicID = 17
	Philosophy    TopicID = 18
	Opinion       TopicID = 19
	Ad            TopicID = 20
)
