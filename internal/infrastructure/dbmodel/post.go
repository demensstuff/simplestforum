package dbmodel

import "time"

// Post is a structure which represents the 'posts' table entry.
type Post struct {
	ID        int64      `db:"id"`
	Text      string     `db:"name"`
	TopicID   int64      `db:"topic_id"`
	UserID    int64      `db:"user_id"`
	CreatedAt time.Time  `db:"created_at" insert:"false"`
	UpdatedAt time.Time  `db:"updated_at" insert:"false"`
	DeletedAt *time.Time `db:"deleted_at" insert:"false"`
}

// PostUpdate is a structure which is used to modify an existing entry in 'posts' table.
type PostUpdate struct {
	Text    *string `db:"text"`
	UserID  *int64  `db:"user_id"`
	TopicID *int64  `db:"topic_id"`
}

// PostFilters is a structure which represents post filters.
type PostFilters struct {
	IDs      []int64 `db:"id" sign:"="`
	UserIDs  []int64 `db:"user_id" sign:"="`
	TopicIDs []int64 `db:"topic_id" sign:"="`
}

// PostDelete is a structure which represents post filters for deletion.
type PostDelete PostFilters
