package usecase

import (
	"simplestforum/internal/domain"
	"simplestforum/internal/domain/entity"
)

// SectionUC is a Section usecase.
type SectionUC struct {
	sectionService SectionAdapter
}

// NewSectionUC instantiates a Section usecase.
func NewSectionUC(sectionService SectionAdapter) *SectionUC {
	return &SectionUC{
		sectionService: sectionService,
	}
}

// Add creates a new Section.
func (uc *SectionUC) Add(sess entity.Session, e *entity.SectionAdd) (*entity.Section, error) {
	if !sess.Level.AtLeast(entity.UserLevelAdmin) {
		return nil, domain.ErrForbidden
	}

	return uc.sectionService.Add(sess, e)
}

// Edit updates an existing Section.
func (uc *SectionUC) Edit(sess entity.Session, e *entity.SectionEdit) (*entity.Section, error) {
	if !sess.Level.AtLeast(entity.UserLevelAdmin) {
		return nil, domain.ErrForbidden
	}

	err := uc.sectionService.Edit(sess, e)
	if err != nil {
		return nil, err
	}

	return uc.ByID(sess, e.ID)
}

// Delete removes an existing Section.
func (uc *SectionUC) Delete(sess entity.Session, id int64) error {
	if !sess.Level.AtLeast(entity.UserLevelAdmin) {
		return domain.ErrForbidden
	}

	return uc.sectionService.Delete(sess, id)
}

// ByID returns a Section by its ID.
func (uc *SectionUC) ByID(sess entity.Session, id int64) (*entity.Section, error) {
	sections, err := uc.All(sess, &entity.SectionFilters{
		IDs: []int64{id},
	}, nil, nil)

	if err != nil {
		return nil, err
	}

	return sections[0], err
}

// All selects all Sections.
func (uc *SectionUC) All(sess entity.Session, f *entity.SectionFilters, p *entity.Pagination, s *entity.SectionSort) ([]*entity.Section, error) {
	return uc.sectionService.All(sess, f, p, s)
}
