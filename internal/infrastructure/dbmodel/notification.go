package dbmodel

import "time"

type Notification struct {
	ID        int64     `db:"id"`
	UserID    int64     `db:"user_id"`
	Text      string    `db:"text"`
	CreatedAt time.Time `db:"created_at" insert:"false"`
}
