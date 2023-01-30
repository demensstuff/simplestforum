package dto

import (
	"simplestforum/internal/delivery/api/apimodel"
	"simplestforum/internal/domain/entity"
	"simplestforum/internal/infrastructure/dbmodel"
)

func NotificationsToRest(e []*entity.Notification) []*apimodel.Notification {
	if e == nil {
		return nil
	}

	notifications := make([]*apimodel.Notification, len(e))

	for i, notification := range e {
		notifications[i] = NotificationToRest(notification)
	}

	return notifications
}

func NotificationToRest(e *entity.Notification) *apimodel.Notification {
	if e == nil {
		return nil
	}

	return &apimodel.Notification{
		ID:        e.ID,
		UserID:    e.UserID,
		Text:      e.Text,
		CreatedAt: e.CreatedAt,
	}
}

func NotificationAddToDB(e *entity.NotificationAdd) *dbmodel.Notification {
	if e == nil {
		return nil
	}

	return &dbmodel.Notification{
		UserID: e.UserID,
		Text:   e.Text,
	}
}

func NotificationFromDB(n *dbmodel.Notification) *entity.Notification {
	if n == nil {
		return nil
	}

	return &entity.Notification{
		ID:        n.ID,
		UserID:    n.UserID,
		Text:      n.Text,
		CreatedAt: n.CreatedAt,
	}
}

func NotificationsFromDB(n []*dbmodel.Notification) []*entity.Notification {
	if n == nil {
		return nil
	}

	e := make([]*entity.Notification, len(n))

	for i, notification := range n {
		e[i] = NotificationFromDB(notification)
	}

	return e
}
