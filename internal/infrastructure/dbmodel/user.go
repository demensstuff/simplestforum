package dbmodel

import (
	"time"
)

// UserInfo is a structure which represents the 'user_info' table entry.
type UserInfo struct {
	UserID    int64   `db:"user_id"`
	Phone     *string `db:"phone"`
	Email     *string `db:"email"`
	FirstName *string `db:"first_name"`
	LastName  *string `db:"last_name"`
}

// User is a structure which represents the 'users' table entry.
type User struct {
	ID          int64      `db:"id"`
	Nickname    string     `db:"nickname"`
	Password    string     `db:"password"`
	ShowInfo    bool       `db:"show_info"`
	Rank        int64      `db:"rank" insert:"false"`
	Level       *string    `db:"level" insert:"false"`
	Restriction *string    `db:"restriction" insert:"false"`
	CountTopics int64      `db:"count_topics" insert:"false"`
	CountPosts  int64      `db:"count_posts" insert:"false"`
	CreatedAt   time.Time  `db:"created_at" insert:"false"`
	UpdatedAt   time.Time  `db:"updated_at" insert:"false"`
	DeletedAt   *time.Time `db:"deleted_at" insert:"false"`
}

// UserWithInfo is a structure which represents a combined entry from the 'users' and 'user_info' table.
type UserWithInfo struct {
	User
	UserInfo
}

// UserUpdate is a structure used to store the optional fields to update a User.
type UserUpdate struct {
	Nickname    *string  `db:"nickname"`
	Password    *string  `db:"password"`
	ShowInfo    *bool    `db:"show_info"`
	Rank        *int64   `db:"rank"`
	CountTopics *int64   `db:"count_topics"`
	CountPosts  *int64   `db:"count_posts"`
	Level       **string `db:"level"`
	Restriction **string `db:"restriction"`
}

// UserFilters is a structure which represents all possible User filters.
type UserFilters struct {
	IDs            []int64 `db:"id" sign:"="`
	RankFrom       *int64  `db:"rank" sign:">="`
	RankTo         *int64  `db:"rank" sign:"<="`
	Level          *string `db:"level" sign:"="`
	Restriction    *string `db:"restriction" sign:"="`
	CountPostsFrom *int64  `db:"count_posts" sign:">="`
	CountPostsTo   *int64  `db:"count_topics" sign:"<="`
}
