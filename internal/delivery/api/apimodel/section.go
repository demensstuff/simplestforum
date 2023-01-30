package apimodel

import "time"

type AddSectionInput struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
}

type EditSectionInput struct {
	ID          int64   `json:"id"`
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

type Section struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	CountTopics int64     `json:"count_topics"`
	Topics      []*Topic  `json:"topics"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type SectionFilters struct {
	Ids []int64 `json:"ids"`
}

type SectionSort struct {
	By    SectionSortBy `json:"by"`
	Order SortOrder     `json:"order"`
}

type SectionSortBy string

const (
	SectionSortByCountTopics SectionSortBy = "COUNT_TOPICS"
	SectionSortByCreatedAt   SectionSortBy = "CREATED_AT"
)
