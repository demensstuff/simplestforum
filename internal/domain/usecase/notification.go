package usecase

import (
	"simplestforum/internal/domain/entity"
)

// NotificationUC is a Notification usecase.
type NotificationUC struct {
	notificationService NotificationAdapter
}

// NewNotificationUC instantiates a Notification usecase.
func NewNotificationUC(notificationService NotificationAdapter) *NotificationUC {
	return &NotificationUC{
		notificationService: notificationService,
	}
}

func (uc *NotificationUC) Clear(sess entity.Session) error {
	return uc.notificationService.Clear(sess, sess.UserID)
}

func (uc *NotificationUC) All(sess entity.Session, p *entity.Pagination) ([]*entity.Notification, error) {
	return uc.notificationService.All(sess, sess.UserID, p)
}
