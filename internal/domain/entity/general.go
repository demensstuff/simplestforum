package entity

// Pagination represents pagination options.
type Pagination struct {
	Limit int64
	Page  int64
}

// SortOrder represents sort directions.
type SortOrder string

const (
	SortOrderAsc  SortOrder = "ASC"
	SortOrderDesc SortOrder = "DESC"
)

const (
	DefaultLimit int64 = 20
	DefaultPage  int64 = 1
)

var DefaultPagination = &Pagination{
	Limit: DefaultLimit,
	Page:  DefaultPage,
}
