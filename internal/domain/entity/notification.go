package entity

import (
	"time"
)

type Notification struct {
	ID        int64
	UserID    int64
	Text      string
	CreatedAt time.Time
}

type NotificationAdd struct {
	UserID int64
	Text   string
}
