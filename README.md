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

## License

This library is under the [MIT](https://github.com/ibllex/go-fractal/blob/main/LICENSE) license.
