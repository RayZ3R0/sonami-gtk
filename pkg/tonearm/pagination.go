package tonearm

type Paginator[T any] interface {
	// Page returns the current page number
	Page() int

	// PageSize returns the number of items per page
	PageSize() int

	// TotalPages returns the total number of pages
	TotalPages() int

	// TotalItems returns the total number of items
	TotalItems() int

	// GetPage returns the items for the specified page
	GetPage(page int) ([]T, error)

	// GetAll returns all items
	GetAll() ([]T, error)
}
