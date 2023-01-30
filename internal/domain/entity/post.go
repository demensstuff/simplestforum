package entity

import "time"

// Post is a general structure representing a Post.
type Post struct {
	ID        int64
	Text      string
	UserID    int64
	User      *User
	TopicID   int64
	Topic     *Topic
	CreatedAt time.Time
	UpdatedAt time.Time
}

type PostAdd struct {
	UserID  int64
	TopicID int64
	Text    string
}

type PostEdit struct {
	ID      int64
	Text    *string
	UserID  *int64
	TopicID *int64
}

type PlainPostByID struct {
	ID         int64
	FetchTopic bool
	FetchUser  bool
}

type PostFilters struct {
	IDs      []int64
	UserIDs  []int64
	TopicIDs []int64
}

type PostDelete PostFilters

type PostSort struct {
	By    PostSortBy
	Order SortOrder
}

type PostSortBy string

const (
	PostSortByTopicID   PostSortBy = "TOPIC_ID"
	PostSortByUserID    PostSortBy = "USER_ID"
	PostSortByCreatedAt PostSortBy = "CREATED_AT"
)

// PostsEntityIDs returns the Ids of the users and topics as slices.
func PostsEntityIDs(posts []*Post) ([]int64, []int64) {
	userIDs := make([]int64, len(posts))
	topicIDs := make([]int64, len(posts))

	for i, post := range posts {
		userIDs[i] = post.UserID
		topicIDs[i] = post.TopicID
	}

	return userIDs, topicIDs
}

// PostsMap returns an id => Post map extracted out of posts.
func PostsMap(posts []*Post) map[int64]*Post {
	postsMap := make(map[int64]*Post, len(posts))

	for _, post := range posts {
		postsMap[post.ID] = post
	}

	return postsMap
}
