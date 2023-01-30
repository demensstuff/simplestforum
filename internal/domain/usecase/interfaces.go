package usecase

import (
	"simplestforum/internal/domain/entity"
)

// UserAdapter represents a set of User Service methods.
type UserAdapter interface {
	entity.Transactionable

	AttachAdapters(TopicAdapter, PostAdapter)

	Add(entity.Session, *entity.UserAdd) (*entity.User, error)
	Edit(entity.Session, *entity.UserEdit) error
	Delete(entity.Session, int64) error
	All(entity.Session, *entity.UserFilters, *entity.Pagination, *entity.UserSort) ([]*entity.User, error)

	ByLoginAndPassword(entity.Session, string, string) (*entity.User, error)
	PlainByID(entity.Session, int64) (*entity.User, error)

	ExistsByID(entity.Session, int64) error
}

// NotificationAdapter represents a set of Notification Service methods.
type NotificationAdapter interface {
	entity.Transactionable

	Add(entity.Session, *entity.NotificationAdd) (*entity.Notification, error)
	Clear(entity.Session, int64) error
	All(entity.Session, int64, *entity.Pagination) ([]*entity.Notification, error)
}

// SectionAdapter represents a set of Section Service methods.
type SectionAdapter interface {
	entity.Transactionable

	AttachAdapters(TopicAdapter)

	Add(entity.Session, *entity.SectionAdd) (*entity.Section, error)
	Edit(entity.Session, *entity.SectionEdit) error
	Delete(entity.Session, int64) error
	All(entity.Session, *entity.SectionFilters, *entity.Pagination, *entity.SectionSort) ([]*entity.Section, error)

	PlainByID(entity.Session, int64) (*entity.Section, error)
	ExistsByID(entity.Session, int64) error
}

// TopicAdapter represents a set of Topic Service methods.
type TopicAdapter interface {
	entity.Transactionable

	AttachAdapters(UserAdapter, SectionAdapter, PostAdapter)

	Add(entity.Session, *entity.TopicAdd) (int64, error)
	Edit(entity.Session, *entity.TopicEdit) error
	MassDelete(entity.Session, *entity.TopicDelete) error
	Delete(entity.Session, int64) error
	All(entity.Session, *entity.TopicFilters, *entity.Pagination, *entity.TopicSort) ([]*entity.Topic, error)

	PlainByID(entity.Session, *entity.PlainTopicByID) (*entity.Topic, error)
	ExistsByID(entity.Session, int64) error
}

// PostAdapter represents a set of Post Service methods.
type PostAdapter interface {
	entity.Transactionable

	AttachAdapters(UserAdapter, TopicAdapter)

	Add(entity.Session, *entity.PostAdd) (int64, error)
	Edit(entity.Session, *entity.PostEdit) error
	Delete(entity.Session, int64) error
	MassDelete(entity.Session, *entity.PostDelete) error
	All(entity.Session, *entity.PostFilters, *entity.Pagination, *entity.PostSort) ([]*entity.Post, error)

	PlainByID(entity.Session, *entity.PlainPostByID) (*entity.Post, error)
}
