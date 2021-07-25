package fractal

// DataArraySerializer array serializer with default resource key
type DataArraySerializer struct {
	ArraySerializer
}

// Collection serialize a collection
func (s *DataArraySerializer) Collection(resourceKey string, data Any) M {
	return s.ArraySerializer.Collection("", data)
}

// Item serialize an item
func (s *DataArraySerializer) Item(resourceKey string, data Any) M {
	return M{
		DefaultResourceKey: data,
	}
}

// Null serialize null resource
func (s *DataArraySerializer) Null() M {
	return M{
		DefaultResourceKey: nil,
	}
}
