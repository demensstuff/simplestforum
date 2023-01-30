package service

import (
	"errors"
	"simplestforum/internal/domain"
	"simplestforum/internal/domain/entity"
	"simplestforum/internal/domain/usecase"

	"golang.org/x/crypto/bcrypt"
)

// UserService represents a User service.
type UserService struct {
	repo UserStorage

	topicAdapter usecase.TopicAdapter
	postAdapter  usecase.PostAdapter

	Service
}

// NewUserService instantiates a UserService.
func NewUserService(repo UserStorage) *UserService {
	return &UserService{
		repo: repo,

		Service: Service{
			repo,
		},
	}
}

func (a *UserService) AttachAdapters(topicAdapter usecase.TopicAdapter, postAdapter usecase.PostAdapter) {
	a.topicAdapter = topicAdapter
	a.postAdapter = postAdapter
}

// Add creates a new User.
func (a *UserService) Add(sess entity.Session, e *entity.UserAdd) (*entity.User, error) {
	// Hash the password
	var err error

	e.Password, err = a.hashPassword(e.Password)
	if err != nil {
		return nil, err
	}

	var user *entity.User

	err = a.DoTransaction(sess, func() error {
		// Check if the nickname is already registered
		err := a.nicknameAlreadyRegistered(sess, e.Nickname)
		if err != nil {
			return err
		}

		// Insert the main info
		user, err = a.repo.Insert(sess, e)
		if err != nil {
			return err
		}

		user.UserInfo = e.UserInfo

		// Insert additional info
		return a.repo.InsertInfo(sess, user.UserInfo, user.ID)
	})

	return user, err
}

// Edit modifies an existing User.
func (a *UserService) Edit(sess entity.Session, e *entity.UserEdit) error {
	// If the password is present, hash it
	if e.Password != nil {
		var err error

		*e.Password, err = a.hashPassword(*e.Password)

		if err != nil {
			return domain.NewErrorWrap(err, domain.ErrCodeInternal, "Cannot hash the password.")
		}
	}

	var user *entity.User

	return a.DoTransaction(sess, func() error {
		var err error

		// Check if the ID is valid
		err = a.ExistsByID(sess, e.ID)
		if err != nil {
			return err
		}

		// If the nickname is present, check if it was already registered
		if e.Nickname != nil {
			err = a.nicknameAlreadyRegistered(sess, *e.Nickname)

			if err != nil {
				return err
			}
		}

		// Update the main info
		err = a.repo.Update(sess, e)
		if err != nil {
			return err
		}

		// Update additional info if needed
		if e.UserInfo != nil {
			return a.repo.UpdateInfo(sess, e.UserInfo, user.ID)
		}

		return nil
	})
}

// Delete removes a single User entry along with the entities (posts and topics) which depend on it.
func (a *UserService) Delete(sess entity.Session, id int64) error {
	return a.DoTransaction(sess, func() error {
		// Check if the ID is valid
		err := a.ExistsByID(sess, id)
		if err != nil {
			return err
		}

		// Delete the user
		err = a.repo.Delete(sess, id)
		if err != nil {
			return err
		}

		// Delete all their topics
		err = a.topicAdapter.MassDelete(sess, &entity.TopicDelete{
			UserIDs: []int64{id},
		})
		if err != nil {
			return err
		}

		// Delete all their posts
		return a.postAdapter.Delete(sess, id)
	})
}

//nolint:godox,gocognit
// All fetches every User row matching the given filters, pagination, sorting and request options.
// TODO: simplify by allowing to attach embedded fields with two functions which manipulate with interfaces.
func (a *UserService) All(sess entity.Session, f *entity.UserFilters, p *entity.Pagination, s *entity.UserSort) ([]*entity.User, error) {
	// If pagination was not set, use default
	if p == nil {
		p = entity.DefaultPagination
	}

	// If sorting options were not set, use default
	if s == nil {
		s = &entity.UserSort{
			By:    entity.UserSortByCreatedAt,
			Order: entity.SortOrderDesc,
		}
	}

	var users []*entity.User

	err := a.DoTransaction(sess, func() error {
		var (
			requestedFields = sess.RequestedFields
			err             error
			infoRequested   = requestedFields.ContainsAny("user_info")
		)

		// Select the users, with or without additional info
		if infoRequested {
			users, err = a.repo.SelectAllWithInfo(sess, f, p, s)
		} else {
			users, err = a.repo.SelectAll(sess, f, p, s)
		}

		if err != nil {
			return err
		}

		// If none were found, it's safe to return
		if len(users) == 0 {
			return domain.NewError(domain.ErrCodeNotFound, "Users not found")
		}

		// If the private info is not supposed to be seen, hide it
		if infoRequested && !sess.Level.AtLeast(entity.UserLevelAdmin) {
			for _, user := range users {
				if !user.ShowInfo && user.ID != sess.UserID {
					user.UserInfo = nil
				}
			}
		}

		// Retrieve users' Ids and build a map id => User to attach any embedded entities
		userIDs := entity.UsersEntityIDs(users)
		usersMap := entity.UsersMap(users)

		// If we wish to fetch topics
		if requestedFields.ContainsAny("topics") {
			var topics []*entity.Topic

			// Recursively change the requested fields to those for topics
			sess.RequestedFields = requestedFields["topics"]

			// Fetch the topics
			topics, err = a.topicAdapter.All(sess, &entity.TopicFilters{
				UserIDs: userIDs,
			}, nil, nil)

			// Put the initial requested fields back
			sess.RequestedFields = requestedFields

			if err != nil {
				return err
			}

			// If successfully, then attach the topics to the respective users
			for _, topic := range topics {
				user := usersMap[topic.UserID]
				user.Topics = append(user.Topics, topic)
			}
		}

		// If we wish to fetch posts
		if requestedFields.ContainsAny("posts") {
			var posts []*entity.Post

			// Recursively change the requested fields to those for posts
			sess.RequestedFields = requestedFields["posts"]

			// Fetch the posts
			posts, err = a.postAdapter.All(sess, &entity.PostFilters{
				UserIDs: userIDs,
			}, nil, nil)

			// Put the initial requested fields back
			sess.RequestedFields = requestedFields

			if err != nil {
				return err
			}

			// If successfully, then attach the posts to the respective users
			for _, post := range posts {
				user := usersMap[post.UserID]
				user.Posts = append(user.Posts, post)
			}
		}

		return nil
	})

	return users, err
}

// ByLoginAndPassword returns a User by its login and password.
func (a *UserService) ByLoginAndPassword(sess entity.Session, nickname, password string) (*entity.User, error) {
	user, hashedPassword, err := a.repo.SelectByNicknameWithPassword(sess, nickname)

	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, domain.NewError(domain.ErrCodeInvalidCredentials, "Invalid login or password")
		}

		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, domain.NewError(domain.ErrCodeInvalidCredentials, "Invalid login or password")
		}

		return nil, domain.NewErrorWrap(err, domain.ErrCodeInternal, "Cannot compare passwords")
	}

	return user, nil
}

// PlainByID returns a User by its ID without any embedded fields.
func (a *UserService) PlainByID(sess entity.Session, id int64) (*entity.User, error) {
	return a.repo.SelectByID(sess, id)
}

// hashPassword attempts to hash the password string and return it (or an error).
func (a *UserService) hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", domain.NewErrorWrap(err, domain.ErrCodeInternal, "Cannot hash the password.")
	}

	return string(hashedPassword), nil
}

// nicknameAlreadyRegistered returns nil if the nickname doesn't exist.
func (a *UserService) nicknameAlreadyRegistered(sess entity.Session, nickname string) error {
	_, _, err := a.repo.SelectByNicknameWithPassword(sess, nickname)

	switch {
	case err == nil:
		return domain.NewError(domain.ErrCodeAlreadyExists, "Nickname %s is already registered", nickname)
	case !errors.Is(err, domain.ErrNotFound):
		return err
	}

	return nil
}

// ExistsByID return nil if the user ID exists.
func (a *UserService) ExistsByID(sess entity.Session, id int64) error {
	_, err := a.PlainByID(sess, id)
	if err != nil {
		var domainErr *domain.Error

		if errors.As(err, domainErr) && domainErr.Is(domain.ErrNotFound) {
			domainErr.SetErrorMessage("User with ID %d not found", id)
		}

		return err
	}

	return nil
}
