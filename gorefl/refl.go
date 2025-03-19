// Package gorefl is a sample project for usage of reflection in Go
package gorefl

import (
	"fmt"
	"reflect"
)

// Copy struct data from src to dst
func Copy(dst, src interface{}) error {
	typeDst := reflect.TypeOf(dst)
	if typeDst.Kind() != reflect.Ptr {
		return fmt.Errorf("dst is not a pointer")
	}
	valDst := reflect.ValueOf(dst).Elem()
	valSrc := reflect.ValueOf(src)
	typeSrc := reflect.TypeOf(src)
	for i := 0; i < valSrc.NumField(); i++ {
		srcField := typeSrc.Field(i)
		dstField := valDst.FieldByName(srcField.Name)
		if !dstField.IsValid() {
			continue
		}
		srcTag := typeSrc.Field(i).Tag
		if srcTag.Get("copy") == "nocopy" {
			continue
		}
		dstField.Set(valSrc.Field(i))
	}
	return nil
}
