package gateway

import (
	"net/http"
	"reflect"
	"strings"
)

type Predicates struct {
	Header HeaderPredicate
	Method MethodPredicate
	Host   string
	Path   PathPredicate
}
type Filter struct {
}
type Route struct {
	Id         string
	Url        string
	Predicates Predicates
	Filters    []interface{}
}

type Routes []*Route

func (this Routes) Match(request *http.Request) *Route {
	for _, r := range this {
		if this.isMatch(r, request) {
			return r
		}
	}
	return nil
}
func (this Routes) isMatch(r *Route, request *http.Request) bool {
	v := reflect.ValueOf(r.Predicates)
	for i := 0; i < v.NumField(); i++ {
		if matcher, ok := v.Field(i).Interface().(PredicateMatcher); ok &&
			strings.Trim(v.Field(i).String(), " ") != "" {
			if !matcher.Match(request) {
				return false
			}
		}
	}

	return true

}
