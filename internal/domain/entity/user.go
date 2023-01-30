package entity

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

const PostsPerRank int64 = 50

// UserAdd is a structure used to insert a new User.
type UserAdd struct {
	Nickname string
	Password string
	ShowInfo bool

	*UserInfo
}

// UserEdit is a structure used to edit an existing User.
type UserEdit struct {
	ID          int64
	Nickname    *string
	Password    *string
	ShowInfo    *bool
	Rank        *int64
	CountTopics *int64
	CountPosts  *int64
	Level       *UserLevel
	Restriction *UserRestriction

	*UserInfo
}

// UserInfo contains secondary information about a User.
type UserInfo struct {
	Phone     *string
	Email     *string
	FirstName *string
	LastName  *string
}

// User is a general structure representing a User.
type User struct {
	ID          int64
	Nickname    string
	ShowInfo    bool
	Rank        int64
	Level       UserLevel
	Restriction UserRestriction
	CountTopics int64
	CountPosts  int64
	Topics      []*Topic
	Posts       []*Post
	CreatedAt   time.Time
	UpdatedAt   time.Time

	*UserInfo
}

func (u *User) PostsUntilNextRank() int64 {
	return u.Rank*PostsPerRank - u.CountPosts
}

type UserFilters struct {
	IDs            []int64
	RankFrom       *int64
	RankTo         *int64
	Level          *UserLevel
	Restriction    *UserRestriction
	CountPostsFrom *int64
	CountPostsTo   *int64
}

// UserSort represents User sorting options.
type UserSort struct {
	By    UserSortBy
	Order SortOrder
}

type UserSortBy string

const (
	UserSortByRank        UserSortBy = "RANK"
	UserSortByCountPosts  UserSortBy = "COUNT_POSTS"
	UserSortByCountTopics UserSortBy = "COUNT_TOPICS"
	UserSortByCreatedAt   UserSortBy = "CREATED_AT"
)

func (l UserLevel) AtLeast(atLeast UserLevel) bool {
	switch atLeast {
	case UserLevelMod:
		return l == UserLevelMod || l == UserLevelAdmin
	case UserLevelAdmin:
		return l == UserLevelAdmin
	}

	return false
}

func (r UserRestriction) AtLeast(atLeast UserRestriction) bool {
	switch atLeast {
	case UserRestrictionReadOnly:
		return r == UserRestrictionBanned || r == UserRestrictionReadOnly
	case UserRestrictionBanned:
		return r == UserRestrictionBanned
	}

	return false
}

// UsersEntityIDs returns the Ids of the users as a slice.
func UsersEntityIDs(users []*User) []int64 {
	ids := make([]int64, len(users))

	for i, user := range users {
		ids[i] = user.ID
	}

	return ids
}

// UsersMap returns an id => User map extracted out of users.
func UsersMap(users []*User) map[int64]*User {
	usersMap := make(map[int64]*User, len(users))

	for _, user := range users {
		usersMap[user.ID] = user
	}

	return usersMap
}
