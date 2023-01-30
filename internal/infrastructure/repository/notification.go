package repository

import (
	"simplestforum/internal/domain/entity"
	"simplestforum/internal/dto"
	"simplestforum/internal/infrastructure/dbmodel"
)

// NotificationRepository represents a Section Repository.
type NotificationRepository struct {
	*DBConn
}

// NewNotificationRepository instantiates a NotificationRepository.
func NewNotificationRepository(db *DBConn) *NotificationRepository {
	return &NotificationRepository{db}
}

// Insert creates a new Notification entry in the database and returns a Notification object.
func (r *NotificationRepository) Insert(sess entity.Session, e *entity.NotificationAdd) (*entity.Notification, error) {
	notification := dto.NotificationAddToDB(e)

	err := r.Wrap(sess, func(tx Gateway) error {
		stmt := tx.InsertInto("notifications").
			Returning("id", "created_at")

		insertNotNil(stmt, notification)

		return stmt.Load(&notification)
	})

	return dto.NotificationFromDB(notification), err
}

// DeleteByUserID removes existing Notifications (softly) by User ID.
func (r *NotificationRepository) DeleteByUserID(sess entity.Session, userID int64) error {
	return r.Wrap(sess, func(tx Gateway) error {
		_, err := tx.DeleteFrom("notifications").
			Where("user_id = ?", userID).
			Exec()

		return err
	})
}

// SelectAllByUserID returns all Notifications attributed to the given user.
func (r *NotificationRepository) SelectAllByUserID(sess entity.Session, p *entity.Pagination, userID int64) ([]*entity.Notification, error) {
	var notifications []*dbmodel.Notification

	err := r.Wrap(sess, func(tx Gateway) error {
		stmt := tx.Select("*").
			From("notifications").
			Where("user_id = ?", userID)

		if p != nil {
			stmt.Paginate(uint64(p.Page), uint64(p.Limit))
		}

		_, err := stmt.Load(&notifications)

		return err
	})

	return dto.NotificationsFromDB(notifications), err
}
