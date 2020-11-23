package gateway

import "strings"

type ValueConfig string

func (this ValueConfig) GetValue() []string {
	return strings.Split(string(this), ",")
}
