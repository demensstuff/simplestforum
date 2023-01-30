package dto

import (
	"simplestforum/internal/delivery/api/apimodel"
	"simplestforum/internal/domain/entity"
	"simplestforum/internal/infrastructure/dbmodel"
)

func SectionsToRest(e []*entity.Section) []*apimodel.Section {
	if e == nil {
		return nil
	}

	sections := make([]*apimodel.Section, len(e))

	for i, section := range e {
		sections[i] = SectionToRest(section)
	}

	return sections
}

func SectionToRest(e *entity.Section) *apimodel.Section {
	if e == nil {
		return nil
	}

	return &apimodel.Section{
		ID:          e.ID,
		Name:        e.Name,
		Description: e.Description,
		CountTopics: e.CountTopics,
		Topics:      TopicsToRest(e.Topics),
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}
}

func SectionAddFromRest(s *apimodel.AddSectionInput) *entity.SectionAdd {
	if s == nil {
		return nil
	}

	return &entity.SectionAdd{
		Name:        s.Name,
		Description: s.Description,
	}
}

func SectionEditFromRest(s *apimodel.EditSectionInput) *entity.SectionEdit {
	if s == nil {
		return nil
	}

	return &entity.SectionEdit{
		ID:          s.ID,
		Name:        s.Name,
		Description: s.Description,
	}
}

func SectionFiltersFromRest(s *apimodel.SectionFilters) *entity.SectionFilters {
	if s == nil {
		return nil
	}

	return &entity.SectionFilters{
		IDs: s.Ids,
	}
}

func SectionSortFromRest(s *apimodel.SectionSort) *entity.SectionSort {
	if s == nil {
		return nil
	}

	return &entity.SectionSort{
		By:    entity.SectionSortBy(s.By),
		Order: entity.SortOrder(s.Order),
	}
}

func SectionAddToDB(e *entity.SectionAdd) *dbmodel.Section {
	if e == nil {
		return nil
	}

	return &dbmodel.Section{
		Name:        e.Name,
		Description: e.Description,
	}
}

func SectionFromDB(s *dbmodel.Section) *entity.Section {
	if s == nil {
		return nil
	}

	return &entity.Section{
		ID:          s.ID,
		Name:        s.Name,
		Description: s.Description,
		CountTopics: s.CountTopics,
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
	}
}

func SectionEditToDB(e *entity.SectionEdit) (*dbmodel.SectionUpdate, int64) {
	if e == nil {
		return nil, 0
	}

	sectionUpdate := &dbmodel.SectionUpdate{
		Name: e.Name,
	}

	var ptr *string

	if e.Description != nil {
		if *e.Description != "" {
			ptr = e.Description
		}

		sectionUpdate.Description = &ptr
	}

	return sectionUpdate, e.ID
}

func SectionFiltersToDB(e *entity.SectionFilters) *dbmodel.SectionFilters {
	if e == nil {
		return nil
	}

	return &dbmodel.SectionFilters{
		IDs: e.IDs,
	}
}

func SectionsFromDB(s []*dbmodel.Section) []*entity.Section {
	e := make([]*entity.Section, len(s))

	for i, section := range s {
		e[i] = SectionFromDB(section)
	}

	return e
}
