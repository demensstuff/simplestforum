package dbmodel

import "time"

// Section is a structure which represents the 'sections' table entry.
type Section struct {
	ID          int64      `db:"id"`
	Name        string     `db:"name"`
	Description *string    `db:"description"`
	CountTopics int64      `db:"count_topics" insert:"false"`
	CreatedAt   time.Time  `db:"created_at" insert:"false"`
	UpdatedAt   time.Time  `db:"updated_at" insert:"false"`
	DeletedAt   *time.Time `db:"deleted_at" insert:"false"`
}

// SectionUpdate is a structure which is used to modify an existing entry in 'sections' table.
type SectionUpdate struct {
	Name        *string  `db:"name"`
	Description **string `db:"description"`
}

// SectionFilters is a structure which represents section filters.
type SectionFilters struct {
	IDs []int64 `db:"id" sign:"="`
}
