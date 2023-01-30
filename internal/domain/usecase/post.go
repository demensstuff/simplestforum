package usecase

import (
	"fmt"
	"simplestforum/internal/domain"
	"simplestforum/internal/domain/entity"
)

// PostUC is a Post usecase.
type PostUC struct {
	postService         PostAdapter
	userService         UserAdapter
	notificationService NotificationAdapter
}

// NewPostUC instantiates a Post usecase.
func NewPostUC(postService PostAdapter, userService UserAdapter, notificationService NotificationAdapter) *PostUC {
	return &PostUC{
		postService:         postService,
		userService:         userService,
		notificationService: notificationService,
	}
}

// Add creates a new Post.
func (uc *PostUC) Add(sess entity.Session, e *entity.PostAdd) (*entity.Post, error) {
	if sess.Restriction.AtLeast(entity.UserRestrictionReadOnly) {
		return nil, domain.ErrRestricted
	}

	e.UserID = sess.UserID

	// Creating a new Post
	postID, err := uc.postService.Add(sess, e)
	if err != nil {
		return nil, err
	}

	var (
		shouldUpdateRank bool
		newRank          int64
	)

	err = uc.userService.DoTransaction(sess, func() error {
		// Getting info about the author
		user, err := uc.userService.PlainByID(sess, e.UserID)
		if err != nil {
			return err
		}

		// Increasing the number of posts
		var (
			newCountPosts = user.CountPosts + 1
			userEdit      = &entity.UserEdit{
				ID:         user.ID,
				CountPosts: &newCountPosts,
			}
		)

		shouldUpdateRank = user.PostsUntilNextRank() == 1

		// If a new newRank is coming up, updating it as well
		if shouldUpdateRank {
			newRank = user.Rank + 1
			userEdit.Rank = &newRank
		}

		// Applying the new data to the User
		return uc.userService.Edit(sess, userEdit)
	})

	if err != nil {
		return uc.ByID(sess, postID)
	}

	// If the newRank was updated, notify
	if shouldUpdateRank {
		_, _ = uc.notificationService.Add(sess, &entity.NotificationAdd{
			UserID: e.UserID,
			Text:   fmt.Sprintf("You achieved the newRank %d, congratulations!", newRank),
		})
	}

	return uc.ByID(sess, postID)
}

// Edit updates an existing Post.
func (uc *PostUC) Edit(sess entity.Session, e *entity.PostEdit) (*entity.Post, error) {
	if sess.Restriction.AtLeast(entity.UserRestrictionReadOnly) {
		return nil, domain.ErrRestricted
	}

	var (
		isMod      = sess.Level.AtLeast(entity.UserLevelMod)
		postBefore *entity.Post
		post       *entity.Post
	)

	err := uc.postService.DoTransaction(sess, func() error {
		var err error

		// If current user doesn't have privileges, fetch the current state of the post to get its author ID
		if !isMod {
			postBefore, err = uc.postService.PlainByID(sess, &entity.PlainPostByID{
				ID:         e.ID,
				FetchTopic: e.TopicID != nil,
				FetchUser:  e.UserID != nil,
			})

			if err != nil {
				return err
			}

			// If it's someone else's post or the 'protected' fields are to be modified, return an error
			if sess.UserID != postBefore.UserID || e.TopicID != nil || e.UserID != nil {
				return domain.ErrForbidden
			}
		}

		// Apply the modification
		err = uc.postService.Edit(sess, e)
		if err != nil {
			return err
		}

		// Fetch the modified post with any embedded fields
		post, err = uc.ByID(sess, e.ID)

		return err
	})

	if err != nil {
		return nil, err
	}

	// If the topic was changed, notify about it
	if e.TopicID != nil {
		_, _ = uc.notificationService.Add(sess, &entity.NotificationAdd{
			UserID: post.UserID,
			Text:   fmt.Sprintf("Your post #%d was moved from topic %s to topic %s", post.ID, postBefore.Topic.Name, post.Topic.Name),
		})
	}

	// If the user was changed, notify both about it
	if e.UserID != nil {
		_, _ = uc.notificationService.Add(sess, &entity.NotificationAdd{
			UserID: postBefore.UserID,
			Text:   fmt.Sprintf("Your post #%d was assigned to user %s", post.ID, post.User.Nickname),
		})

		_, _ = uc.notificationService.Add(sess, &entity.NotificationAdd{
			UserID: post.UserID,
			Text:   fmt.Sprintf("Post #%d was assigned from user %s to you", post.ID, postBefore.User.Nickname),
		})
	}

	return post, nil
}

// Delete removes an existing Post.
func (uc *PostUC) Delete(sess entity.Session, id int64) error {
	if !sess.Level.AtLeast(entity.UserLevelMod) {
		return domain.ErrForbidden
	}

	var post *entity.Post

	err := uc.postService.DoTransaction(sess, func() error {
		var err error

		// Fetching the post to get its author ID
		post, err = uc.postService.PlainByID(sess, &entity.PlainPostByID{
			ID: id,
		})

		if err != nil {
			return err
		}

		// Deleting the post
		return uc.postService.Delete(sess, id)
	})

	if err != nil {
		return err
	}

	// Notifying the author
	_, _ = uc.notificationService.Add(sess, &entity.NotificationAdd{
		UserID: post.UserID,
		Text:   fmt.Sprintf("Your post #%d was removed", post.ID),
	})

	return nil
}

// ByID returns a Posts by its ID.
func (uc *PostUC) ByID(sess entity.Session, id int64) (*entity.Post, error) {
	posts, err := uc.All(sess, &entity.PostFilters{
		IDs: []int64{id},
	}, nil, nil)

	if err != nil {
		return nil, err
	}

	return posts[0], err
}

// All selects all Posts.
func (uc *PostUC) All(sess entity.Session, f *entity.PostFilters, p *entity.Pagination, s *entity.PostSort) ([]*entity.Post, error) {
	return uc.postService.All(sess, f, p, s)
}
