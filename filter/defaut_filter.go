package filter

import (
	"github.com/pantianying/dubbo-go-proxy/common/constant"
	ct "github.com/pantianying/dubbo-go-proxy/context"
)

func init() {
	SetFilter(constant.CommonFilterName, NewCommonFilter)
}

type CommonFilter struct{}

func NewCommonFilter() Filter {
	return &CommonFilter{}
}
func (f *CommonFilter) OnRequest(ctx ct.ProxyContext) (ret int) {
	return
}
