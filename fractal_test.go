package fractal_test

import (
	"encoding/json"
	"testing"

	"github.com/ibllex/go-fractal"
	"github.com/ibllex/go-fractal/pagination"
	"github.com/stretchr/testify/assert"
)

type User struct {
	ID   int
	Name string
}

type UserTransformer struct {
	fractal.BaseTransformer
}

func (t *UserTransformer) Transform(data fractal.Any) fractal.M {
	result := fractal.M{}

	if u := t.toUser(data); u != nil {
		result["id"] = u.ID
		result["name"] = u.Name
	}

	return result
}

func (t *UserTransformer) toUser(data fractal.Any) *User {

	switch u := data.(type) {
	case *User:
		return u
	case User:
		return &u
	}

	return nil
}

func NewUserTransformer() *UserTransformer {
	t := &UserTransformer{}
	return t
}

type Category struct {
	ID      int
	Name    string
	Creator *User
	Books   []*Book
}

type CategoryTransformer struct {
	fractal.BaseTransformer
	primitive bool
}

func (t *CategoryTransformer) SetPrimitive(primitive bool) *CategoryTransformer {
	t.primitive = primitive
	return t
}

func (t *CategoryTransformer) Transform(data fractal.Any) fractal.M {
	result := fractal.M{}

	if c := t.toCategory(data); c != nil {
		result["id"] = c.ID
		result["name"] = c.Name
	}

	return result
}

func (t *CategoryTransformer) toCategory(data fractal.Any) *Category {

	switch c := data.(type) {
	case *Category:
		return c
	case Category:
		return &c
	}

	return nil
}

func (t *CategoryTransformer) Include(includeName string, data fractal.Any, params fractal.P) fractal.Resource {

	switch includeName {
	case "creator":
		return t.includeUser(data, params)
	case "books":
		return t.includeBooks(data, params)
	}

	return nil
}

func (t *CategoryTransformer) includeBooks(data fractal.Any, params fractal.P) fractal.Resource {

	if c := t.toCategory(data); c != nil {

		books := make([]interface{}, len(c.Books))
		for i := range c.Books {
			books[i] = c.Books[i]
		}

		opts := []fractal.ModResourceOption{
			fractal.WithData(books),
			fractal.WithTransformer(NewBookTransformer()),
		}

		return t.PrimitiveCollection(opts...)
	}

	return nil
}

func (t *CategoryTransformer) includeUser(data fractal.Any, params fractal.P) fractal.Resource {

	if c := t.toCategory(data); c != nil {
		if c.Creator != nil {

			opts := []fractal.ModResourceOption{
				fractal.WithData(c.Creator),
				fractal.WithTransformer(NewUserTransformer()),
			}

			if t.primitive {
				return t.Primitive(opts...)
			}

			return t.Item(opts...)
		}
	}

	return nil
}

func NewCategoryTransformer() *CategoryTransformer {
	t := &CategoryTransformer{}
	t.SetIncluder(t).SetAvailableIncludes([]string{"creator", "books"})
	return t
}

type Book struct {
	ID       int
	Title    string
	Year     int
	Author   string
	Category *Category
}

type BookTransformer struct {
	fractal.BaseTransformer
	primitive bool
}

func (t *BookTransformer) SetPrimitive(primitive bool) *BookTransformer {
	t.primitive = primitive
	return t
}

func (t *BookTransformer) Transform(data fractal.Any) fractal.M {
	result := fractal.M{}

	if b := t.toBook(data); b != nil {
		result["id"] = b.ID
		result["title"] = "'" + b.Title + "'"
		result["year"] = b.Year
		result["author"] = b.Author
	}

	return result
}

func (t *BookTransformer) toBook(data fractal.Any) *Book {

	switch b := data.(type) {
	case *Book:
		return b
	case Book:
		return &b
	}

	return nil
}

func (t *BookTransformer) Include(includeName string, data fractal.Any, params fractal.P) fractal.Resource {

	switch includeName {
	case "category":
		return t.includeCategory(data, params)
	}

	return nil
}

func (t *BookTransformer) includeCategory(data fractal.Any, params fractal.P) fractal.Resource {

	if b := t.toBook(data); b != nil {
		if b.Category != nil {
			opts := []fractal.ModResourceOption{
				fractal.WithData(b.Category),
				fractal.WithTransformer(
					NewCategoryTransformer().SetPrimitive(t.primitive),
				),
			}

			if t.primitive {
				return t.Primitive(opts...)
			}

			return t.Item(opts...)
		}
	}

	return nil
}

func NewBookTransformer() *BookTransformer {
	t := &BookTransformer{}
	t.SetIncluder(t).SetAvailableIncludes([]string{"category"})
	return t
}

func TestItem(t *testing.T) {
	data := Book{1, "Hogfather", 1998, "Philip K Dick", &Category{}}

	manager := fractal.NewManager(nil)
	manager.SetSerializer(&fractal.ArraySerializer{})
	resource := fractal.NewItem(
		fractal.WithData(data),
		fractal.WithResourceKey("item"),
		fractal.WithTransformer(NewBookTransformer()),
	)

	expected := fractal.M{"item": fractal.M{
		"id":     data.ID,
		"title":  "'" + data.Title + "'",
		"year":   data.Year,
		"author": data.Author,
	}}

	expectedJson, _ := json.Marshal(expected)

	t.Run("to map", func(t *testing.T) {

		actual, err := manager.CreateData(resource, nil).ToMap()

		assert.Nil(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("to json", func(t *testing.T) {
		actual, err := manager.CreateData(resource, nil).ToJSON()

		assert.Nil(t, err)
		assert.Equal(t, string(expectedJson), actual)
	})
}

func TestCollection(t *testing.T) {
	books := []fractal.Any{
		Book{1, "Hogfather", 1998, "Philip K Dick", &Category{}},
		Book{2, "Game Of Kill Everyone", 2014, "George R. R. Satan", &Category{}},
	}

	manager := fractal.NewManager(nil)
	resource := fractal.NewCollection(
		fractal.WithData(books),
		fractal.WithTransformer(NewBookTransformer()),
	)

	transformed := []fractal.Any{}
	for _, data := range books {
		b, _ := data.(Book)
		transformed = append(transformed,
			fractal.M{
				"id":     b.ID,
				"title":  "'" + b.Title + "'",
				"year":   b.Year,
				"author": b.Author,
			})
	}

	expected := fractal.M{"data": transformed}
	expectedJson, _ := json.Marshal(expected)

	t.Run("to map", func(t *testing.T) {

		actual, err := manager.CreateData(resource, nil).ToMap()

		assert.Nil(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("to json", func(t *testing.T) {
		actual, err := manager.CreateData(resource, nil).ToJSON()

		assert.Nil(t, err)
		assert.Equal(t, string(expectedJson), actual)
	})
}

func TestPagination(t *testing.T) {
	books := []fractal.Any{
		Book{1, "Hogfather", 1998, "Philip K Dick", &Category{}},
		Book{2, "Game Of Kill Everyone", 2014, "George R. R. Satan", &Category{}},
	}

	page := pagination.NewLengthAwarePaginator(
		books, 10, 2,
		pagination.WithPath("https://www.example.com/books/?user=example"),
		pagination.WithQuery(map[string]string{
			"cat": "1",
		}),
		pagination.WithCurrentPage(2),
	)

	manager := fractal.NewManager(nil)
	resource := fractal.NewCollection(
		fractal.WithData(books),
		fractal.WithTransformer(NewBookTransformer()),
	).SetPaginator(page)

	transformed := []fractal.Any{}
	for _, data := range books {
		b, _ := data.(Book)
		transformed = append(transformed,
			fractal.M{
				"id":     b.ID,
				"title":  "'" + b.Title + "'",
				"year":   b.Year,
				"author": b.Author,
			})
	}

	expected := fractal.M{
		"data": transformed,
		"meta": fractal.M{
			"pagination": fractal.M{
				"total":        uint(10),
				"count":        uint(2),
				"per_page":     uint(2),
				"current_page": uint(2),
				"total_pages":  uint(5),
				"links": map[string]string{
					"previous": "https://www.example.com/books/?user=example&cat=1&page=1",
					"next":     "https://www.example.com/books/?user=example&cat=1&page=3",
				},
			},
		},
	}

	actual, err := manager.CreateData(resource, nil).ToMap()

	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func TestInclude(t *testing.T) {
	cat := &Category{ID: 1, Name: "novel", Creator: &User{ID: 1, Name: "Tamas"}}
	book := Book{1, "Hogfather", 1998, "Philip K Dick", cat}

	manager := fractal.NewManager(nil)
	manager.ParseIncludes([]string{"category", "category.creator"})

	resource := fractal.NewItem(
		fractal.WithData(book),
		fractal.WithTransformer(NewBookTransformer()),
	)

	expected := fractal.M{"data": fractal.M{
		"id":     book.ID,
		"title":  "'" + book.Title + "'",
		"year":   book.Year,
		"author": book.Author,
		"category": fractal.M{
			"data": fractal.M{
				"id":   cat.ID,
				"name": cat.Name,
				"creator": fractal.M{
					"data": fractal.M{
						"id":   cat.Creator.ID,
						"name": cat.Creator.Name,
					},
				},
			},
		},
	}}

	actual, err := manager.CreateData(resource, nil).ToMap()

	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func TestPrimitiveInclude(t *testing.T) {
	cat := &Category{ID: 1, Name: "novel", Creator: &User{ID: 1, Name: "Tamas"}}
	book := Book{1, "Hogfather", 1998, "Philip K Dick", cat}

	manager := fractal.NewManager(nil)
	manager.ParseIncludes([]string{"category", "category.creator"})

	resource := fractal.NewItem(
		fractal.WithData(book),
		fractal.WithTransformer(NewBookTransformer().SetPrimitive(true)),
	)

	expected := fractal.M{"data": fractal.M{
		"id":     book.ID,
		"title":  "'" + book.Title + "'",
		"year":   book.Year,
		"author": book.Author,
		"category": fractal.M{
			"id":   cat.ID,
			"name": cat.Name,
		},
	}}

	actual, err := manager.CreateData(resource, nil).ToMap()

	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func TestPrimitiveCollectionInclude(t *testing.T) {

	books := []*Book{
		{1, "Hogfather", 1998, "Philip K Dick", &Category{}},
		{2, "Game Of Kill Everyone", 2014, "George R. R. Satan", &Category{}},
	}

	cat := &Category{ID: 1, Name: "novel", Books: books}

	manager := fractal.NewManager(nil)
	manager.ParseIncludes([]string{"books"})

	resource := fractal.NewItem(
		fractal.WithData(cat),
		fractal.WithTransformer(NewCategoryTransformer().SetPrimitive(true)),
	)

	expected := fractal.M{"data": fractal.M{
		"id":   cat.ID,
		"name": cat.Name,
		"books": []fractal.Any{
			map[string]interface{}{
				"id":     books[0].ID,
				"title":  "'" + books[0].Title + "'",
				"year":   books[0].Year,
				"author": books[0].Author,
			},
			map[string]interface{}{
				"id":     books[1].ID,
				"title":  "'" + books[1].Title + "'",
				"year":   books[1].Year,
				"author": books[1].Author,
			},
		},
	}}

	actual, err := manager.CreateData(resource, nil).ToMap()

	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}
