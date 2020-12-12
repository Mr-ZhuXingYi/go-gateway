package gateway

import (
	"net/http"
	"reflect"
	"sort"
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
	GetOrder() int
}

type ServerWebExchange struct { //exchange 用来保存请求和响应等上下文
	Request *http.Request
}

func BuildServerWebExchange(req *http.Request) *ServerWebExchange {
	return &ServerWebExchange{Request: req}
}

type ResponseFilter func(*http.Response)
type ResponseFilters []ResponseFilter

func (this ResponseFilters) Filter(r *http.Response) {
	for _, f := range this {
		f(r)
	}
}

type GatewayFilter func(exchange *ServerWebExchange) ResponseFilter

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

func (this SimpleFilter) getClass() FilterFactory {
	vlist := strings.Split(string(this), "=")
	if f, ok := FilterMap.Load(vlist[0]); ok {
		return f.(FilterFactory)
	}
	return nil
}

// 正式反代之前做过滤（请求处理)
func (this *Route) FilterRequest(exchange *ServerWebExchange) ResponseFilters {
	//for _, f := range this.Filters {
	//	v := reflect.ValueOf(f)
	//	if v.Kind() == reflect.String {
	//		if gate_filter := SimpleFilter(v.String()).filter(); gate_filter != nil {
	//			gate_filter(exchange)
	//		}
	//	}
	//}
	ret := make([]ResponseFilter, 0)

	this.sortFilter()
	//fmt.Println(this.orderFilters,this.Filters) // 测试用
	for _, f := range this.orderFilters {
		//目前这部分代码 只支持string类型的simplefilter
		rspFilter := f.(SimpleFilter).filter()(exchange)
		if rspFilter != nil {
			ret = append(ret, rspFilter)
		}
	}
	return ret
}

func (this *Route) sortFilter() {
	this.orderFilters = make([]interface{}, 0)
	for _, f := range this.Filters { // 字符串
		v := reflect.ValueOf(f)
		if v.Kind() == reflect.String {
			if obj := SimpleFilter(v.String()).getClass(); obj != nil {
				this.orderFilters = append(this.orderFilters, SimpleFilter(v.String()))
			}
		}
	}
	//排序 (稳定排序)
	sort.SliceStable(this.orderFilters, func(i, j int) bool {
		return this.orderFilters[i].(SimpleFilter).getClass().GetOrder() <
			this.orderFilters[j].(SimpleFilter).getClass().GetOrder()
	})
}
