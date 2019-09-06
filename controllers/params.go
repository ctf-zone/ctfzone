package controllers

import (
	"reflect"
	"strings"

	"github.com/gorilla/schema"
)

var decoder = schema.NewDecoder()

func init() {
	decoder.RegisterConverter(map[string]interface{}{}, func(input string) reflect.Value {
		m := make(map[string]interface{})
		for _, pair := range strings.Split(input, ";") {
			parts := strings.Split(pair, ":")
			if len(parts) == 2 {
				m[parts[0]] = parts[1]
			}
		}
		return reflect.ValueOf(m)
	})
}
