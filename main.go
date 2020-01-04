package main

import (
	_ "github.com/pantianying/dubbo-go-proxy/filter"
	"github.com/pantianying/dubbo-go-proxy/proxy/http"
)

// for example
func main() {
	http.Run()
}
