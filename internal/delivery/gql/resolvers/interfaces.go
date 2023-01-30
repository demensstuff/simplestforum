package resolvers

import (
	"simplestforum/internal/domain/entity"
)

// UserInteractor is an abstract User usecase.
type UserInteractor interface {
	Add(entity.Session, *entity.UserAdd) (*entity.User, error)
	Edit(entity.Session, *entity.UserEdit) (*entity.User, error)
	Delete(entity.Session, int64) error
	ByID(entity.Session, int64) (*entity.User, error)
	All(entity.Session, *entity.UserFilters, *entity.Pagination, *entity.UserSort) ([]*entity.User, error)
}

// TopicInteractor is an abstract Topic usecase.
type TopicInteractor interface {
	Add(entity.Session, *entity.TopicAdd) (*entity.Topic, error)
	Edit(entity.Session, *entity.TopicEdit) (*entity.Topic, error)
	Delete(entity.Session, int64) error
	ByID(entity.Session, int64) (*entity.Topic, error)
	All(entity.Session, *entity.TopicFilters, *entity.Pagination, *entity.TopicSort) ([]*entity.Topic, error)
}

// SectionInteractor is an abstract Section usecase.
type SectionInteractor interface {
	Add(entity.Session, *entity.SectionAdd) (*entity.Section, error)
	Edit(entity.Session, *entity.SectionEdit) (*entity.Section, error)
	Delete(entity.Session, int64) error
	ByID(entity.Session, int64) (*entity.Section, error)
	All(entity.Session, *entity.SectionFilters, *entity.Pagination, *entity.SectionSort) ([]*entity.Section, error)
}

// PostInteractor is an abstract Post usecase.
type PostInteractor interface {
	Add(entity.Session, *entity.PostAdd) (*entity.Post, error)
	Edit(entity.Session, *entity.PostEdit) (*entity.Post, error)
	Delete(entity.Session, int64) error
	ByID(entity.Session, int64) (*entity.Post, error)
	All(entity.Session, *entity.PostFilters, *entity.Pagination, *entity.PostSort) ([]*entity.Post, error)
}

// NotificationInteractor is an abstract Notification usecase.
type NotificationInteractor interface {
	Clear(entity.Session) error
	All(entity.Session, *entity.Pagination) ([]*entity.Notification, error)
}
