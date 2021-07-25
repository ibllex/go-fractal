package fractal

// BaseTransformer All Transformer classes should extend this to utilize the convenience methods
// collection() and item(), and make the availableIncludes property available.
// Extend it and add a `Transform()` method to transform any default or included data
// into a basic array.
type BaseTransformer struct {
	// Resources that can be included if requested.
	availableIncludes []string
	// Include resources without needing it to be requested.
	defaultIncludes []string
	// The transformer should know about the current scope, so we can fetch relevant params
	currentScope *Scope
	// The transformer should have an includer to perform custom includes
	includer Includer
}

// Transform perform transform
func (t *BaseTransformer) Transform(data Any) M {
	if m, ok := data.(M); ok {
		return m
	}
	return M{}
}

// GetAvailableIncludes getter for availableIncludes
func (t *BaseTransformer) GetAvailableIncludes() []string {
	return t.availableIncludes
}

// GetDefaultIncludes getter for defaultIncludes
func (t *BaseTransformer) GetDefaultIncludes() []string {
	return t.defaultIncludes
}

// GetCurrentScope getter for current scope
func (t *BaseTransformer) GetCurrentScope() *Scope {
	return t.currentScope
}

// SetAvailableIncludes setter for availableIncludes
func (t *BaseTransformer) SetAvailableIncludes(includes []string) Transformer {
	t.availableIncludes = includes
	return t
}

// SetDefaultIncludes setter for default includes
func (t *BaseTransformer) SetDefaultIncludes(includes []string) Transformer {
	t.defaultIncludes = includes
	return t
}

// SetCurrentScope setter for current scope
func (t *BaseTransformer) SetCurrentScope(scope *Scope) Transformer {
	t.currentScope = scope
	return t
}

// SetIncluder setter for includer
func (t *BaseTransformer) SetIncluder(includer Includer) *BaseTransformer {
	t.includer = includer
	return t
}

// ProcessIncludedResources is fired to loop through available includes,
// see if any of them are requested and permitted for this scope.
func (t *BaseTransformer) ProcessIncludedResources(scope *Scope, data Any) M {
	includedData := M{}

	includes := t.figureOutWhichIncludes(scope)

	for _, include := range includes {
		includedData = t.includeResourceIfAvailable(
			scope, data, includedData, include,
		)
	}

	return includedData
}

// Include a resource only if it is available on the method
func (t *BaseTransformer) includeResourceIfAvailable(scope *Scope, data Any, includeData M, include string) M {

	resource := t.callIncludeMethod(scope, include, data)

	if resource != nil {
		childScope := scope.EmbedChildScope(include, resource)

		if _, ok := childScope.GetResource().(*Primitive); ok {
			includeData[include], _ = childScope.TransformPrimitiveResource()
		} else if _, ok := childScope.GetResource().(*PrimitiveCollection); ok {
			includeData[include], _ = childScope.TransformPrimitiveResource()
		} else {
			includeData[include], _ = childScope.ToMap()
		}
	}

	return includeData
}

// Call Include Method
func (t *BaseTransformer) callIncludeMethod(scope *Scope, includeName string, data Any) Resource {

	if t.includer == nil {
		return nil
	}

	scopeIdentifier := scope.GetIdentifier(includeName)
	params := scope.GetManager().GetIncludeParams(scopeIdentifier)

	return t.includer.Include(includeName, data, params)
}

// Figure out which includes we need
func (t *BaseTransformer) figureOutWhichIncludes(scope *Scope) []string {
	includes := t.GetDefaultIncludes()

	for _, include := range t.GetAvailableIncludes() {
		if scope.IsRequested(include) {
			includes = append(includes, include)
		}
	}

	target := includes[:0]

	for _, include := range includes {
		if !scope.IsExcluded(include) {
			target = append(target, include)
		}
	}

	return target
}

// Item create a new item resource object.
func (t *BaseTransformer) Item(opts ...ModResourceOption) *Item {
	return NewItem(opts...)
}

// Collection create a new collection resource object.
func (t *BaseTransformer) Collection(opts ...ModResourceOption) *Collection {
	return NewCollection(opts...)
}

// Primitive create a new primitive resource object.
func (t *BaseTransformer) Primitive(opts ...ModResourceOption) *Primitive {
	return NewPrimitive(opts...)
}

// Primitive create a new primitive resource object.
func (t *BaseTransformer) PrimitiveCollection(opts ...ModResourceOption) *PrimitiveCollection {
	return NewPrimitiveCollection(opts...)
}
