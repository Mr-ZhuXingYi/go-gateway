package gateway

import (
	"fmt"
	"github.com/micro/go-micro/v2/config"
	//"github.com/micro/go-micro/v2/util/log"
)

func InitConfig() Routes {
	configFile := "gateway.yaml"
	err := config.LoadFile(configFile)
	if err != nil {
		fmt.Println(err)
	}
	r := make(Routes, 0)
	err = config.Get("routes").Scan(&r)
	if err != nil {
		fmt.Println(err)
	}

	return r
}
