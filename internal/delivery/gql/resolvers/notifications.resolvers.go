package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.24

import (
	"context"
	"simplestforum/internal/delivery/api/apimodel"
	"simplestforum/internal/domain"
	"simplestforum/internal/domain/entity"
	"simplestforum/internal/dto"
)

// ClearNotifications is the resolver for the clearNotifications field.
func (r *mutationResolver) ClearNotifications(ctx context.Context) (bool, error) {
	sess := entity.GetSession(ctx)
	if !sess.IsAuthorized() {
		return false, domain.ErrNotAuthorized
	}

	err := r.Notification.Clear(sess)

	return err == nil, err
}

// ShowNotifications is the resolver for the showNotifications field.
func (r *queryResolver) ShowNotifications(ctx context.Context, p *apimodel.Pagination) ([]*apimodel.Notification, error) {
	sess := entity.GetSession(ctx)
	if !sess.IsAuthorized() {
		return nil, domain.ErrNotAuthorized
	}

	notifications, err := r.Notification.All(sess, dto.PaginationFromRest(p))
	if err != nil {
		return nil, err
	}

	return dto.NotificationsToRest(notifications), nil
}
