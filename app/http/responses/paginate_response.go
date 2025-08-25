package responses

type Pagination[T any] struct {
	Data    []T   `json:"data"`
	Total   int64 `json:"total"`
	Page    int   `json:"page"`
	PerPage int   `json:"per_page"`
}

func NewPagination[T any](data []T, page int, perPage int, total int64) Pagination[T] {
	return Pagination[T]{
		Data:    data,
		Total:   total,
		Page:    page,
		PerPage: perPage,
	}
}
