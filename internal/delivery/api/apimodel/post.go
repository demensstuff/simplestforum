package apimodel

import "time"

type AddPostInput struct {
	TopicID int64  `json:"topic_id"`
	Text    string `json:"text"`
}

type EditPostInput struct {
	ID      int64   `json:"id"`
	Text    *string `json:"text"`
	UserID  *int64  `json:"user_id"`
	TopicID *int64  `json:"topic_id"`
}

type Post struct {
	ID        int64     `json:"id"`
	Text      string    `json:"text"`
	UserID    int64     `json:"user_id"`
	User      *User     `json:"user"`
	TopicID   int64     `json:"topic_id"`
	Topic     *Topic    `json:"topic"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PostFilters struct {
	Ids      []int64 `json:"ids"`
	UserIds  []int64 `json:"user_ids"`
	TopicIds []int64 `json:"topic_ids"`
}

type PostSort struct {
	By    PostSortBy `json:"by"`
	Order SortOrder  `json:"order"`
}

type PostSortBy string

const (
	PostSortByTopicID   PostSortBy = "TOPIC_ID"
	PostSortByUserID    PostSortBy = "USER_ID"
	PostSortByCreatedAt PostSortBy = "CREATED_AT"
)
