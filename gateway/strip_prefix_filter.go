package gateway

import (
	"log"
	"strconv"
	"strings"
)

type StripPrefixFilter struct{}

func init() {
	RegisterFilter("StripPrefix", NewStripPrefixFilter())
}
func (this *StripPrefixFilter) Apply(config interface{}) GatewayFilter {
	return func(exchange *ServerWebExchange) {
		path := exchange.Request.URL.Path
		//   /v1/course  ==>  [ v1 course]

		defIndex := 1

		config := ValueConfig(config.(string)).GetValue()
		i, err := strconv.Atoi(config[0])
		if err != nil {
			log.Fatal(err)
		} else {
			defIndex = i
		}
		path_list := strings.Split(path, "/")
		exchange.Request.URL.Path = strings.Join(path_list[defIndex+1:], "/")
	}
}
func NewStripPrefixFilter() *StripPrefixFilter {
	return &StripPrefixFilter{}
}
