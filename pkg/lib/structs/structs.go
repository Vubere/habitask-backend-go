package structs

type PaginationAndSort struct {
	Page          int    `json:"page"`
	PerPage       int    `json:"per_page"`
	SortBy        string `json:"sort_by"`
	SortDirection string `json:"sort_direction"`
}

func (p PaginationAndSort) GetOffset() int {
	return (p.Page - 1) * p.PerPage
}
