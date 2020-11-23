package gateway

import (
	"net/http"
	"reflect"
	"strings"
	"sync"
)

var FilterMap sync.Map // 存储 所有filter
//注册过滤器
func RegisterFilter(name string, value FilterFactory) {
	FilterMap.Store(name, value)
}

type FilterFactory interface { //所有过滤器的接口
	Apply(config interface{}) GatewayFilter
}
type ServerWebExchange struct { //exchange 用来保存请求和响应等上下文
	Request *http.Request
}

func BuildServerWebExchange(req *http.Request) *ServerWebExchange {
	return &ServerWebExchange{Request: req}
}

type GatewayFilter func(exchange *ServerWebExchange)

// 字符串类型的Filter配置
type SimpleFilter string

func (this SimpleFilter) filter() GatewayFilter {
	vlist := strings.Split(string(this), "=")
	if len(vlist) != 2 {
		return nil
	}
	if f, ok := FilterMap.Load(vlist[0]); ok {
		return f.(FilterFactory).Apply(vlist[1])
	}
	return nil
}

// 正式反代之前做过滤（处理)
func (this *Route) FilterBefore(exchange *ServerWebExchange) {
	for _, f := range this.Filters {
		v := reflect.ValueOf(f)
		if v.Kind() == reflect.String {
			if gate_filter := SimpleFilter(v.String()).filter(); gate_filter != nil {
				gate_filter(exchange)
			}
		}
	}
}
