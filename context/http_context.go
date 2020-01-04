package context

import (
	"github.com/pantianying/dubbo-go-proxy/common/constant"
	"github.com/pantianying/dubbo-go-proxy/common/util"
	"io/ioutil"
	"net/http"
)

type httpContext struct {
	*baseContext
	r       *http.Request
	w       http.ResponseWriter
	bodyMap map[string]interface{}
}

func NewhttpContext(w http.ResponseWriter, r *http.Request) ProxyContext {
	ctx := &httpContext{
		baseContext: NewBaseContext([]string{constant.CommonFilterName}),
		r:           r,
		w:           w,
		bodyMap:     make(map[string]interface{}),
	}
	if body, err := ioutil.ReadAll(r.Body); err == nil && len(body) > 0 {
		err = util.ParseJsonByStruct(body, &ctx.bodyMap)
	}
	return ctx
}
