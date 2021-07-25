package fractal

// Primitive resource can store any primitive data, like a string,
// integer, float etc
type Primitive struct {
	Item
}

// NewPrimitive create new primitive resource
func NewPrimitive(opts ...ModResourceOption) *Primitive {
	opt := &ResourceOption{}

	for _, mod := range opts {
		mod(opt)
	}

	return &Primitive{
		Item: Item{
			data:        opt.data,
			resourceKey: opt.resourceKey,
			transformer: opt.transformer,
		},
	}
}
