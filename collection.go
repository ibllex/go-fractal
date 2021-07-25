package fractal

// Collection resource can store be a collection of any sort of data
type Collection struct {
	Item
	paginator Paginator
	cursor    Cursor
}

// GetPaginator get the paginator instance
func (c *Collection) GetPaginator() Paginator {
	return c.paginator
}

// HasPaginator if the resource has a paginator implementation
func (c *Collection) HasPaginator() bool {
	return c.paginator != nil
}

// GetCursor get the cursor instance
func (c *Collection) GetCursor() Cursor {
	return c.cursor
}

// HasCursor if the resource has a cursor implementation
func (c *Collection) HasCursor() bool {
	return c.cursor != nil
}

// SetPaginator set the paginator instance
func (c *Collection) SetPaginator(paginator Paginator) *Collection {
	c.paginator = paginator
	return c
}

// SetCursor set the cursor instance
func (c *Collection) SetCursor(cursor Cursor) *Collection {
	c.cursor = cursor
	return c
}

// NewCollection create new collection resource
func NewCollection(opts ...ModResourceOption) *Collection {
	opt := &ResourceOption{}

	for _, mod := range opts {
		mod(opt)
	}

	return &Collection{
		Item: Item{
			data:        opt.data,
			resourceKey: opt.resourceKey,
			transformer: opt.transformer,
		},
	}
}
