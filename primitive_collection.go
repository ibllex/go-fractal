package fractal

// PrimitiveCollection resource can store be a collection of any sort of primitive data
type PrimitiveCollection struct {
	Collection
}

// NewPrimitiveCollection create new primitive collection resource
func NewPrimitiveCollection(opts ...ModResourceOption) *PrimitiveCollection {
	opt := &ResourceOption{}

	for _, mod := range opts {
		mod(opt)
	}

	return &PrimitiveCollection{
		Collection: Collection{
			Item: Item{
				data:        opt.data,
				resourceKey: opt.resourceKey,
				transformer: opt.transformer,
			},
		},
	}
}
