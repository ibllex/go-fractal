package pagination

import (
	"math"
	"strconv"
	"strings"
)

// LengthAwarePaginator default paginator
type LengthAwarePaginator struct {
	items       []interface{}
	total       uint
	lastPage    uint
	perPage     uint
	currentPage uint
	pageName    string
	path        string
	fragment    string
	query       map[string]string
}

// GetItems get current items
func (p *LengthAwarePaginator) GetItems() []interface{} {
	return p.items
}

// SetItems set current items
func (p *LengthAwarePaginator) SetItems(items []interface{}) {
	p.items = items
}

// GetCurrentPage get the current page
func (p *LengthAwarePaginator) GetCurrentPage() uint {
	return p.currentPage
}

// GetLastPage get the last page
func (p *LengthAwarePaginator) GetLastPage() uint {
	return p.lastPage
}

// GetTotal get the total number of items being paginated
func (p *LengthAwarePaginator) GetTotal() uint {
	return p.total
}

// GetCount get the number of all items on the current page
func (p *LengthAwarePaginator) GetCount() uint {
	return uint(len(p.items))
}

// GetPerPage the number of items shown per page
func (p *LengthAwarePaginator) GetPerPage() uint {
	return p.perPage
}

// GetURL get the URL for a given page number
func (p *LengthAwarePaginator) GetURL(page uint) string {
	if page <= 0 {
		page = 1
	}

	query := p.GetPath()
	sep := "?"
	if strings.Contains(query, sep) {
		sep = "&"
	}
	query += sep

	for k, v := range p.query {
		query += (k + "=" + v + "&")
	}

	query += (p.pageName + "=" + strconv.FormatUint(uint64(page), 10))
	return query + p.buildFragment()
}

// GetPath get the base path for paginator generated URLs.
func (p *LengthAwarePaginator) GetPath() string {
	if p.path == "" {
		return "/"
	}
	return p.path
}

// SetPath set the base path for paginator generated URLs.
func (p *LengthAwarePaginator) SetPath(path string) *LengthAwarePaginator {
	p.path = path
	return p
}

func (p *LengthAwarePaginator) buildFragment() string {
	if p.fragment == "" {
		return ""
	}
	return "#" + p.fragment
}

func (p *LengthAwarePaginator) setCurrentPage(currentPage uint) *LengthAwarePaginator {
	if currentPage == 0 {
		currentPage = 1
	}

	p.currentPage = currentPage
	return p
}

// ModLengthAwarePaginator function to modify LengthAwarePaginator
type ModLengthAwarePaginator func(p *LengthAwarePaginator)

// WithCurrentPage is an easy way to set current page for paginator
func WithCurrentPage(currentPage uint) ModLengthAwarePaginator {
	return func(p *LengthAwarePaginator) {
		p.setCurrentPage(currentPage)
	}
}

// WithPageName is an easy way to set page name for paginator
func WithPageName(pageName string) ModLengthAwarePaginator {
	return func(p *LengthAwarePaginator) {
		p.pageName = pageName
	}
}

// WithPath is an easy way to set path for paginator
func WithPath(path string) ModLengthAwarePaginator {
	return func(p *LengthAwarePaginator) {
		if path != "/" {
			path = strings.TrimRight(path, "/")
		}

		p.path = path
	}
}

// WithQuery is an easy way to set query for paginator
func WithQuery(query map[string]string) ModLengthAwarePaginator {
	return func(p *LengthAwarePaginator) {
		p.query = query
	}
}

// WithFragment is an easy way to set fragment for paginator
func WithFragment(fragment string) ModLengthAwarePaginator {
	return func(p *LengthAwarePaginator) {
		p.fragment = fragment
	}
}

// NewLengthAwarePaginator create NewLengthAwarePaginator instance
func NewLengthAwarePaginator(items []interface{}, total uint, perPage uint, mods ...ModLengthAwarePaginator) *LengthAwarePaginator {
	lastPage := math.Ceil(float64(total) / float64(perPage))
	if lastPage < 1 {
		lastPage = 1
	}

	paginator := &LengthAwarePaginator{
		items:       items,
		total:       total,
		perPage:     perPage,
		pageName:    "page",
		lastPage:    uint(lastPage),
		currentPage: 1,
	}

	for _, mod := range mods {
		mod(paginator)
	}

	return paginator
}
