package service

import (
	"simplestforum/internal/domain/entity"
	"simplestforum/internal/domain/usecase"
)

// Service is an interface which wraps the Transactioner.
type Service struct {
	entity.Transactioner
}

// Storages is a struct which contains all abstract Repositories.
type Storages struct {
	User         UserStorage
	Section      SectionStorage
	Topic        TopicStorage
	Post         PostStorage
	Notification NotificationStorage
}

// DoTransaction allows to wrap multiple service calls into a transaction.
func (s *Service) DoTransaction(sess entity.Session, f func() error) (err error) {
	if sess.Transaction != nil {
		return f()
	}

	tx, err := s.NewTransaction(sess.Ctx)
	if err != nil {
		return err
	}

	defer tx.RollbackUnlessCommitted()

	sess.Transaction = tx
	err = f()
	sess.Transaction = nil

	if err != nil {
		return err
	}

	return tx.Commit()
}

// NewServices creates a list of all abstract Services.
func NewServices(r *Storages) *usecase.Adapters {
	a := &usecase.Adapters{
		User:         NewUserService(r.User),
		Section:      NewSectionService(r.Section),
		Topic:        NewTopicService(r.Topic),
		Post:         NewPostService(r.Post),
		Notification: NewNotificationService(r.Notification),
	}

	a.User.AttachAdapters(a.Topic, a.Post)
	a.Section.AttachAdapters(a.Topic)
	a.Topic.AttachAdapters(a.User, a.Section, a.Post)
	a.Post.AttachAdapters(a.User, a.Topic)

	return a
}
