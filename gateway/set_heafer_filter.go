package gateway

func init() {
	RegisterFilter("SetHeader", NewSetHeaderFilter())
}

type SetHeaderFilter struct{}

func (this *SetHeaderFilter) Apply(config interface{}) GatewayFilter {
	return func(exchange *ServerWebExchange) {
		p := NameConfig(config.(string))
		if headers := p.GetValue(); headers != nil {
			for _, header := range headers {
				exchange.Request.Header.Set(header.Namme, header.Value)
			}
		}
	}
}

func NewSetHeaderFilter() *SetHeaderFilter {
	return &SetHeaderFilter{}
}
