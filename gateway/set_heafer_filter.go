package gateway

import "net/http"

func init() {
	RegisterFilter("SetHeader", NewSetHeaderFilter())
}

type SetHeaderFilter struct{}

func (this *SetHeaderFilter) Apply(config interface{}) GatewayFilter {
	return func(exchange *ServerWebExchange) ResponseFilter {
		p := NameConfig(config.(string))
		if headers := p.GetValue(); headers != nil {
			for _, header := range headers {
				exchange.Request.Header.Set(header.Namme, header.Value)
			}
		}
		return func(response *http.Response) {
			response.Header.Add("age", "19")
		}
	}
}

func (this *SetHeaderFilter) GetOrder() int {
	return 2
}

func NewSetHeaderFilter() *SetHeaderFilter {
	return &SetHeaderFilter{}
}
