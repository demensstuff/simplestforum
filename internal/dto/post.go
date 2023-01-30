package dto

import (
	"simplestforum/internal/delivery/api/apimodel"
	"simplestforum/internal/domain/entity"
	"simplestforum/internal/infrastructure/dbmodel"
)

func PostsToRest(e []*entity.Post) []*apimodel.Post {
	if e == nil {
		return nil
	}

	posts := make([]*apimodel.Post, len(e))

	for i, post := range e {
		posts[i] = PostToRest(post)
	}

	return posts
}

func PostToRest(e *entity.Post) *apimodel.Post {
	if e == nil {
		return nil
	}

	return &apimodel.Post{
		ID:        e.ID,
		Text:      e.Text,
		UserID:    e.UserID,
		User:      UserToRest(e.User),
		TopicID:   e.TopicID,
		Topic:     TopicToRest(e.Topic),
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}

func PostAddFromRest(p *apimodel.AddPostInput) *entity.PostAdd {
	if p == nil {
		return nil
	}

	return &entity.PostAdd{
		TopicID: p.TopicID,
		Text:    p.Text,
	}
}

func PostEditFromRest(p *apimodel.EditPostInput) *entity.PostEdit {
	if p == nil {
		return nil
	}

	return &entity.PostEdit{
		ID:      p.ID,
		Text:    p.Text,
		UserID:  p.UserID,
		TopicID: p.TopicID,
	}
}

func PostFiltersFromRest(p *apimodel.PostFilters) *entity.PostFilters {
	if p == nil {
		return nil
	}

	return &entity.PostFilters{
		IDs:      p.Ids,
		UserIDs:  p.UserIds,
		TopicIDs: p.TopicIds,
	}
}

func PostSortFromRest(p *apimodel.PostSort) *entity.PostSort {
	if p == nil {
		return nil
	}

	return &entity.PostSort{
		By:    entity.PostSortBy(p.By),
		Order: entity.SortOrder(p.Order),
	}
}

func PostAddToDB(e *entity.PostAdd) *dbmodel.Post {
	if e == nil {
		return nil
	}

	return &dbmodel.Post{
		Text:    e.Text,
		TopicID: e.TopicID,
		UserID:  e.UserID,
	}
}

func PostFromDB(p *dbmodel.Post) *entity.Post {
	if p == nil {
		return nil
	}

	return &entity.Post{
		ID:        p.ID,
		Text:      p.Text,
		UserID:    p.UserID,
		TopicID:   p.TopicID,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}

func PostEditToDB(e *entity.PostEdit) (*dbmodel.PostUpdate, int64) {
	if e == nil {
		return nil, 0
	}

	return &dbmodel.PostUpdate{
		Text:    e.Text,
		UserID:  e.UserID,
		TopicID: e.TopicID,
	}, e.ID
}

func PostFiltersToDB(e *entity.PostFilters) *dbmodel.PostFilters {
	if e == nil {
		return nil
	}

	return &dbmodel.PostFilters{
		IDs:      e.IDs,
		UserIDs:  e.UserIDs,
		TopicIDs: e.TopicIDs,
	}
}

func PostDeleteToDB(e *entity.PostDelete) *dbmodel.PostDelete {
	return (*dbmodel.PostDelete)(PostFiltersToDB((*entity.PostFilters)(e)))
}

func PostsFromDB(p []*dbmodel.Post) []*entity.Post {
	e := make([]*entity.Post, len(p))

	for i, Post := range p {
		e[i] = PostFromDB(Post)
	}

	return e
}
