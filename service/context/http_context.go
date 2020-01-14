package context

import (
	"github.com/pantianying/dubbo-go-proxy/common/constant"
	"github.com/pantianying/dubbo-go-proxy/common/errcode"
	"github.com/pantianying/dubbo-go-proxy/common/logger"
	"github.com/pantianying/dubbo-go-proxy/common/util"
	"github.com/pantianying/dubbo-go-proxy/dubbo"
	"github.com/pantianying/dubbo-go-proxy/service"
	"io/ioutil"
	"net/http"
	"strings"
)

type httpContext struct {
	*service.BaseContext
	r       *http.Request
	w       http.ResponseWriter
	bodyMap map[string]interface{}

	mdKey  *service.MetadataIdentifier
	mdInfo *service.MetadataInfo

	metadataCenter service.MetadataCenter
}

func NewHttpContext(w http.ResponseWriter, r *http.Request) service.ProxyContext {
	ctx := &httpContext{
		BaseContext: service.NewBaseContext([]service.Filter{
			service.GetFilter(constant.MatchFilterName)}),
		r:       r,
		w:       w,
		bodyMap: make(map[string]interface{}),
	}
	if body, err := ioutil.ReadAll(r.Body); err == nil && len(body) > 0 {
		err = util.ParseJsonByStruct(body, &ctx.bodyMap)
	}
	return ctx
}
func (hc *httpContext) Match() (*service.MetadataIdentifier, *service.MetadataInfo, int) {
	//todo 需要支持更复杂的match，包括restful，等等
	if hc.mdKey != nil {
		return hc.mdKey, hc.mdInfo, errcode.Success
	}
	ss := strings.Split(hc.r.URL.Path, "/")
	if len(ss) < 2 {
		return nil, nil, errcode.NotFind
	}
	application := ss[0]
	serviceInterfaceName := ss[1]
	group := hc.r.FormValue("group")
	version := hc.r.FormValue("version")
	hc.mdKey = &service.MetadataIdentifier{
		ServiceInterface: serviceInterfaceName,
		Version:          version,
		Group:            group,
		Application:      application,
	}
	hc.mdInfo = hc.metadataCenter.GetProviderMetaData(hc.mdKey)
	return hc.mdKey, hc.mdInfo, errcode.Success
}
func (hc *httpContext) InvokeData() *dubbo.InvokeData {
	if hc.mdKey == nil {
		logger.Error("httpContext.MetadataIdentifier is nil, check if had do match()")
		return nil
	}
	method := hc.getMethod()
	parameterTypes, ok := hc.getParameterTypes(method)
	if !ok {
		return nil
	}

	invokeData := &dubbo.InvokeData{
		InterfaceName:  hc.mdKey.ServiceInterface,
		Group:          hc.mdKey.Group,
		Version:        hc.mdKey.Version,
		Method:         method,
		ParameterTypes: parameterTypes,
		//ReqData:
	}
	return invokeData
}
func (hc *httpContext) getMethod() string {
	//todo
	return hc.bodyMap["method"].(string)
}
func (hc *httpContext) getParameterTypes(method string) ([]string, bool) {
	//todo 优化
	if hc.mdInfo != nil && len(hc.mdInfo.Methods) > 0 {
		for _, v := range hc.mdInfo.Methods {
			if v.Name == method {
				return v.ParameterTypes, true
			}
		}
	}
	// 元数据中获取不到从body中获取
	if p, ok := hc.bodyMap["parameterTypes"]; ok {
		parameterTypes, ok := p.([]string)
		return parameterTypes, ok
	}
	return nil, false
}
func (hc *httpContext) getReqData() ([]interface{}, bool) {
	//todo
	return nil, false
}
