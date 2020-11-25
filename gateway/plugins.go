package gateway

import (
	"github.com/robertkrimen/otto"
	"io/ioutil"
	"log"
	"path"
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
	Name string
	Main otto.Value
}

func (self *FilterPlugin) Filter(exchange *ServerWebExchange) {
	this, _ := otto.ToValue(nil)
	_, err := self.Main.Call(this, exchange) //调用 闭包 ,传递上下文参数
	if err != nil {
		log.Println(err)
	}
}
func (self *FilterPlugin) Apply(config interface{}) GatewayFilter {
	return func(exchange *ServerWebExchange) {
		self.Filter(exchange)
	}
}
func loadPlugin(js *otto.Otto) *FilterPlugin {
	filter_name, err := js.Call("name", nil)
	if err != nil {
		log.Println(err)
		return nil
	}
	filter_main, err := js.Call("main", nil)
	if err != nil || !filter_main.IsFunction() {
		log.Println(err)
		return nil
	}
	log.Println(filter_name.ToString())
	return &FilterPlugin{Name: filter_name.String(), Main: filter_main}
}
func loadPlugins(dirname string) []*FilterPlugin {
	ret := make([]*FilterPlugin, 0)
	fileInfos, _ := ioutil.ReadDir(dirname)
	for _, file := range fileInfos {
		if !file.IsDir() && path.Ext(file.Name()) == ".js" {
			js := otto.New()
			_, err := js.Run(readFile(dirname + "/" + file.Name()))
			if err != nil {
				log.Println(err)
				continue
			}
			if p := loadPlugin(js); p != nil {
				ret = append(ret, p)
			}
		}
	}
	return ret

}
