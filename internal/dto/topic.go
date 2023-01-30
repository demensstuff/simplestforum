package dto

import (
	"simplestforum/internal/delivery/api/apimodel"
	"simplestforum/internal/domain/entity"
	"simplestforum/internal/infrastructure/dbmodel"
)

func TopicsToRest(e []*entity.Topic) []*apimodel.Topic {
	if e == nil {
		return nil
	}

	topics := make([]*apimodel.Topic, len(e))

	for i, topic := range e {
		topics[i] = TopicToRest(topic)
	}

	return topics
}

func TopicToRest(e *entity.Topic) *apimodel.Topic {
	if e == nil {
		return nil
	}

	return &apimodel.Topic{
		ID:         e.ID,
		SectionID:  e.SectionID,
		Section:    SectionToRest(e.Section),
		Name:       e.Name,
		UserID:     e.UserID,
		User:       UserToRest(e.User),
		CountPosts: e.CountPosts,
		Posts:      PostsToRest(e.Posts),
		CreatedAt:  e.CreatedAt,
		UpdatedAt:  e.UpdatedAt,
	}
}

func TopicAddFromRest(t *apimodel.AddTopicInput) *entity.TopicAdd {
	if t == nil {
		return nil
	}

	return &entity.TopicAdd{
		SectionID: t.SectionID,
		Name:      t.Name,
	}
}

func TopicEditFromRest(t *apimodel.EditTopicInput) *entity.TopicEdit {
	if t == nil {
		return nil
	}

	return &entity.TopicEdit{
		ID:        t.ID,
		UserID:    t.UserID,
		SectionID: t.SectionID,
		Name:      t.Name,
	}
}

func TopicFiltersFromRest(t *apimodel.TopicFilters) *entity.TopicFilters {
	if t == nil {
		return nil
	}

	return &entity.TopicFilters{
		IDs:        t.Ids,
		UserIDs:    t.UserIds,
		SectionIDs: t.SectionIds,
	}
}

func TopicSortFromRest(t *apimodel.TopicSort) *entity.TopicSort {
	if t == nil {
		return nil
	}

	return &entity.TopicSort{
		By:    entity.TopicSortBy(t.By),
		Order: entity.SortOrder(t.Order),
	}
}

func TopicAddToDB(e *entity.TopicAdd) *dbmodel.Topic {
	if e == nil {
		return nil
	}

	return &dbmodel.Topic{
		SectionID: e.SectionID,
		Name:      e.Name,
		UserID:    e.UserID,
	}
}

func TopicFromDB(t *dbmodel.Topic) *entity.Topic {
	if t == nil {
		return nil
	}

	return &entity.Topic{
		ID:         t.ID,
		SectionID:  t.SectionID,
		Name:       t.Name,
		UserID:     t.UserID,
		CountPosts: t.CountPosts,
		CreatedAt:  t.CreatedAt,
		UpdatedAt:  t.UpdatedAt,
	}
}

func TopicEditToDB(e *entity.TopicEdit) (*dbmodel.TopicUpdate, int64) {
	if e == nil {
		return nil, 0
	}

	return &dbmodel.TopicUpdate{
		UserID:    e.UserID,
		SectionID: e.SectionID,
		Name:      e.Name,
	}, e.ID
}

func TopicFiltersToDB(e *entity.TopicFilters) *dbmodel.TopicFilters {
	if e == nil {
		return nil
	}

	return &dbmodel.TopicFilters{
		IDs:        e.IDs,
		UserIDs:    e.UserIDs,
		SectionIDs: e.SectionIDs,
	}
}

func TopicDeleteToDB(e *entity.TopicDelete) *dbmodel.TopicDelete {
	return (*dbmodel.TopicDelete)(TopicFiltersToDB((*entity.TopicFilters)(e)))
}

func TopicsFromDB(t []*dbmodel.Topic) []*entity.Topic {
	e := make([]*entity.Topic, len(t))

	for i, topic := range t {
		e[i] = TopicFromDB(topic)
	}

	return e
}
