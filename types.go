package fractal

// Paginator interface
type Paginator interface {
	GetCurrentPage() uint
	GetLastPage() uint
	GetTotal() uint
	GetCount() uint
	GetPerPage() uint
	GetURL(page uint) string
}

// Cursor interface
type Cursor interface {
	GetCurrent() string
	GetPrev() string
	GetNext() string
	GetCount() uint
}

// Resource interface
type Resource interface {
	GetResourceKey() string
	SetResourceKey(key string)
	GetData() Any
	SetData(data Any)
	GetTransformer() Transformer
	SetTransformer(Transformer) Resource
	GetMeta() M
	SetMeta(M) Resource
	GetMetaValue(key string) Any
	SetMetaValue(key string, value Any) Resource
}

// Transformer interface
type Transformer interface {
	Transform(data Any) M
	GetAvailableIncludes() []string
	SetAvailableIncludes(includes []string) Transformer
	GetDefaultIncludes() []string
	SetDefaultIncludes(includes []string) Transformer
	GetCurrentScope() *Scope
	SetCurrentScope(scope *Scope) Transformer
	ProcessIncludedResources(scope *Scope, data Any) M
}

// Includer interface
type Includer interface {
	Include(includeName string, data Any, params P) Resource
}

// Serializer interface
type Serializer interface {
	Collection(resourceKey string, data Any) M
	Item(resourceKey string, data Any) M
	Null() M
	IncludeData(resource Resource, data Any) Any
	Meta(meta M) M
	Paginator(paginator Paginator) M
	Cursor(cursor Cursor) M
	MergeIncludes(transformed, included M) M
	SideloadIncludes() bool
	InjectData(data, rawIncluded Any) Any
	InjectAvailableIncludeData(data M, availableIncludes []string) M
	FilterIncludes(included, data Any) Any
}

// ScopeFactory interface
type ScopeFactory interface {
	CreateScopeFor(manager *Manager, resource Resource, opts ...ModScopeOption) *Scope
	CreateChildScopeFor(manager *Manager, parentScope *Scope, resource Resource, opts ...ModScopeOption) *Scope
}
