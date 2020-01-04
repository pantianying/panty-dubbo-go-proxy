package context

import (
	"context"
)

type ProxyContext interface {
	//base
	NextFilterName() string

	//your
	//InterfaceName() string
}
type baseContext struct {
	filter []string
	ctx    context.Context
}

func NewBaseContext(filter []string) *baseContext {
	return &baseContext{
		filter: filter,
		ctx:    context.Background(),
	}
}
func (h *baseContext) NextFilterName() string {
	if len(h.filter) > 0 {
		f := h.filter[0]
		h.filter = h.filter[1:]
		return f
	}
	return ""
}
