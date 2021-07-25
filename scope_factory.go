package fractal

// DefaultScopeFactory default scope factory
type DefaultScopeFactory struct {
	//
}

// CreateScopeFor create new scope
func (f *DefaultScopeFactory) CreateScopeFor(manager *Manager, resource Resource, opts ...ModScopeOption) *Scope {
	return NewScope(manager, resource, opts...)
}

// CreateChildScopeFor create new child scope by given parent
func (f *DefaultScopeFactory) CreateChildScopeFor(manager *Manager, parentScope *Scope, resource Resource, opts ...ModScopeOption) *Scope {

	scope := f.CreateScopeFor(manager, resource, opts...)

	// This will be the new children list of parents (parents parents, plus the parent)
	scopeArray := parentScope.GetParentScopes()
	scopeArray = append(scopeArray, parentScope.GetScopeIdentifier())

	scope.SetParentScopes(scopeArray)
	return scope
}
