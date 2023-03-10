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

// AddSection is the resolver for the addSection field.
func (r *mutationResolver) AddSection(ctx context.Context, s apimodel.AddSectionInput) (*apimodel.Section, error) {
	sess := entity.GetSession(ctx)
	if !sess.IsAuthorized() {
		return nil, domain.ErrNotAuthorized
	}

	section, err := r.Section.Add(sess, dto.SectionAddFromRest(&s))
	if err != nil {
		return nil, err
	}

	return dto.SectionToRest(section), nil
}

// EditSection is the resolver for the editSection field.
func (r *mutationResolver) EditSection(ctx context.Context, s apimodel.EditSectionInput) (*apimodel.Section, error) {
	sess := entity.GetSession(ctx)
	if !sess.IsAuthorized() {
		return nil, domain.ErrNotAuthorized
	}

	section, err := r.Section.Edit(sess, dto.SectionEditFromRest(&s))
	if err != nil {
		return nil, err
	}

	return dto.SectionToRest(section), nil
}

// DeleteSection is the resolver for the deleteSection field.
func (r *mutationResolver) DeleteSection(ctx context.Context, id int64) (bool, error) {
	sess := entity.GetSession(ctx)
	if !sess.IsAuthorized() {
		return false, domain.ErrNotAuthorized
	}

	err := r.Section.Delete(sess, id)

	return err == nil, err
}

// ShowSection is the resolver for the showSection field.
func (r *queryResolver) ShowSection(ctx context.Context, id int64) (*apimodel.Section, error) {
	sess := entity.GetSession(ctx)

	section, err := r.Section.ByID(sess, id)
	if err != nil {
		return nil, err
	}

	return dto.SectionToRest(section), nil
}

// ShowSections is the resolver for the showSections field.
func (r *queryResolver) ShowSections(ctx context.Context, f *apimodel.SectionFilters, p *apimodel.Pagination, s *apimodel.SectionSort) ([]*apimodel.Section, error) {
	sess := entity.GetSession(ctx)

	sections, err := r.Section.All(sess, dto.SectionFiltersFromRest(f), dto.PaginationFromRest(p), dto.SectionSortFromRest(s))
	if err != nil {
		return nil, err
	}

	return dto.SectionsToRest(sections), nil
}
