package usecase

import (
	"fmt"
	"simplestforum/internal/domain"
	"simplestforum/internal/domain/entity"
)

// UserUC is a User usecase.
type UserUC struct {
	userService         UserAdapter
	notificationService NotificationAdapter
}

// NewUserUC instantiates a User usecase.
func NewUserUC(userService UserAdapter, notificationService NotificationAdapter) *UserUC {
	return &UserUC{
		userService:         userService,
		notificationService: notificationService,
	}
}

// Add creates a new User.
func (uc *UserUC) Add(sess entity.Session, e *entity.UserAdd) (*entity.User, error) {
	// Adding a user
	user, err := uc.userService.Add(sess, e)
	if err != nil {
		return nil, err
	}

	// Adding a welcome notification
	_, _ = uc.notificationService.Add(sess, &entity.NotificationAdd{
		UserID: user.ID,
		Text:   "Welcome to the forum!",
	})

	return user, nil
}

// Edit updates an existing User.
func (uc *UserUC) Edit(sess entity.Session, e *entity.UserEdit) (*entity.User, error) {
	var (
		shouldRetrieveUserBefore = e.Level != nil || e.Restriction != nil
		protectedFieldsChanged   = shouldRetrieveUserBefore || e.Rank != nil || e.CountPosts != nil || e.CountTopics != nil
		userBefore               *entity.User
		user                     *entity.User
	)

	err := uc.userService.DoTransaction(sess, func() error {
		var err error

		// If we're editing another user or 'protected' fields, and we're not the admin, return an error
		if (e.ID != sess.UserID || protectedFieldsChanged) && !sess.Level.AtLeast(entity.UserLevelAdmin) {
			return domain.ErrForbidden
		}

		// If the level or restriction were changed, fetch the current state of the User for the previous values
		if shouldRetrieveUserBefore {
			userBefore, err = uc.userService.PlainByID(sess, e.ID)
			if err != nil {
				return err
			}
		}

		// Apply the modifications
		err = uc.userService.Edit(sess, e)
		if err != nil {
			return err
		}

		// Fetch the modified user with any embedded fields
		user, err = uc.ByID(sess, e.ID)

		return err
	})

	if err != nil {
		return nil, err
	}

	// If the level was changed, notify about it
	if e.Level != nil && userBefore.Level != user.Level {
		_, _ = uc.notificationService.Add(sess, &entity.NotificationAdd{
			UserID: user.ID,
			Text:   fmt.Sprintf("Your privilege level has been changed to %s", user.Level),
		})
	}

	// If the restriction was changed, notify about it
	if e.Restriction != nil && userBefore.Restriction != user.Restriction {
		_, _ = uc.notificationService.Add(sess, &entity.NotificationAdd{
			UserID: user.ID,
			Text:   fmt.Sprintf("Your restriction level has been changed to %s", user.Restriction),
		})
	}

	return user, nil
}

// Delete removes an existing User.
func (uc *UserUC) Delete(sess entity.Session, id int64) error {
	if !sess.Level.AtLeast(entity.UserLevelAdmin) {
		return domain.ErrForbidden
	}

	return uc.userService.Delete(sess, id)
}

// ByID returns a User by its ID.
func (uc *UserUC) ByID(sess entity.Session, id int64) (*entity.User, error) {
	users, err := uc.All(sess, &entity.UserFilters{
		IDs: []int64{id},
	}, nil, nil)

	if err != nil {
		return nil, err
	}

	return users[0], err
}

// All selects all Users.
func (uc *UserUC) All(sess entity.Session, f *entity.UserFilters, p *entity.Pagination, s *entity.UserSort) ([]*entity.User, error) {
	return uc.userService.All(sess, f, p, s)
}
