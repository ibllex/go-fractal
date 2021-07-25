package fractal

// Item resource can store any mixed data
type Item struct {
	data        Any
	meta        M
	resourceKey string
	transformer Transformer
}

// GetResourceKey get the resource key
func (item *Item) GetResourceKey() string {
	return item.resourceKey
}

// SetResourceKey set the resource key
func (item *Item) SetResourceKey(key string) {
	item.resourceKey = key
}

// GetData get the data
func (item *Item) GetData() Any {
	return item.data
}

// SetData set the data
func (item *Item) SetData(data Any) {
	item.data = data
}

// GetMeta get the meta data
func (item *Item) GetMeta() M {
	return item.meta
}

// SetMeta set the meta data
func (item *Item) SetMeta(meta M) Resource {
	item.meta = meta
	return item
}

// GetMetaValue get the meta value
func (item *Item) GetMetaValue(key string) Any {

	if item.meta == nil {
		return nil
	}

	if v, ok := item.meta[key]; ok {
		return v
	}
	return nil
}

// SetMetaValue get the meta value
func (item *Item) SetMetaValue(key string, value Any) Resource {
	if item.meta == nil {
		item.meta = make(M)
	}

	item.meta[key] = value
	return item
}

// GetTransformer get the transformer
func (item *Item) GetTransformer() Transformer {
	return item.transformer
}

// SetTransformer set the transformer
func (item *Item) SetTransformer(transformer Transformer) Resource {
	item.transformer = transformer
	return item
}

// NewItem create new item resource
func NewItem(opts ...ModResourceOption) *Item {
	opt := &ResourceOption{}

	for _, mod := range opts {
		mod(opt)
	}

	return &Item{
		data:        opt.data,
		resourceKey: opt.resourceKey,
		transformer: opt.transformer,
	}
}
