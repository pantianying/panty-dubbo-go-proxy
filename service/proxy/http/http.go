package http

import (
	"github.com/pantianying/dubbo-go-proxy/common/config"
	"github.com/pantianying/dubbo-go-proxy/common/errcode"
	"github.com/pantianying/dubbo-go-proxy/common/logger"
	"github.com/pantianying/dubbo-go-proxy/dubbo"
	"github.com/pantianying/dubbo-go-proxy/service"
	ct "github.com/pantianying/dubbo-go-proxy/service/context"
	"github.com/pantianying/dubbo-go-proxy/service/metadata/redis"
	"io"
	"net/http"
	"time"
)

var srv http.Server
var mc = redis.NewRedisMetaDataCenter()

func Run() {
	startHttpServer()
}
func startHttpServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", commonHandle)
	srv = http.Server{
		Addr:           config.Config.HttpListenAddr,
		Handler:        mux,
		ReadTimeout:    20 * time.Second,
		WriteTimeout:   20 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	err := srv.ListenAndServe()
	if err != nil {
		if err != http.ErrServerClosed {
			logger.Errorf("http.ListenAndServe err:%v", err)
		}
	}
}
func commonHandle(w http.ResponseWriter, r *http.Request) {
	setJsHeader(w, r)
	if r.Method == http.MethodOptions {
		return
	}
	var (
		ret         int
		responseStr string
		ctx         = ct.NewHttpContext(w, r)
	)
	defer func() {
		responseStr = getRsp(ctx, ret)
		//返回结果
		io.WriteString(w, responseStr)

	}()
	filter := ctx.NextFilter()
	for filter != nil {
		ret = filter.OnRequest(ctx)
		if ret != errcode.Success {
			return
		}
		filter = ctx.NextFilter()
	}
	dubbo.Client.Call(*ctx.InvokeData())
	return

}

func getRsp(ctx service.ProxyContext, ret int) string {
	if ret == errcode.Success {
		//todo
		return ""
	}
	return errcode.GetMsg(ret)
}

func setJsHeader(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Set("Access-Control-Allow-Headers", r.Header.Get("Access-Control-Request-Headers"))
	w.Header().Add("Access-Control-Allow-Headers", "access_token,Content-Type")

	w.Header().Add("Access-Control-Allow-Methods", "POST")
	w.Header().Add("Access-Control-Allow-Methods", "OPTIONS")
	w.Header().Add("Access-Control-Allow-Methods", "GET")
	w.Header().Add("Access-Control-Allow-Methods", "DELETE")
	w.Header().Add("Access-Control-Allow-Methods", "PUT")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
}
