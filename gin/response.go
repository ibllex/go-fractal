package gin

import "github.com/ibllex/fractal"

type Response struct {
	Ctx      *Context
	Fractal  *fractal.Manager
	Resource fractal.Resource
	Status   int
}
