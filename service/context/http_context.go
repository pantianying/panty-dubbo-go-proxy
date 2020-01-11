package context

import (
	"github.com/pantianying/dubbo-go-proxy/common/constant"
	"github.com/pantianying/dubbo-go-proxy/common/util"
	"github.com/pantianying/dubbo-go-proxy/service"
	"io/ioutil"
	"net/http"
)

type httpContext struct {
	*service.BaseContext
	r       *http.Request
	w       http.ResponseWriter
	bodyMap map[string]interface{}
}

func NewHttpContext(w http.ResponseWriter, r *http.Request) service.ProxyContext {
	ctx := &httpContext{
		BaseContext: service.NewBaseContext([]service.Filter{service.GetFilter(constant.CommonFilterName)}),
		r:           r,
		w:           w,
		bodyMap:     make(map[string]interface{}),
	}
	if body, err := ioutil.ReadAll(r.Body); err == nil && len(body) > 0 {
		err = util.ParseJsonByStruct(body, &ctx.bodyMap)
	}
	return ctx
}
