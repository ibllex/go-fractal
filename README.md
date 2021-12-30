# Go-Fractal

go-fractal is the go version of [league/fractal](https://github.com/thephpleague/fractal) which is is written in PHP, fractal provides a presentation and transformation layer for complex data output, the like found in RESTful APIs, and works really well with JSON. Think of this as a view layer for your JSON/YAML/etc.

## Install

```bash
go get -u github.com/ibllex/go-fractal
```

## Quick usage

> For the sake of simplicity, this example has been put together as though it was one file. In reality you would spread the manager initiation, data collection and JSON conversion into separate parts of your application.

### Transform an item

```go
package main

import (
	"fmt"

	"github.com/ibllex/go-fractal"
)

type Book struct {
	ID     int
	Title  string
	Year   int
	Author string
}

func main() {

	// This is the data to be converted
	data := &Book{1, "Hogfather", 1998, "Philip K Dick"}

	// Create a top level instance
	manager := fractal.NewManager(nil)
	manager.SetSerializer(&fractal.ArraySerializer{})

	// This "Transformer" can be a callback or a new instance of a fractal.Transformer object
	transformer := fractal.T(func(t *fractal.BaseTransformer, data fractal.Any) fractal.M {
		result := fractal.M{}

		if b, ok := data.(*Book); ok {
			result["id"] = b.ID
			result["title"] = "'" + b.Title + "'"
			result["year"] = b.Year
			result["author"] = b.Author
		}

		return result
	})

	// Convert the single book into an item resource
	resource := fractal.NewItem(
		fractal.WithData(data),
		fractal.WithTransformer(transformer),
	)

	// Turn all of that into a JSON string
	json, _ := manager.CreateData(resource, nil).ToJSON()

	fmt.Println(json)
	// Outputs: {"data":{"author":"Philip K Dick","id":1,"title":"'Hogfather'","year":1998}}
}
```

### Transform a collection of data

```go
package main

import (
	"fmt"

	"github.com/ibllex/go-fractal"
)

type Book struct {
	ID     int
	Title  string
	Year   int
	Author string
}

func main() {

	// This is the data to be converted
	books := []fractal.Any{
		&Book{1, "Hogfather", 1998, "Philip K Dick"},
		&Book{2, "Game Of Kill Everyone", 2014, "George R. R. Satan"},
	}

	// Create a top level instance
	manager := fractal.NewManager(nil)
	manager.SetSerializer(&fractal.ArraySerializer{})

	// This "Transformer" can be a callback or a new instance of a fractal.Transformer object
	transformer := fractal.T(func(t *fractal.BaseTransformer, data fractal.Any) fractal.M {
		result := fractal.M{}

		if b, ok := data.(*Book); ok {
			result["id"] = b.ID
			result["title"] = "'" + b.Title + "'"
			result["year"] = b.Year
			result["author"] = b.Author
		}

		return result
	})

	// Convert the books into a collection resource
	resource := fractal.NewCollection(
		fractal.WithData(books),
		fractal.WithTransformer(transformer),
	)

	// Turn all of that into a JSON string
	json, _ := manager.CreateData(resource, nil).ToJSON()

	fmt.Println(json)
	// Outputs: {"data":[{"author":"Philip K Dick","id":1,"title":"'Hogfather'","year":1998},{"author":"George R. R. Satan","id":2,"title":"'Game Of Kill Everyone'","year":2014}]}
}
```

### Transform with pagination

```go
package main

import (
	"fmt"

	"github.com/ibllex/go-fractal"
	"github.com/ibllex/go-fractal/pagination"
)

type Book struct {
	ID     int
	Title  string
	Year   int
	Author string
}

func main() {

	// This is the data to be converted
	books := []fractal.Any{
		&Book{1, "Hogfather", 1998, "Philip K Dick"},
		&Book{2, "Game Of Kill Everyone", 2014, "George R. R. Satan"},
	}

	// Create an paginator
	paginator := pagination.NewLengthAwarePaginator(
		books, 10, 2,
		pagination.WithPath("https://www.example.com/books/?user=example"),
		pagination.WithCurrentPage(2),
	)

	// Create a top level instance
	manager := fractal.NewManager(nil)
	manager.SetSerializer(&fractal.ArraySerializer{})

	// This "Transformer" can be a callback or a new instance of a fractal.Transformer object
	transformer := fractal.T(func(t *fractal.BaseTransformer, data fractal.Any) fractal.M {
		result := fractal.M{}

		if b, ok := data.(*Book); ok {
			result["id"] = b.ID
			result["title"] = "'" + b.Title + "'"
			result["year"] = b.Year
			result["author"] = b.Author
		}

		return result
	})

	// Convert the books into a collection resource with a paginator
	resource := fractal.NewCollection(
		fractal.WithData(books),
		fractal.WithTransformer(transformer),
	).SetPaginator(paginator)

	// Turn all of that into a JSON string
	json, _ := manager.CreateData(resource, nil).ToJSON()

	fmt.Println(json)
	// Outputs: {"data":[{"author":"Philip K Dick","id":1,"title":"'Hogfather'","year":1998},{"author":"George R. R. Satan","id":2,"title":"'Game Of Kill Everyone'","year":2014}],"meta":{"pagination":{"count":2,"current_page":2,"links":{"next":"https://www.example.com/books/?user=example&page=3","previous":"https://www.example.com/books/?user=example&page=1"},"per_page":2,"total":10,"total_pages":5}}}
}
```

### Use with gin

```go
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ibllex/go-fractal"
	fgin "github.com/ibllex/go-fractal/gin"
	"github.com/ibllex/go-fractal/pagination"
)

type Book struct {
	ID     int
	Title  string
	Year   int
	Author string
}

type BookTransformer struct {
	*fractal.BaseTransformer
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

func NewBookTransformer() *BookTransformer {
	return &BookTransformer{&fractal.BaseTransformer{}}
}

func GetBook(c *fgin.Context) {
	// This is the data to be converted
	data := &Book{1, "Hogfather", 1998, "Philip K Dick"}

	// Turn all of that into a JSON output
	c.Item(data, NewBookTransformer())
	// Outputs: {"data":{"author":"Philip K Dick","id":1,"title":"'Hogfather'","year":1998}}
}

func GetBooks(c *fgin.Context) {
	// This is the data to be converted
	books := []fractal.Any{
		&Book{1, "Hogfather", 1998, "Philip K Dick"},
		&Book{2, "Game Of Kill Everyone", 2014, "George R. R. Satan"},
	}

	// Create an paginator
	paginator := pagination.NewLengthAwarePaginator(
		books, 10, 2,
		pagination.WithPath("https://www.example.com/books/?user=example"),
		pagination.WithCurrentPage(2),
	)

	// Turn all of that into a JSON output with pagination
	c.Paginator(paginator, NewBookTransformer())
	// Outputs: {"data":[{"author":"Philip K Dick","id":1,"title":"'Hogfather'","year":1998},{"author":"George R. R. Satan","id":2,"title":"'Game Of Kill Everyone'","year":2014}],"meta":{"pagination":{"count":2,"current_page":2,"links":{"next":"https://www.example.com/books/?user=example\u0026page=3","previous":"https://www.example.com/books/?user=example\u0026page=1"},"per_page":2,"total":10,"total_pages":5}}}

	// You can also use c.Collection() to transform a collection without pagination
	// c.Collection(books, NewBookTransformer())
	// Outputs: {"data":[{"author":"Philip K Dick","id":1,"title":"'Hogfather'","year":1998},{"author":"George R. R. Satan","id":2,"title":"'Game Of Kill Everyone'","year":2014}]}
}

func NotFound(c *fgin.Context) {
	// A convenient way to return a 404 status,
	// We also provide many other methods to return error messages easily,
	// For more infomation: https://github.com/ibllex/go-fractal/blob/main/gin/context.go
	c.ErrorNotFound()
	// Outputs: {"message":"Not Found","status_code":404}
}

func main() {
	r := gin.Default()

	r.GET("/books", fgin.H(GetBooks))
	r.GET("/book", fgin.H(GetBook))
	r.GET("/404", fgin.H(NotFound))

	srv := http.Server{
		Addr:    fmt.Sprintf(":%d", 8080),
		Handler: r,
	}

	log.Printf("HTTP Service Started at %s", srv.Addr)
	srv.ListenAndServe()
}
```

## License

This library is under the [MIT](https://github.com/ibllex/go-fractal/blob/main/LICENSE) license.
