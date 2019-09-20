package server

import (
	"fmt"
	"strings"
)

var UrlMissingParam = &customErr{code: "MISSING_PARAM", txt: ""}

type customErr struct {
	code string
	txt  string
	val  []interface{}
}

func (e *customErr) Error() string {
	txt := e.txt
	val := e.val
	if len(txt) > 0 {
		txt = fmt.Sprintf(" %v", txt)
	}
	if val != nil && len(val) > 0 {
		values := []string{}
		for _, v := range val {
			values = append(values, fmt.Sprintf("%v", v))
		}
		txt = fmt.Sprintf("%v %v", txt, strings.Join(values, ","))
	}
	return fmt.Sprintf("[%v]%v", e.code, txt)
}

func (e *customErr) Val(i ...interface{}) *customErr {
	e.val = append(e.val, i...)
	return e
}
