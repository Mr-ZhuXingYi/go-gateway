package gateway

import (
	"github.com/robertkrimen/otto"
	"io/ioutil"
	"log"
	"path"
	"strings"
)

func init() {
	plugins := loadPlugins("plugins")
	for _, plugin := range plugins {
		RegisterFilter(plugin.Name, plugin)
	}
}
func readFile(file string) string {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return ""
	}
	return string(b)
}

type FilterPlugin struct {
	Name  string
	Main  otto.Value
	Order int
}

func (self *FilterPlugin) Filter(params ...interface{}) {
	this, _ := otto.ToValue(nil)
	_, err := self.Main.Call(this, params...) //调用 闭包 ,传递上下文参数
	if err != nil {
		log.Println("self.Main.Call(this, params):", err)
	}
}
func (self *FilterPlugin) Apply(config interface{}) GatewayFilter {
	return func(exchange *ServerWebExchange) ResponseFilter {
		slis := strings.Split(config.(string), ",")
		params := make([]interface{}, 0)
		params = append(params, exchange)
		for _, s := range slis {
			params = append(params, s)
		}
		self.Filter(params...)
		return nil
	}
}

func (self *FilterPlugin) GetOrder() int {
	return self.Order
}

func loadPlugin(js *otto.Otto) *FilterPlugin {
	filter_name, err := js.Call("name", nil)
	if err != nil {
		log.Println("err := js.Call(\"name\", nil):", err)
		return nil
	}
	filter_main, err := js.Call("main", nil)
	if err != nil || !filter_main.IsFunction() {
		log.Println("js.Call(\"main\", nil):", err)
		return nil
	}

	filter_order, err := js.Call("order", nil)
	if !filter_order.IsNumber() || err != nil {
		return nil
	}

	order, _ := filter_order.ToInteger()

	//log.Println(filter_name.ToString())
	return &FilterPlugin{Name: filter_name.String(), Main: filter_main, Order: int(order)}
}
func loadPlugins(dirname string) []*FilterPlugin {
	ret := make([]*FilterPlugin, 0)
	fileInfos, _ := ioutil.ReadDir(dirname)
	for _, file := range fileInfos {
		if !file.IsDir() && path.Ext(file.Name()) == ".js" {
			js := otto.New()
			_, err := js.Run(readFile(dirname + "/" + file.Name()))
			if err != nil {
				log.Println("js.Run:", err)
				continue
			}
			if p := loadPlugin(js); p != nil {
				ret = append(ret, p)
			}
		}
	}
	return ret

}
