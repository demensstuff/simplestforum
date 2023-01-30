package dto

import (
	"simplestforum/internal/delivery/api/apimodel"
	"simplestforum/internal/domain/entity"
)

func PaginationFromRest(p *apimodel.Pagination) *entity.Pagination {
	if p == nil {
		return nil
	}

	return &entity.Pagination{
		Limit: p.Limit,
		Page:  p.Page,
	}
}

// SortColumnToDB returns a column correlating to the entity constant.
func SortColumnToDB(e string) string {
	switch e {
	case string(entity.UserSortByCreatedAt):
		return "created_at"
	case string(entity.UserSortByRank):
		return "rank"
	case string(entity.UserSortByCountPosts):
		return "count_posts"
	case string(entity.UserSortByCountTopics):
		return "count_topics"
	case string(entity.TopicSortBySectionID):
		return "section_id"
	case string(entity.TopicSortByUserID):
		return "user_id"
	}

	return ""
}
