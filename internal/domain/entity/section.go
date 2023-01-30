package entity

import "time"

// Section is a general structure representing a Section.
type Section struct {
	ID          int64
	Name        string
	Description *string
	CountTopics int64
	Topics      []*Topic
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type SectionAdd struct {
	Name        string
	Description *string
}

type SectionEdit struct {
	ID          int64
	Name        *string
	Description *string
}

type SectionFilters struct {
	IDs []int64
}

type SectionSort struct {
	By    SectionSortBy
	Order SortOrder
}

type SectionSortBy string

const (
	SectionSortByCountTopics SectionSortBy = "COUNT_TOPICS"
	SectionSortByCreatedAt   SectionSortBy = "CREATED_AT"
)

// SectionsEntityIDs returns the Ids of the sections as a slice.
func SectionsEntityIDs(sections []*Section) []int64 {
	ids := make([]int64, len(sections))

	for i, section := range sections {
		ids[i] = section.ID
	}

	return ids
}

// SectionsMap returns an id => Section map extracted out of sections.
func SectionsMap(sections []*Section) map[int64]*Section {
	sectionsMap := make(map[int64]*Section, len(sections))

	for _, section := range sections {
		sectionsMap[section.ID] = section
	}

	return sectionsMap
}
