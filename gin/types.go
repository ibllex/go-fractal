package gin

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/ibllex/go-fractal"
)

// Paginator paginator with items
type Paginator interface {
	fractal.Paginator
	GetItems() []interface{}
	SetItems(items []interface{})
}

// Request request with binging
type Request interface {
	Messages() map[string]string
}

// HandlerFunc fractal gin handler
type HandlerFunc func(*Context)

// Callback modify response
type Callback func(*Response)

// H fractal handler wrapper for gin defaulr handler
func H(h HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		manager := fractal.NewManager(nil)
		manager.SetSerializer(&fractal.ArraySerializer{})
		manager.ParseIncludes(strings.Split(c.Query("include"), ","))
		ctx := &Context{c, manager}
		h(ctx)
	}
}

// ErrorOption options for error
type ErrorOption struct {
	Message string
	Errors  map[string][]string
}

// ModErrorOption function to modify error option
type ModErrorOption func(*ErrorOption)

// WithMessage change error message
func WithMessage(message string) ModErrorOption {
	return func(opt *ErrorOption) {
		opt.Message = message
	}
}

// WithError change error message by error
func WithError(err error) ModErrorOption {
	return func(opt *ErrorOption) {
		opt.Message = err.Error()
	}
}

// WithError change error message by error
func WithValidatorError(err error, req Request) ModErrorOption {
	errors := map[string][]string{}
	if errs, ok := err.(validator.ValidationErrors); ok {
		messages := map[string]string{}
		if req != nil {
			messages = req.Messages()
		}

		for _, e := range errs {
			f := strings.ToLower(strings.Join([]string{e.Field(), e.Tag()}, "."))
			if msg, ok := messages[f]; ok {
				errors[e.Field()] = append(errors[f], msg)
			} else {
				errors[e.Field()] = append(errors[f], "validation for '"+f+"' failed")
			}
		}
	}

	return func(opt *ErrorOption) {
		opt.Message = "The given data was invalid."
		opt.Errors = errors
	}
}
