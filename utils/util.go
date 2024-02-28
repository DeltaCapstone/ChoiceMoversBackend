package utils

import (
	"fmt"
	"reflect"
)

func StructToMap(data interface{}, tag string) map[string]interface{} {
	result := make(map[string]interface{})
	value := reflect.ValueOf(data)
	typ := reflect.TypeOf(data)

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		tag := field.Tag.Get(tag) // You can customize this tag based on your needs

		// If the tag is not empty, use it as the key in the map
		if tag != "" {
			result[tag] = value.Field(i).Interface()
			fmt.Println(tag, result)
		}
	}
	return result
}
