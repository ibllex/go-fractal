package gin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ibllex/go-fractal"
)

// Context gin context with fractal extension
type Context struct {
	*gin.Context
	manager *fractal.Manager
}

func (c *Context) NoContent() {
	c.Status(http.StatusNoContent)
}

func (c *Context) getResponse(resource fractal.Resource, status int, callbacks ...Callback) *Response {
	rsp := &Response{c, c.manager, resource, status}

	for _, callback := range callbacks {
		callback(rsp)
	}

	return rsp
}

func (c *Context) renderResource(resource fractal.Resource, callbacks ...Callback) {

	rsp := c.getResponse(resource, http.StatusOK, callbacks...)

	data, err := c.manager.CreateData(
		resource, nil,
		fractal.WithIdentifier(resource.GetResourceKey()),
	).ToMap()

	if err != nil {
		c.ErrorInternal(WithMessage(err.Error()))
	} else {
		c.JSON(rsp.Status, data)
	}
}

func (c *Context) Collection(items []fractal.Any, transformer fractal.Transformer, callbacks ...Callback) {
	resource := fractal.NewCollection(
		fractal.WithData(items),
		fractal.WithTransformer(transformer),
	)

	c.renderResource(resource)
}

func (c *Context) Item(item fractal.Any, transformer fractal.Transformer, callbacks ...Callback) {
	resource := fractal.NewItem(
		fractal.WithData(item),
		fractal.WithTransformer(transformer),
	)

	c.renderResource(resource)
}

func (c *Context) Paginator(paginator Paginator, transformer fractal.Transformer, callbacks ...Callback) {
	resource := fractal.NewCollection(
		fractal.WithData(paginator.GetItems()),
		fractal.WithTransformer(transformer),
	).SetPaginator(paginator)

	c.renderResource(resource)
}

func (c *Context) getErrorOption(opt *ErrorOption, mods ...ModErrorOption) *ErrorOption {
	for _, mod := range mods {
		mod(opt)
	}
	return opt
}

// Error return an error
func (c *Context) Error(opt *ErrorOption, status int) {
	result := map[string]interface{}{
		"message":     opt.Message,
		"status_code": status,
	}

	if len(opt.Errors) > 0 {
		result["errors"] = opt.Errors
	}

	c.JSON(status, result)
}

// ErrorNotFound return a 404 error
func (c *Context) ErrorNotFound(mods ...ModErrorOption) {
	opt := &ErrorOption{Message: "Not Found"}
	c.Error(c.getErrorOption(opt, mods...), http.StatusNotFound)
}

// AbortNotFound return a 404 error and abort
func (c *Context) AbortNotFound(mods ...ModErrorOption) {
	c.ErrorNotFound(mods...)
	c.Abort()
}

// ErrorBadRequest return a 400 error
func (c *Context) ErrorBadRequest(mods ...ModErrorOption) {
	opt := &ErrorOption{Message: "Bad Request"}
	c.Error(c.getErrorOption(opt, mods...), http.StatusBadRequest)
}

// AbortBadRequest return a 400 error and abort
func (c *Context) AbortBadRequest(mods ...ModErrorOption) {
	c.ErrorBadRequest(mods...)
	c.Abort()
}

// ErrorForbidden return a 403 error
func (c *Context) ErrorForbidden(mods ...ModErrorOption) {
	opt := &ErrorOption{Message: "Forbidden"}
	c.Error(c.getErrorOption(opt, mods...), http.StatusForbidden)
}

// AbortForbidden return a 403 error and abort
func (c *Context) AbortForbidden(mods ...ModErrorOption) {
	c.ErrorForbidden(mods...)
	c.Abort()
}

// ErrorInternal return a 500 error
func (c *Context) ErrorInternal(mods ...ModErrorOption) {
	opt := &ErrorOption{Message: "Internal Error"}
	c.Error(c.getErrorOption(opt, mods...), http.StatusInternalServerError)
}

// AbortInternal return a 500 error and abort
func (c *Context) AbortInternal(mods ...ModErrorOption) {
	c.ErrorInternal(mods...)
	c.Abort()
}

// ErrorUnauthorized return a 401 error
func (c *Context) ErrorUnauthorized(mods ...ModErrorOption) {
	opt := &ErrorOption{Message: "Unauthorized"}
	c.Error(c.getErrorOption(opt, mods...), http.StatusUnauthorized)
}

// AbortUnauthorized return a 401 error and abort
func (c *Context) AbortUnauthorized(mods ...ModErrorOption) {
	c.ErrorUnauthorized(mods...)
	c.Abort()
}

// ErrorMethodNotAllowed return a 405 error
func (c *Context) ErrorMethodNotAllowed(mods ...ModErrorOption) {
	opt := &ErrorOption{Message: "Method Not Allowed"}
	c.Error(c.getErrorOption(opt, mods...), http.StatusMethodNotAllowed)
}

// AbortMethodNotAllowed return a 405 error and abort
func (c *Context) AbortMethodNotAllowed(mods ...ModErrorOption) {
	c.ErrorMethodNotAllowed(mods...)
	c.Abort()
}

// ErrorUnprocessable return a 422 error
func (c *Context) ErrorUnprocessable(mods ...ModErrorOption) {
	opt := &ErrorOption{Message: "Unprocessable Entity"}
	c.Error(c.getErrorOption(opt, mods...), http.StatusUnprocessableEntity)
}

// AbortUnprocessable return a 422 error and abort
func (c *Context) AbortUnprocessable(mods ...ModErrorOption) {
	c.ErrorUnprocessable(mods...)
	c.Abort()
}
