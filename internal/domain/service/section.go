package service

import (
	"errors"
	"simplestforum/internal/domain"
	"simplestforum/internal/domain/entity"
	"simplestforum/internal/domain/usecase"
)

// SectionService represents a Section service.
type SectionService struct {
	repo SectionStorage

	topicAdapter usecase.TopicAdapter

	Service
}

// NewSectionService instantiates a SectionService.
func NewSectionService(repo SectionStorage) *SectionService {
	return &SectionService{
		repo: repo,

		Service: Service{
			repo,
		},
	}
}

func (a *SectionService) AttachAdapters(topicAdapter usecase.TopicAdapter) {
	a.topicAdapter = topicAdapter
}

// Add creates a new Section.
func (a *SectionService) Add(sess entity.Session, e *entity.SectionAdd) (*entity.Section, error) {
	return a.repo.Insert(sess, e)
}

// Edit modifies an existing Section.
func (a *SectionService) Edit(sess entity.Session, e *entity.SectionEdit) error {
	// Check if the ID is valid
	err := a.ExistsByID(sess, e.ID)
	if err != nil {
		return err
	}

	// Update the section
	return a.repo.Update(sess, e)
}

// Delete removes a single Section entry along with the entities (topics) which depend on it.
func (a *SectionService) Delete(sess entity.Session, id int64) error {
	return a.DoTransaction(sess, func() error {
		// Check if the ID is valid
		err := a.ExistsByID(sess, id)
		if err != nil {
			return err
		}

		// Delete the section
		err = a.repo.Delete(sess, id)
		if err != nil {
			return err
		}

		// Delete all its topics
		return a.topicAdapter.MassDelete(sess, &entity.TopicDelete{
			UserIDs: []int64{id},
		})
	})
}

// All fetches every Section row matching the given filters, pagination, sorting and request options.
func (a *SectionService) All(sess entity.Session, f *entity.SectionFilters, p *entity.Pagination, s *entity.SectionSort) ([]*entity.Section, error) {
	// If pagination was not set, use default
	if p == nil {
		p = entity.DefaultPagination
	}

	// If sorting options were not set, use default
	if s == nil {
		s = &entity.SectionSort{
			By:    entity.SectionSortByCreatedAt,
			Order: entity.SortOrderDesc,
		}
	}

	var sections []*entity.Section

	err := a.DoTransaction(sess, func() error {
		var err error

		// Select the sections
		sections, err = a.repo.SelectAll(sess, f, p, s)

		if err != nil {
			return err
		}

		// If none were found, it's safe to return
		if len(sections) == 0 {
			return domain.NewError(domain.ErrCodeNotFound, "Sections not found")
		}

		// Retrieve Ids of the sections and build a map id => Section to attach any embedded entities
		sectionIDs := entity.SectionsEntityIDs(sections)
		sectionsMap := entity.SectionsMap(sections)
		requestedFields := sess.RequestedFields

		// If we wish to fetch topics
		if requestedFields.ContainsAny("topics") {
			var topics []*entity.Topic

			// Recursively change the requested fields to those for topics
			sess.RequestedFields = requestedFields["topics"]

			// Fetch the topics
			topics, err = a.topicAdapter.All(sess, &entity.TopicFilters{
				SectionIDs: sectionIDs,
			}, nil, nil)

			// Put the initial requested fields back
			sess.RequestedFields = requestedFields

			if err != nil {
				return err
			}

			// If successfully, then attach the topics to the respective sections
			if err == nil {
				for _, topic := range topics {
					section := sectionsMap[topic.SectionID]
					section.Topics = append(section.Topics, topic)
				}
			}
		}

		return nil
	})

	return sections, err
}

// PlainByID returns a Section by its ID without any embedded fields.
func (a *SectionService) PlainByID(sess entity.Session, id int64) (*entity.Section, error) {
	return a.repo.SelectByID(sess, id)
}

// ExistsByID return nil if the section ID exists.
func (a *SectionService) ExistsByID(sess entity.Session, id int64) error {
	_, err := a.PlainByID(sess, id)
	if err != nil {
		var domainErr *domain.Error

		if errors.As(err, domainErr) && domainErr.Is(domain.ErrNotFound) {
			domainErr.SetErrorMessage("Section with ID %d not found", id)
		}

		return err
	}

	return nil
}
