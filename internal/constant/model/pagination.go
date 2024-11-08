package model

type Pagination struct {
	CurrentPage  int `json:"current_page"`
	TotalPages   int `json:"total_page"`
	TotalItems   int `json:"total_items"`
	ItemsPerPage int `json:"items_per_page"`
}
