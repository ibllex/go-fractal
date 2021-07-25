package fractal

import (
	"encoding/json"
	"errors"
	"strings"
)

// Scope acts as a tracker, relating a specific resource in a specific
// context. For example, the same resource could be attached to multiple scopes.
// There are root scopes, parent scopes and child scopes.
type Scope struct {
	identifier        string
	manager           *Manager
	resource          Resource
	availableIncludes []string
	parentScopes      []string
}

// ToJSON convert the current data for this scope to json.
func (s *Scope) ToJSON() (string, error) {
	m, err := s.ToMap()
	if err != nil {
		return "", err
	}

	str, err := json.Marshal(m)
	return string(str), err
}

// ToMap convert the current data for this scope to a map.
func (s *Scope) ToMap() (M, error) {

	rawData, _, err := s.executeResourceTransformers()
	if err != nil {
		return nil, err
	}

	serializer := s.manager.GetSerializer()
	data := s.serializeResource(serializer, rawData)

	// If the serializer wants the includes to be side-loaded then we'll
	// serialize the included data and merge it with the data.
	// if serializer.SideloadIncludes() {
	// 	//
	// }

	if len(s.availableIncludes) > 0 {
		data = serializer.InjectAvailableIncludeData(data, s.availableIncludes)
	}

	if c, ok := s.resource.(*Collection); ok {
		var pagination M

		if c.HasPaginator() {
			pagination = serializer.Paginator(c.GetPaginator())
		}

		for k, v := range pagination {
			s.resource.SetMetaValue(k, v)
		}
	}

	meta := serializer.Meta(s.resource.GetMeta())
	if data == nil {
		if len(meta) != 0 {
			return meta, nil
		}

		return nil, nil
	}

	for k, v := range meta {
		data[k] = v
	}

	return data, nil
}

// TransformPrimitiveResource transformer a primitive resource
func (s *Scope) TransformPrimitiveResource() (Any, error) {
	transformer := s.resource.GetTransformer()
	transformer.SetCurrentScope(s)
	data := s.resource.GetData()

	switch s.resource.(type) {
	case *Primitive:
		data = transformer.Transform(data)
		return data, nil
	case *PrimitiveCollection:
		transformedData := []Any{}
		anyCollection, ok := data.([]Any)

		if !ok {
			return nil, errors.New(
				"the data of primitive collection resource should be []interface{} or []Any",
			)
		}

		for _, d := range anyCollection {
			transformedData = append(transformedData, transformer.Transform(d))
		}

		return transformedData, nil
	}

	return nil, errors.New("argument resource should be an instance of fractal.Primitive or fractal.PrimitiveCollection")
}

func (s *Scope) executeResourceTransformers() (Any, Any, error) {
	transformer := s.resource.GetTransformer()
	data := s.resource.GetData()

	switch s.resource.(type) {
	case *Item:
		transformedData, includedData := s.fireTransformer(transformer, data)
		return transformedData, includedData, nil
	case *Collection:
		transformedData := []Any{}
		includedData := []Any{}
		anyCollection, ok := data.([]Any)

		if !ok {
			return nil, nil, errors.New(
				"the data of collection resource should be []interface{} or []Any",
			)
		}

		for _, d := range anyCollection {
			transformed, included := s.fireTransformer(transformer, d)
			transformedData = append(transformedData, transformed)
			includedData = append(includedData, included)
		}

		return transformedData, includedData, nil
	}

	return nil, nil, errors.New(
		"argument resource should be an instance of fractal.Item or fractal.Collection",
	)
}

func (s *Scope) serializeResource(serializer Serializer, data Any) M {
	resourceKey := s.resource.GetResourceKey()

	switch s.resource.(type) {
	case *Item:
		return serializer.Item(resourceKey, data)
	case *Collection:
		return serializer.Collection(resourceKey, data)
	}

	return serializer.Null()
}

func (s *Scope) fireTransformer(transformer Transformer, data Any) (M, M) {
	var includedData M

	transformer.SetCurrentScope(s)
	transformedData := transformer.Transform(data)

	if s.transformerHasIncludes(transformer) {
		includedData = s.fireIncludedTransformers(transformer, data)
		transformedData = s.manager.GetSerializer().MergeIncludes(transformedData, includedData)
	}

	// Stick only with requested fields
	transformedData = s.filterFieldsets(transformedData)
	return transformedData, includedData
}

func (s *Scope) transformerHasIncludes(transformer Transformer) bool {

	defaultIncludes := transformer.GetDefaultIncludes()
	availableIncludes := transformer.GetAvailableIncludes()

	return len(defaultIncludes) != 0 || len(availableIncludes) != 0
}

func (s *Scope) fireIncludedTransformers(transformer Transformer, data Any) M {
	s.availableIncludes = transformer.GetAvailableIncludes()
	return transformer.ProcessIncludedResources(s, data)
}

// EmbedChildScope embed a scope as a child of the current scope
func (s *Scope) EmbedChildScope(identifier string, resource Resource) *Scope {
	return s.manager.CreateData(resource, s, WithIdentifier(identifier))
}

// Filter the provided data with the requested filter fieldset for
// the scope resource
func (s *Scope) filterFieldsets(data M) M {
	if !s.hasFilterFieldset() {
		return data
	}

	// serializer := s.manager.GetSerializer()
	// requestedFieldset := s.getFilterFieldset()

	return data
}

// GetScopeIdentifier get the current identifier
func (s *Scope) GetScopeIdentifier() string {
	return s.identifier
}

// GetIdentifier get the unique identifier for this scope
func (s *Scope) GetIdentifier(appendIdentifier string) string {

	identifierParts := s.parentScopes
	identifierParts = append(identifierParts, s.identifier)

	if appendIdentifier != "" {
		identifierParts = append(identifierParts, appendIdentifier)
	}

	return strings.Join(identifierParts, ".")
}

// GetParentScopes getter for parent scopes
func (s *Scope) GetParentScopes() []string {
	return s.parentScopes
}

// SetParentScopes setter for parent scopes
func (s *Scope) SetParentScopes(scopes []string) *Scope {
	s.parentScopes = scopes
	return s
}

// GetManager getter for manager
func (s *Scope) GetManager() *Manager {
	return s.manager
}

// Return the requested filter fieldset for the scope resource
func (s *Scope) getFilterFieldset() []string {
	return s.manager.GetFieldset(s.getResourceType())
}

func (s *Scope) hasFilterFieldset() bool {
	return len(s.getFilterFieldset()) != 0
}

func (s *Scope) getResourceType() string {
	return s.resource.GetResourceKey()
}

// GetResource getter for resource
func (s *Scope) GetResource() Resource {
	return s.resource
}

// IsRequested  Check if - in relation to the current scope - this specific segment is allowed.
// That means, if a.b.c is requested and the current scope is a.b, then c is allowed. If the current
// scope is a then c is not allowed, even if it is there and potentially transformable.
func (s *Scope) IsRequested(checkScopeSegment string) bool {

	scopeArray := []string{checkScopeSegment}

	if len(s.parentScopes) > 0 {
		scopeArray = s.parentScopes[1:]
		scopeArray = append(scopeArray, s.identifier, checkScopeSegment)
	}

	scopeString := strings.Join(scopeArray, ".")

	for _, exclude := range s.manager.GetRequestedIncludes() {
		if exclude == scopeString {
			return true
		}
	}

	return false
}

// IsExcluded Check if - in relation to the current scope - this specific segment should
// be excluded. That means, if a.b.c is excluded and the current scope is a.b,
// then c will not be allowed in the transformation whether it appears in
// the list of default or available, requested includes.
func (s *Scope) IsExcluded(checkScopeSegment string) bool {

	scopeArray := []string{checkScopeSegment}

	if len(s.parentScopes) > 0 {
		scopeArray = s.parentScopes[1:]
		scopeArray = append(scopeArray, s.identifier, checkScopeSegment)
	}

	scopeString := strings.Join(scopeArray, ".")

	for _, exclude := range s.manager.GetRequestedExcludes() {
		if exclude == scopeString {
			return true
		}
	}

	return false
}

// ScopeOption options for scope object
type ScopeOption struct {
	Identifier string
}

// ModScopeOption function to modify scope option
type ModScopeOption func(option *ScopeOption)

// WithIdentifier is an easy way to set identifier for scope
func WithIdentifier(identifier string) ModScopeOption {
	return func(option *ScopeOption) {
		option.Identifier = identifier
	}
}

// NewScope create new scope
func NewScope(manager *Manager, resource Resource, opts ...ModScopeOption) *Scope {
	opt := ScopeOption{}

	for _, mod := range opts {
		mod(&opt)
	}

	return &Scope{
		manager:    manager,
		resource:   resource,
		identifier: opt.Identifier,
	}
}
