package filter

import (
	"github.com/pantianying/dubbo-go-proxy/common/constant"
	"github.com/pantianying/dubbo-go-proxy/service"
)

func init() {
	service.SetFilter(constant.CommonFilterName, NewCommonFilter)
}

type CommonFilter struct{}

func NewCommonFilter() service.Filter {
	return &CommonFilter{}
}
func (f *CommonFilter) OnRequest(ctx service.ProxyContext) (ret int) {
	return
}
