package entity

import "time"

// Topic is a general structure representing a Topic.
type Topic struct {
	ID         int64
	SectionID  int64
	Section    *Section
	Name       string
	UserID     int64
	User       *User
	CountPosts int64
	Posts      []*Post
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type TopicAdd struct {
	UserID    int64
	SectionID int64
	Name      string
}

type TopicEdit struct {
	ID        int64
	UserID    *int64
	SectionID *int64
	Name      *string
}

type PlainTopicByID struct {
	ID           int64
	FetchSection bool
	FetchUser    bool
}

type TopicFilters struct {
	IDs        []int64
	UserIDs    []int64
	SectionIDs []int64
}

type TopicDelete TopicFilters

type TopicSort struct {
	By    TopicSortBy
	Order SortOrder
}

type TopicSortBy string

const (
	TopicSortBySectionID  TopicSortBy = "SECTION_ID"
	TopicSortByUserID     TopicSortBy = "USER_ID"
	TopicSortByCountPosts TopicSortBy = "COUNT_POSTS"
	TopicSortByCreatedAt  TopicSortBy = "CREATED_AT"
)

// TopicsEntityIDs returns the Ids of the topics, users and sections as slices.
func TopicsEntityIDs(topics []*Topic) ([]int64, []int64, []int64) {
	ids := make([]int64, len(topics))
	userIDs := make([]int64, len(topics))
	sectionIDs := make([]int64, len(topics))

	for i, topic := range topics {
		ids[i] = topic.ID
		userIDs[i] = topic.UserID
		sectionIDs[i] = topic.SectionID
	}

	return ids, userIDs, sectionIDs
}

// TopicsMap returns an id => Topic map extracted out of topics.
func TopicsMap(topics []*Topic) map[int64]*Topic {
	topicsMap := make(map[int64]*Topic, len(topics))

	for _, topic := range topics {
		topicsMap[topic.ID] = topic
	}

	return topicsMap
}
