package dbmodel

import "time"

// Topic is a structure which represents the 'topics' table entry.
type Topic struct {
	ID         int64      `db:"id"`
	SectionID  int64      `db:"section_id"`
	Name       string     `db:"name"`
	UserID     int64      `db:"user_id"`
	CountPosts int64      `db:"count_posts" insert:"false"`
	CreatedAt  time.Time  `db:"created_at" insert:"false"`
	UpdatedAt  time.Time  `db:"updated_at" insert:"false"`
	DeletedAt  *time.Time `db:"deleted_at" insert:"false"`
}

// TopicUpdate is a structure which is used to modify an existing entry in 'topics' table.
type TopicUpdate struct {
	UserID    *int64  `db:"user_id"`
	SectionID *int64  `db:"section_id"`
	Name      *string `db:"name"`
}

// TopicFilters is a structure which represents topic filters.
type TopicFilters struct {
	IDs        []int64 `db:"id" sign:"="`
	UserIDs    []int64 `db:"user_id" sign:"="`
	SectionIDs []int64 `db:"section_id" sign:"="`
}

// TopicDelete is a structure which represents topic filters for deletion.
type TopicDelete TopicFilters
