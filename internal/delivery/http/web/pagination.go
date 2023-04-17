package web

type PaginationDTO struct {
	TotalRows   int64 `json:"total_rows"`
	Limit       int64 `json:"limit,omitempty"`
	CurrentPage int64 `json:"current_page,omitempty"`
	TotalPages  int64 `json:"total_pages"`
	Rows        any   `json:"rows"`
}
