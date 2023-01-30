package service

import (
	"errors"
	"simplestforum/internal/domain"
	"simplestforum/internal/domain/entity"
	"simplestforum/internal/domain/usecase"
)

// PostService represents a Section service.
type PostService struct {
	repo PostStorage

	userAdapter  usecase.UserAdapter
	topicAdapter usecase.TopicAdapter

	Service
}

// NewPostService instantiates a PostService.
func NewPostService(repo PostStorage) *PostService {
	return &PostService{
		repo: repo,

		Service: Service{
			repo,
		},
	}
}

func (a *PostService) AttachAdapters(userAdapter usecase.UserAdapter, topicAdapter usecase.TopicAdapter) {
	a.userAdapter = userAdapter
	a.topicAdapter = topicAdapter
}

// Add creates a new Post.
func (a *PostService) Add(sess entity.Session, e *entity.PostAdd) (int64, error) {
	var id int64

	err := a.DoTransaction(sess, func() error {
		// Checking if the topic exists
		err := a.topicAdapter.ExistsByID(sess, e.TopicID)
		if err != nil {
			return err
		}

		// Inserting the post
		id, err = a.repo.Insert(sess, e)

		return err
	})

	return id, err
}

// Edit modifies an existing Post.
func (a *PostService) Edit(sess entity.Session, e *entity.PostEdit) error {
	return a.DoTransaction(sess, func() error {
		// Check if the ID is valid
		err := a.existsByID(sess, e.ID)
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

		// If a new Topic ID is provided, check if the Topic exists
		if e.TopicID != nil {
			err = a.topicAdapter.ExistsByID(sess, *e.TopicID)

			if err != nil {
				return err
			}
		}

		// Update the post
		return a.repo.Update(sess, e)
	})
}

// Delete removes a single Post entry.
func (a *PostService) Delete(sess entity.Session, id int64) error {
	return a.DoTransaction(sess, func() error {
		// Check if the ID is valid
		err := a.existsByID(sess, id)
		if err != nil {
			return err
		}

		// Delete the post
		return a.repo.Delete(sess, id)
	})
}

// MassDelete removes multiple Post entries.
func (a *PostService) MassDelete(sess entity.Session, e *entity.PostDelete) error {
	return a.DoTransaction(sess, func() error {
		var idsToDelete []int64

		// Get the Ids of topics which are about to be deleted
		if e.IDs != nil && e.TopicIDs == nil && e.UserIDs == nil {
			idsToDelete = e.IDs
		} else {
			var err error

			idsToDelete, err = a.repo.IDsToDelete(sess, e)

			if err != nil {
				return err
			}
		}

		// Delete the posts
		return a.repo.Delete(sess, idsToDelete...)
	})
}

// All fetches every Post row matching the given filters, pagination, sorting and request options.
func (a *PostService) All(sess entity.Session, f *entity.PostFilters, p *entity.Pagination, s *entity.PostSort) ([]*entity.Post, error) {
	// If pagination was not set, use default
	if p == nil {
		p = entity.DefaultPagination
	}

	// If sorting options were not set, use default
	if s == nil {
		s = &entity.PostSort{
			By:    entity.PostSortByCreatedAt,
			Order: entity.SortOrderDesc,
		}
	}

	var posts []*entity.Post

	err := a.DoTransaction(sess, func() error {
		var err error

		// Select the posts
		posts, err = a.repo.SelectAll(sess, f, p, s)

		if err != nil {
			return err
		}

		// If none were found, it's safe to return
		if len(posts) == 0 {
			return domain.NewError(domain.ErrCodeNotFound, "Posts not found")
		}

		// Retrieve Ids of the users and topics, and build a map id => Post to attach any embedded entities
		userIDs, topicIDs := entity.PostsEntityIDs(posts)
		postsMap := entity.PostsMap(posts)
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

			// If successfully, then attach the users to the respective posts
			userMap := entity.UsersMap(users)

			for _, post := range postsMap {
				post.User = userMap[post.UserID]
			}
		}

		// If we wish to fetch sections
		if requestedFields.ContainsAny("topic") {
			var topics []*entity.Topic

			// Recursively change the requested fields to those for topics
			sess.RequestedFields = requestedFields["topic"]

			// Fetch the topics
			topics, err = a.topicAdapter.All(sess, &entity.TopicFilters{
				IDs: topicIDs,
			}, nil, nil)

			// Put the initial requested fields back
			sess.RequestedFields = requestedFields

			if err != nil {
				return err
			}

			// If successfully, then attach the users to the respective topics
			topicsMap := entity.TopicsMap(topics)

			for _, post := range postsMap {
				post.Topic = topicsMap[post.TopicID]
			}
		}

		return nil
	})

	return posts, err
}

// PlainByID returns a Post by its ID.
func (a *PostService) PlainByID(sess entity.Session, e *entity.PlainPostByID) (*entity.Post, error) {
	var post *entity.Post

	err := a.DoTransaction(sess, func() error {
		var err error

		post, err = a.repo.SelectByID(sess, e.ID)

		if err != nil {
			return err
		}

		if e.FetchTopic {
			post.Topic, err = a.topicAdapter.PlainByID(sess, &entity.PlainTopicByID{
				ID: post.TopicID,
			})

			if err != nil {
				return err
			}
		}

		if e.FetchUser {
			post.User, err = a.userAdapter.PlainByID(sess, post.UserID)

			if err != nil {
				return err
			}
		}

		return nil
	})

	return post, err
}

// existsByID return nil if the topic ID exists.
func (a *PostService) existsByID(sess entity.Session, id int64) error {
	_, err := a.PlainByID(sess, &entity.PlainPostByID{
		ID: id,
	})

	if err != nil {
		var domainErr *domain.Error

		if errors.As(err, domainErr) && domainErr.Is(domain.ErrNotFound) {
			domainErr.SetErrorMessage("Post with ID %d not found", id)
		}

		return err
	}

	return nil
}
