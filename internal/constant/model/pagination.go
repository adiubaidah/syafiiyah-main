package model

type Pagination struct {
	CurrentPage  int32 `json:"current_page"`
	TotalPages   int32 `json:"total_page"`
	TotalItems   int64 `json:"total_items"`
	ItemsPerPage int32 `json:"items_per_page"`
}
