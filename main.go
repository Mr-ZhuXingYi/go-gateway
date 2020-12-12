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
		r := gateway.InitConfig()
		if router := r.Match(request); router != nil {

			exchange := gateway.BuildServerWebExchange(request)

			RspFilters := router.FilterRequest(exchange)

			fmt.Println(request.URL.Query()) //测试
			remote, err := url.Parse(router.Url)
			if err != nil {
				fmt.Println(err)
			}
			p := httputil.NewSingleHostReverseProxy(remote)
			p.ModifyResponse = func(response *http.Response) error {
				RspFilters.Filter(response)
				return nil
			}
			p.ServeHTTP(writer, request)
		} else {
			writer.WriteHeader(http.StatusBadRequest)
		}
	})
	if err := http.ListenAndServe(":8081", nil); err != nil {
		fmt.Println(err)
	}
}
