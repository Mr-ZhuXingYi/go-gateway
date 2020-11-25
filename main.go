package main

import (
	"fmt"
	"go-gateway/gateway"

	"net/http"
	"net/http/httputil"
	"net/url"
)

func main() {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Println(request.URL.Path)
		r := gateway.InitConfig()
		if router := r.Match(request); router != nil {
			exchange := gateway.BuildServerWebExchange(request)
			router.FilterBefore(exchange)
			fmt.Println(request.Header)
			remote, err := url.Parse(router.Url)
			if err != nil {
				fmt.Println(err)
			}
			p := httputil.NewSingleHostReverseProxy(remote)
			p.ServeHTTP(writer, request)
		} else {
			writer.WriteHeader(http.StatusBadRequest)
		}
	})
	if err := http.ListenAndServe(":8081", nil); err != nil {
		fmt.Println(err)
	}
}
