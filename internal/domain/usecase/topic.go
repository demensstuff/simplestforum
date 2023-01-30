package usecase

import (
	"fmt"
	"simplestforum/internal/domain"
	"simplestforum/internal/domain/entity"
)

// TopicUC is a Topic usecase.
type TopicUC struct {
	topicService        TopicAdapter
	userService         UserAdapter
	notificationService NotificationAdapter
}

// NewTopicUC instantiates a Topic usecase.
func NewTopicUC(topicService TopicAdapter, userService UserAdapter, notificationService NotificationAdapter) *TopicUC {
	return &TopicUC{
		topicService:        topicService,
		userService:         userService,
		notificationService: notificationService,
	}
}

// Add creates a new Topic.
func (uc *TopicUC) Add(sess entity.Session, e *entity.TopicAdd) (*entity.Topic, error) {
	if sess.Restriction.AtLeast(entity.UserRestrictionReadOnly) {
		return nil, domain.ErrRestricted
	}

	e.UserID = sess.UserID

	// Inserting the topic
	topicID, err := uc.topicService.Add(sess, e)
	if err != nil {
		return nil, err
	}

	// Increasing the user's topic count
	inc := int64(1)

	_ = uc.userService.Edit(sess, &entity.UserEdit{
		ID:          e.UserID,
		CountTopics: &inc,
	})

	return uc.ByID(sess, topicID)
}

// Edit updates an existing Topic.
func (uc *TopicUC) Edit(sess entity.Session, e *entity.TopicEdit) (*entity.Topic, error) {
	if sess.Restriction.AtLeast(entity.UserRestrictionReadOnly) {
		return nil, domain.ErrRestricted
	}

	var (
		isMod       = sess.Level.AtLeast(entity.UserLevelMod)
		topicBefore *entity.Topic
		topic       *entity.Topic
	)

	err := uc.topicService.DoTransaction(sess, func() error {
		var err error

		// If current user doesn't have privileges, fetch the current state of the topic to get its author ID
		if !isMod {
			topicBefore, err = uc.topicService.PlainByID(sess, &entity.PlainTopicByID{
				ID:           e.ID,
				FetchSection: e.SectionID != nil,
				FetchUser:    e.UserID != nil,
			})

			if err != nil {
				return err
			}

			// If it's someone else's topic or the 'protected' fields are to be modified, return an error
			if sess.UserID != topicBefore.UserID || e.SectionID != nil || e.UserID != nil {
				return domain.ErrForbidden
			}
		}

		// Apply the modification
		err = uc.topicService.Edit(sess, e)
		if err != nil {
			return err
		}

		// Fetch the modified topic with any embedded fields
		topic, err = uc.ByID(sess, e.ID)

		return err
	})

	if err != nil {
		return nil, err
	}

	// If the section was changed, notify about it
	if e.SectionID != nil {
		_, _ = uc.notificationService.Add(sess, &entity.NotificationAdd{
			UserID: topic.UserID,
			Text: fmt.Sprintf("Your topic %s was moved from section %s to section %s", topic.Name, topicBefore.Section.Name,
				topic.Section.Name),
		})
	}

	// If the user was changed, notify both about it
	if e.UserID != nil {
		_, _ = uc.notificationService.Add(sess, &entity.NotificationAdd{
			UserID: topicBefore.UserID,
			Text:   fmt.Sprintf("Your topic %s was assigned to user %s", topic.Name, topic.User.Nickname),
		})

		_, _ = uc.notificationService.Add(sess, &entity.NotificationAdd{
			UserID: topic.UserID,
			Text:   fmt.Sprintf("Topic %s was assigned from user %s to you", topic.Name, topicBefore.User.Nickname),
		})
	}

	return topic, nil
}

// Delete removes an existing Topic.
func (uc *TopicUC) Delete(sess entity.Session, id int64) error {
	if !sess.Level.AtLeast(entity.UserLevelMod) {
		return domain.ErrForbidden
	}

	var topic *entity.Topic

	err := uc.topicService.DoTransaction(sess, func() error {
		var err error

		// Fetching the topic to get its author ID
		topic, err = uc.topicService.PlainByID(sess, &entity.PlainTopicByID{
			ID: id,
		})

		if err != nil {
			return err
		}

		// Deleting the topic
		return uc.topicService.Delete(sess, id)
	})

	if err != nil {
		return err
	}

	// Notifying the author
	_, _ = uc.notificationService.Add(sess, &entity.NotificationAdd{
		UserID: topic.UserID,
		Text:   fmt.Sprintf("Your topic %s was removed", topic.Name),
	})

	return nil
}

// ByID returns a Topic by its ID.
func (uc *TopicUC) ByID(sess entity.Session, id int64) (*entity.Topic, error) {
	topics, err := uc.All(sess, &entity.TopicFilters{
		IDs: []int64{id},
	}, nil, nil)

	if err != nil {
		return nil, err
	}

	return topics[0], err
}

// All selects all Sections.
func (uc *TopicUC) All(sess entity.Session, f *entity.TopicFilters, p *entity.Pagination, s *entity.TopicSort) ([]*entity.Topic, error) {
	return uc.topicService.All(sess, f, p, s)
}
