package apimodel

type Pagination struct {
	Limit int64 `json:"limit"`
	Page  int64 `json:"page"`
}

type SortOrder string

const (
	SortOrderAsc  SortOrder = "ASC"
	SortOrderDesc SortOrder = "DESC"
)
