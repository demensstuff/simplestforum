package service

import (
	"simplestforum/internal/domain/entity"
)

// UserStorage is an interface which declares methods to interact with any User storage.
type UserStorage interface {
	entity.Transactioner

	Insert(entity.Session, *entity.UserAdd) (*entity.User, error)
	Update(entity.Session, *entity.UserEdit) error
	Delete(entity.Session, int64) error
	SelectAll(entity.Session, *entity.UserFilters, *entity.Pagination, *entity.UserSort) ([]*entity.User, error)
	SelectAllWithInfo(entity.Session, *entity.UserFilters, *entity.Pagination, *entity.UserSort) ([]*entity.User, error)
	SelectByID(entity.Session, int64) (*entity.User, error)

	InsertInfo(entity.Session, *entity.UserInfo, int64) error
	UpdateInfo(entity.Session, *entity.UserInfo, int64) error
	SelectByNicknameWithPassword(entity.Session, string) (*entity.User, string, error)
}

// SectionStorage is an interface which declares methods to interact with any Section storage.
type SectionStorage interface {
	entity.Transactioner

	Insert(entity.Session, *entity.SectionAdd) (*entity.Section, error)
	Update(entity.Session, *entity.SectionEdit) error
	Delete(entity.Session, int64) error
	SelectByID(entity.Session, int64) (*entity.Section, error)
	SelectAll(entity.Session, *entity.SectionFilters, *entity.Pagination, *entity.SectionSort) ([]*entity.Section, error)
}

// TopicStorage is an interface which declares methods to interact with any Topic storage.
type TopicStorage interface {
	entity.Transactioner

	Insert(entity.Session, *entity.TopicAdd) (int64, error)
	Update(entity.Session, *entity.TopicEdit) error
	Delete(entity.Session, ...int64) error
	SelectByID(entity.Session, int64) (*entity.Topic, error)
	SelectAll(entity.Session, *entity.TopicFilters, *entity.Pagination, *entity.TopicSort) ([]*entity.Topic, error)

	IDsToDelete(entity.Session, *entity.TopicDelete) ([]int64, error)
}

// PostStorage is an interface which declares methods to interact with any Post storage.
type PostStorage interface {
	entity.Transactioner

	Insert(entity.Session, *entity.PostAdd) (int64, error)
	Update(entity.Session, *entity.PostEdit) error
	Delete(entity.Session, ...int64) error
	SelectByID(entity.Session, int64) (*entity.Post, error)
	SelectAll(entity.Session, *entity.PostFilters, *entity.Pagination, *entity.PostSort) ([]*entity.Post, error)

	IDsToDelete(entity.Session, *entity.PostDelete) ([]int64, error)
}

// NotificationStorage is an interface which declares methods to interact with any Notification storage.
type NotificationStorage interface {
	entity.Transactioner

	Insert(entity.Session, *entity.NotificationAdd) (*entity.Notification, error)
	DeleteByUserID(entity.Session, int64) error
	SelectAllByUserID(entity.Session, *entity.Pagination, int64) ([]*entity.Notification, error)
}
