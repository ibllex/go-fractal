package fractal

// ArraySerializer array serializer
type ArraySerializer struct {
	//
}

// Collection serialize a collection
func (s *ArraySerializer) Collection(resourceKey string, data Any) M {
	if resourceKey == "" {
		resourceKey = DefaultResourceKey
	}

	return M{
		resourceKey: data,
	}
}

// Item serialize an item
func (s *ArraySerializer) Item(resourceKey string, data Any) M {
	if resourceKey == "" {
		resourceKey = DefaultResourceKey
	}

	return M{
		resourceKey: data,
	}
}

// Null serialize null resource
func (s *ArraySerializer) Null() M {
	return nil
}

// IncludeData serialize include resource
func (s *ArraySerializer) IncludeData(resource Resource, data Any) Any {
	return data
}

// Meta serialize the meta data
func (s *ArraySerializer) Meta(meta M) M {
	if len(meta) == 0 {
		return nil
	}

	return M{
		"meta": meta,
	}
}

// Paginator serialize the paginator
func (s *ArraySerializer) Paginator(paginator Paginator) M {
	currentPage := paginator.GetCurrentPage()
	lastPage := paginator.GetLastPage()

	pagination := M{
		"total":        paginator.GetTotal(),
		"count":        paginator.GetCount(),
		"per_page":     paginator.GetPerPage(),
		"current_page": currentPage,
		"total_pages":  lastPage,
	}

	links := map[string]string{}

	if currentPage > 1 {
		links["previous"] = paginator.GetURL(currentPage - 1)
	}

	if currentPage < lastPage {
		links["next"] = paginator.GetURL(currentPage + 1)
	}

	pagination["links"] = links

	return M{
		"pagination": pagination,
	}
}

// Cursor serialize the cursor
func (s *ArraySerializer) Cursor(cursor Cursor) M {
	data := M{
		"current": cursor.GetCurrent(),
		"prev":    cursor.GetPrev(),
		"next":    cursor.GetNext(),
		"count":   cursor.GetCount(),
	}

	return M{
		"cursor": data,
	}
}

// MergeIncludes merge include data with transformed data
func (s *ArraySerializer) MergeIncludes(transformed, included M) M {

	// If the serializer does not want the includes to be side-loaded then
	// the included data must be merged with the transformed data.
	if !s.SideloadIncludes() {
		for k, v := range included {
			transformed[k] = v
		}
	}

	return transformed
}

// SideloadIncludes indicates if includes should be side-loaded.
func (s *ArraySerializer) SideloadIncludes() bool {
	return false
}

// InjectData is a hook for the serializer to inject custom data based on the relationships of the resource
func (s *ArraySerializer) InjectData(data, rawIncluded Any) Any {
	return data
}

// InjectAvailableIncludeData is a hook for the serializer to inject custom data based on the available includes of the resource
func (s *ArraySerializer) InjectAvailableIncludeData(data M, availableIncludes []string) M {
	return data
}

// FilterIncludes is hook for the serializer to modify the final list of includes.
func (s *ArraySerializer) FilterIncludes(included, data Any) Any {
	return included
}
