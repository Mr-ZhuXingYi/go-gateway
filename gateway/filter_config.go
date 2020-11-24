package gateway

import "strings"

type ValueConfig string

func (this ValueConfig) GetValue() []string {
	return strings.Split(string(this), ",")
}

type NameConfigObject struct {
	Namme string
	Value string
}
type NameConfig string

func (this NameConfig) GetValue() []*NameConfigObject {
	slist := strings.Split(string(this), ",")
	if len(slist) < 2 || len(slist)%2 != 0 {
		return nil
	}

	ret := make([]*NameConfigObject, 0)
	for i := 0; i < len(slist); i += 2 {
		ret = append(ret, &NameConfigObject{Namme: slist[i], Value: slist[i+1]})
	}
	return ret
}
