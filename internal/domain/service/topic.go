package service

import (
	"errors"
	"simplestforum/internal/domain"
	"simplestforum/internal/domain/entity"
	"simplestforum/internal/domain/usecase"
)

// TopicService represents a Section service.
type TopicService struct {
	repo TopicStorage

	userAdapter    usecase.UserAdapter
	sectionAdapter usecase.SectionAdapter
	postAdapter    usecase.PostAdapter

	Service
}

// NewTopicService instantiates a TopicService.
func NewTopicService(repo TopicStorage) *TopicService {
	return &TopicService{
		repo: repo,

		Service: Service{
			repo,
		},
	}
}

func (a *TopicService) AttachAdapters(userAdapter usecase.UserAdapter, sectionAdapter usecase.SectionAdapter, postAdapter usecase.PostAdapter) {
	a.userAdapter = userAdapter
	a.sectionAdapter = sectionAdapter
	a.postAdapter = postAdapter
}

// Add creates a new Topic.
func (a *TopicService) Add(sess entity.Session, e *entity.TopicAdd) (int64, error) {
	var id int64

	err := a.DoTransaction(sess, func() error {
		// Checking if the section exists
		err := a.sectionAdapter.ExistsByID(sess, e.SectionID)
		if err != nil {
			return err
		}

		// Inserting the topic
		id, err = a.repo.Insert(sess, e)

		return err
	})

	return id, err
}

// Edit modifies an existing Topic.
func (a *TopicService) Edit(sess entity.Session, e *entity.TopicEdit) error {
	return a.DoTransaction(sess, func() error {
		// Check if the ID is valid
		err := a.ExistsByID(sess, e.ID)
		if err != nil {
			return err
		}

		// If a new User ID is provided, check if the User exists
		if e.UserID != nil {
			err = a.userAdapter.ExistsByID(sess, *e.UserID)

			if err != nil {
				return err
			}
		}

		// If a new Section ID is provided, check if the Section exists
		if e.SectionID != nil {
			err = a.sectionAdapter.ExistsByID(sess, *e.SectionID)

			if err != nil {
				return err
			}
		}

		// Update the topic
		return a.repo.Update(sess, e)
	})
}

// Delete removes a single Topic entry along with the entities (posts) which depend on it.
func (a *TopicService) Delete(sess entity.Session, id int64) error {
	return a.DoTransaction(sess, func() error {
		// Check if the ID is valid
		err := a.ExistsByID(sess, id)
		if err != nil {
			return err
		}

		// Delete the topic
		err = a.repo.Delete(sess, id)
		if err != nil {
			return err
		}

		// Delete all its topics
		return a.postAdapter.MassDelete(sess, &entity.PostDelete{
			TopicIDs: []int64{id},
		})
	})
}

// MassDelete removes multiple Topic entries along with the entities (posts) which depend on them.
func (a *TopicService) MassDelete(sess entity.Session, e *entity.TopicDelete) error {
	return a.DoTransaction(sess, func() error {
		var (
			idsToDelete []int64
			err         error
		)

		// Get the Ids of topics which are about to be deleted
		if e.IDs != nil && e.SectionIDs == nil && e.UserIDs == nil {
			idsToDelete = e.IDs
		} else {
			idsToDelete, err = a.repo.IDsToDelete(sess, e)

			if err != nil {
				return err
			}
		}

		// Delete the topics
		err = a.repo.Delete(sess, idsToDelete...)
		if err != nil {
			return err
		}

		// Delete all their topics
		return a.postAdapter.MassDelete(sess, &entity.PostDelete{
			TopicIDs: idsToDelete,
		})
	})
}

// All fetches every Topic row matching the given filters, pagination, sorting and request options.
func (a *TopicService) All(sess entity.Session, f *entity.TopicFilters, p *entity.Pagination, s *entity.TopicSort) ([]*entity.Topic, error) {
	// If pagination was not set, use default
	if p == nil {
		p = entity.DefaultPagination
	}

	// If sorting options were not set, use default
	if s == nil {
		s = &entity.TopicSort{
			By:    entity.TopicSortByCreatedAt,
			Order: entity.SortOrderDesc,
		}
	}

	var topics []*entity.Topic

	err := a.DoTransaction(sess, func() error {
		var err error

		// Select the topics
		topics, err = a.repo.SelectAll(sess, f, p, s)

		if err != nil {
			return err
		}

		// If none were found, it's safe to return
		if len(topics) == 0 {
			return domain.NewError(domain.ErrCodeNotFound, "Topics not found")
		}

		// Retrieve Ids of the topics, users and sections, and build a map id => Topic to attach any embedded entities
		topicIDs, userIDs, sectionIDs := entity.TopicsEntityIDs(topics)
		topicsMap := entity.TopicsMap(topics)
		requestedFields := sess.RequestedFields

		// If we wish to fetch users
		if requestedFields.ContainsAny("user") {
			var users []*entity.User

			// Recursively change the requested fields to those for users
			sess.RequestedFields = requestedFields["user"]

			// Fetch the users
			users, err = a.userAdapter.All(sess, &entity.UserFilters{
				IDs: userIDs,
			}, nil, nil)

			// Put the initial requested fields back
			sess.RequestedFields = requestedFields

			if err != nil {
				return err
			}

			// If successfully, then attach the users to the respective topics
			usersMap := entity.UsersMap(users)

			for _, topic := range topicsMap {
				topic.User = usersMap[topic.UserID]
			}
		}

		// If we wish to fetch sections
		if requestedFields.ContainsAny("section") {
			var sections []*entity.Section

			// Recursively change the requested fields to those for sections
			sess.RequestedFields = requestedFields["section"]

			// Fetch the sections
			sections, err = a.sectionAdapter.All(sess, &entity.SectionFilters{
				IDs: sectionIDs,
			}, nil, nil)

			// Put the initial requested fields back
			sess.RequestedFields = requestedFields

			if err != nil {
				return err
			}

			// If successfully, then attach the users to the respective topics
			sectionsMap := entity.SectionsMap(sections)

			for _, topic := range topicsMap {
				topic.Section = sectionsMap[topic.SectionID]
			}
		}

		// If we wish to fetch posts
		if requestedFields.ContainsAny("posts") {
			var posts []*entity.Post

			// Recursively change the requested fields to those for posts
			sess.RequestedFields = requestedFields["posts"]

			// Fetch the posts
			posts, err = a.postAdapter.All(sess, &entity.PostFilters{
				TopicIDs: topicIDs,
			}, nil, nil)

			// Put the initial requested fields back
			sess.RequestedFields = requestedFields

			if err != nil {
				return err
			}

			// If successfully, then attach the posts to the respective users
			for _, post := range posts {
				topic := topicsMap[post.TopicID]
				topic.Posts = append(topic.Posts, post)
			}
		}

		return nil
	})

	return topics, err
}

// PlainByID returns a Topic by its ID.
func (a *TopicService) PlainByID(sess entity.Session, e *entity.PlainTopicByID) (*entity.Topic, error) {
	var topic *entity.Topic

	err := a.DoTransaction(sess, func() error {
		var err error

		topic, err = a.repo.SelectByID(sess, e.ID)

		if err != nil {
			return err
		}

		if e.FetchSection {
			topic.Section, err = a.sectionAdapter.PlainByID(sess, topic.SectionID)

			if err != nil {
				return err
			}
		}

		if e.FetchUser {
			topic.User, err = a.userAdapter.PlainByID(sess, topic.UserID)

			if err != nil {
				return err
			}
		}

		return nil
	})

	return topic, err
}

// ExistsByID return nil if the topic ID exists.
func (a *TopicService) ExistsByID(sess entity.Session, id int64) error {
	_, err := a.PlainByID(sess, &entity.PlainTopicByID{
		ID: id,
	})

	if err != nil {
		var domainErr *domain.Error

		if errors.As(err, domainErr) && domainErr.Is(domain.ErrNotFound) {
			domainErr.SetErrorMessage("Topic with ID %d not found", id)
		}

		return err
	}

	return nil
}
