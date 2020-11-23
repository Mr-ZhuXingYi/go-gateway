package gateway

import (
	"log"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"
)

type PredicateMatcher interface {
	Match(request *http.Request) bool
}

type HeaderPredicate string

func (this HeaderPredicate) Match(request *http.Request) bool {
	param := request.Header
	s := string(this)
	if strings.Trim(s, " ") == "" {
		return true
	}
	predicates := strings.Split(s, ",")
	if len(predicates) < 2 || len(predicates)%2 != 0 {
		return true
	}

	for i := 0; i < len(predicates); i += 2 {
		key, pattern := predicates[i], predicates[i+1]
		if v, ok := param[key]; ok {
			reg, err := regexp.Compile(pattern)
			if err != nil {
				log.Fatal(err)
				return true
			}
			if !reg.MatchString(v[0]) {
				return false
			}

		} else {
			return false
		}

	}

	return true
}

type PathPredicate string

func (this PathPredicate) Match(request *http.Request) bool {
	param := request.URL.Path
	if strings.Trim(string(this), " ") == "" {
		return true
	}
	m, err := filepath.Match(string(this), param)
	if err != nil || !m {
		return false
	}
	return true
}

type MethodPredicate string

func (this MethodPredicate) Match(req *http.Request) bool {
	param := req.Method
	s := string(this)
	slist := strings.Split(s, ",")
	if len(slist) == 0 {
		return true
	}
	for _, item := range slist {
		if item == param {
			return true
		}
	}
	return false
}
