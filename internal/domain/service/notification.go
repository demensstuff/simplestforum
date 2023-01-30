package service

import (
	"simplestforum/internal/domain/entity"
)

// NotificationService represents a Section service.
type NotificationService struct {
	repo NotificationStorage

	Service
}

// NewNotificationService instantiates a NotificationService.
func NewNotificationService(repo NotificationStorage) *NotificationService {
	return &NotificationService{
		repo: repo,

		Service: Service{
			repo,
		},
	}
}

func (a *NotificationService) Add(sess entity.Session, e *entity.NotificationAdd) (*entity.Notification, error) {
	return a.repo.Insert(sess, e)
}

func (a *NotificationService) Clear(sess entity.Session, userID int64) error {
	return a.repo.DeleteByUserID(sess, userID)
}

func (a *NotificationService) All(sess entity.Session, userID int64, p *entity.Pagination) ([]*entity.Notification, error) {
	return a.repo.SelectAllByUserID(sess, p, userID)
}
