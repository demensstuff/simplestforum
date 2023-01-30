package usecase

import "simplestforum/internal/delivery/gql/resolvers"

// NewAdapters creates a list of all abstract Usecases.
func NewAdapters(s *Adapters) *resolvers.Interactors {
	return &resolvers.Interactors{
		User:         NewUserUC(s.User, s.Notification),
		Section:      NewSectionUC(s.Section),
		Topic:        NewTopicUC(s.Topic, s.User, s.Notification),
		Post:         NewPostUC(s.Post, s.User, s.Notification),
		Notification: NewNotificationUC(s.Notification),
	}
}

// Adapters is a list of all abstract Services.
type Adapters struct {
	User         UserAdapter
	Notification NotificationAdapter
	Section      SectionAdapter
	Topic        TopicAdapter
	Post         PostAdapter
}
