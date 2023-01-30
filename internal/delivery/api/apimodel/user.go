package apimodel

import "time"

// UserLevel represents User privileges.
type UserLevel string

// UserRestriction represents User restrictions.
type UserRestriction string

const (
	UserLevelNone  UserLevel = "NONE"
	UserLevelMod   UserLevel = "MOD"
	UserLevelAdmin UserLevel = "ADMIN"
)

const (
	UserRestrictionNone     UserRestriction = "NONE"
	UserRestrictionBanned   UserRestriction = "BANNED"
	UserRestrictionReadOnly UserRestriction = "READONLY"
)

// User is a general structure representing a User.
type User struct {
	ID          int64           `json:"id"`
	Nickname    string          `json:"nickname"`
	ShowInfo    bool            `json:"show_info"`
	Rank        int64           `json:"rank"`
	Level       UserLevel       `json:"level"`
	Restriction UserRestriction `json:"restriction"`
	UserInfo    *UserInfo       `json:"user_info"`
	CountTopics int64           `json:"count_topics"`
	CountPosts  int64           `json:"count_posts"`
	Topics      []*Topic        `json:"topics"`
	Posts       []*Post         `json:"posts"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

// AddUserInput is a structure which represents the input to create a new User.
type AddUserInput struct {
	Nickname  string  `json:"nickname"`
	Password  string  `json:"password"`
	ShowInfo  bool    `json:"show_info"`
	Phone     *string `json:"phone"`
	Email     *string `json:"email"`
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
}

// EditUserInput is a structure which represents the input to edit an existing User.
type EditUserInput struct {
	ID          *int64           `json:"id"`
	Nickname    *string          `json:"nickname"`
	Password    *string          `json:"password"`
	ShowInfo    *bool            `json:"show_info"`
	Phone       *string          `json:"phone"`
	Email       *string          `json:"email"`
	FirstName   *string          `json:"first_name"`
	LastName    *string          `json:"last_name"`
	Rank        *int64           `json:"rank"`
	CountTopics *int64           `json:"count_topics"`
	CountPosts  *int64           `json:"count_posts"`
	Level       *UserLevel       `json:"level"`
	Restriction *UserRestriction `json:"restriction"`
}

// UserInfo contains secondary information about a User.
type UserInfo struct {
	Phone     *string `json:"phone"`
	Email     *string `json:"email"`
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
}

type UserFilters struct {
	Ids            []int64          `json:"ids"`
	RankFrom       *int64           `json:"rank_from"`
	RankTo         *int64           `json:"rank_to"`
	Level          *UserLevel       `json:"level"`
	Restriction    *UserRestriction `json:"restriction"`
	CountPostsFrom *int64           `json:"count_posts_from"`
	CountPostsTo   *int64           `json:"count_posts_to"`
}

type UserSort struct {
	By    UserSortBy `json:"by"`
	Order SortOrder  `json:"order"`
}

type UserSortBy string

const (
	UserSortByRank        UserSortBy = "RANK"
	UserSortByCountPosts  UserSortBy = "COUNT_POSTS"
	UserSortByCountTopics UserSortBy = "COUNT_TOPICS"
	UserSortByCreatedAt   UserSortBy = "CREATED_AT"
)
