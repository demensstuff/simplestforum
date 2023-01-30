package apimodel

import "time"

type AddTopicInput struct {
	SectionID int64  `json:"section_id"`
	Name      string `json:"name"`
}

type EditTopicInput struct {
	ID        int64   `json:"id"`
	UserID    *int64  `json:"user_id"`
	SectionID *int64  `json:"section_id"`
	Name      *string `json:"name"`
}

type Topic struct {
	ID         int64     `json:"id"`
	SectionID  int64     `json:"section_id"`
	Section    *Section  `json:"section"`
	Name       string    `json:"name"`
	UserID     int64     `json:"user_id"`
	User       *User     `json:"user"`
	CountPosts int64     `json:"count_posts"`
	Posts      []*Post   `json:"posts"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type TopicFilters struct {
	Ids        []int64 `json:"ids"`
	UserIds    []int64 `json:"user_ids"`
	SectionIds []int64 `json:"section_ids"`
}

type TopicSort struct {
	By    TopicSortBy `json:"by"`
	Order SortOrder   `json:"order"`
}

type TopicSortBy string

const (
	TopicSortBySectionID  TopicSortBy = "SECTION_ID"
	TopicSortByUserID     TopicSortBy = "USER_ID"
	TopicSortByCountPosts TopicSortBy = "COUNT_POSTS"
	TopicSortByCreatedAt  TopicSortBy = "CREATED_AT"
)
