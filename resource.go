package fractal

// ResourceOption options for resource object
type ResourceOption struct {
	data        Any
	resourceKey string
	transformer Transformer
}

// ModResourceOption function to modify resource option
type ModResourceOption func(option *ResourceOption)

// WithData is an easy way to set data for resource
func WithData(data Any) ModResourceOption {
	return func(option *ResourceOption) {
		option.data = data
	}
}

// WithTransformer is an easy way to set transformer for resource
func WithTransformer(transformer Transformer) ModResourceOption {
	return func(option *ResourceOption) {
		option.transformer = transformer
	}
}

// WithResourceKey is an easy way to set resource key for resource
func WithResourceKey(key string) ModResourceOption {
	return func(option *ResourceOption) {
		option.resourceKey = key
	}
}
