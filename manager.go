package fractal

import "strings"

// Manager allows users to create the "root scope" easily
type Manager struct {
	// Factory used to create new configured scopes.
	scopeFactory       ScopeFactory
	serializer         Serializer
	requestedFieldsets map[string][]string
	requestedIncludes  []string
	requestedExcludes  []string
	includeParams      map[string]P
	// Upper limit to how many levels of included data are allowed.
	recursionLimit int
}

// CreateData is main method to kick this all off.
// Make a resource then pass it over, and use toJson()
func (m *Manager) CreateData(resource Resource, parentScope *Scope, opts ...ModScopeOption) *Scope {
	if parentScope != nil {
		return m.scopeFactory.CreateChildScopeFor(m, parentScope, resource, opts...)
	}
	return m.scopeFactory.CreateScopeFor(m, resource, opts...)
}

// ParseIncludes parse include params
func (m *Manager) ParseIncludes(includes []string) *Manager {
	// Wipe these before we go again
	m.requestedIncludes = []string{}
	// m.includeParams = map[string][]string{}

	var includeName, allModifiersStr string

	for _, include := range includes {
		includeName, allModifiersStr = m.explodeInclude(include, ":")
		allModifiersStr, _ = m.explodeInclude(allModifiersStr, ".")

		// Trim it down to a cool level of recursion
		includeName = m.trimToAcceptableRecursionLevel(includeName)

		if m.hasRequestInclude(includeName) {
			continue
		}
		m.requestedIncludes = append(m.requestedIncludes, includeName)

		// No params
		if allModifiersStr == "" {
			continue
		}

		// TODO: parse include params
	}

	m.autoIncludeParents()
	return m
}

func (m *Manager) explodeInclude(include string, sep string) (string, string) {
	results := strings.SplitN(include, sep, 2)
	if len(results) < 2 {
		return results[0], ""
	}
	return results[0], results[1]
}

// Auto-include Parents
// Look at the requested includes and automatically include the parents if they
// are not explicitly requested. E.g: [foo, bar.baz] becomes [foo, bar, bar.baz]
func (m *Manager) autoIncludeParents() {

	parsed := []string{}

	for _, include := range m.requestedIncludes {
		nested := strings.Split(include, ".")
		part := nested[0]
		parsed = append(parsed, part)
		nested = nested[1:]

		for len(nested) > 0 {
			part += "." + nested[0]
			parsed = append(parsed, part)
			nested = nested[1:]
		}
	}

	m.requestedIncludes = parsed
}

// Trim to Acceptable Recursion Level
func (m *Manager) trimToAcceptableRecursionLevel(includeName string) string {

	levels := strings.Split(includeName, ".")

	if len(levels) > m.recursionLimit {
		levels = levels[:m.recursionLimit]
	}

	return strings.Join(levels, ".")
}

// SetRecursionLimit set recursion limit
func (m *Manager) SetRecursionLimit(limit int) *Manager {
	m.recursionLimit = limit
	return m
}

// GetSerializer get data serializer and
// return DataArraySerializer if no serializer set
func (m *Manager) GetSerializer() Serializer {
	if m.serializer == nil {
		m.SetSerializer(&DataArraySerializer{})
	}
	return m.serializer
}

// SetSerializer set data serializer
func (m *Manager) SetSerializer(serializer Serializer) *Manager {
	m.serializer = serializer
	return m
}

// GetRequestedFieldsets get requested fieldsets
func (m *Manager) GetRequestedFieldsets() map[string][]string {
	return m.requestedFieldsets
}

// GetFieldset Get fieldset params for the specified type
func (m *Manager) GetFieldset(fieldType string) (fieldset []string) {
	if m.requestedFieldsets == nil {
		return
	}

	if v, ok := m.requestedFieldsets[fieldType]; ok {
		return v
	}

	return
}

// GetRequestedIncludes get requested includes
func (m *Manager) GetRequestedIncludes() []string {
	return m.requestedIncludes
}

func (m *Manager) hasRequestInclude(include string) bool {
	for _, i := range m.requestedIncludes {
		if i == include {
			return true
		}
	}
	return false
}

// GetIncludeParams get include params
func (m *Manager) GetIncludeParams(identifier string) P {
	if m.includeParams == nil {
		return nil
	}

	if p, ok := m.includeParams[identifier]; ok {
		return p
	}

	return nil
}

// GetRequestedExcludes get requested excludes
func (m *Manager) GetRequestedExcludes() []string {
	return m.requestedExcludes
}

// NewManager create new manager
func NewManager(scopeFactory ScopeFactory) *Manager {
	if scopeFactory == nil {
		scopeFactory = &DefaultScopeFactory{}
	}

	return &Manager{
		scopeFactory:   scopeFactory,
		recursionLimit: 10,
	}
}
