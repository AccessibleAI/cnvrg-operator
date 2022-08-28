package controllers

import "reflect"

type CnvrgSpecBoolTransformer struct{}

func ContainsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

func RemoveString(slice []string, s string) (result []string) {
	for _, item := range slice {
		if item == s {
			continue
		}
		result = append(result, item)
	}
	return
}

func (t CnvrgSpecBoolTransformer) Transformer(typ reflect.Type) func(dst, src reflect.Value) error {
	if typ == reflect.TypeOf(true) {
		return func(dst, src reflect.Value) error {
			if dst.CanSet() {
				// always set boolean value
				// e.g always do the WithOverwriteWithEmptyValue
				// but only for booleans
				dst.Set(src)
			}
			return nil
		}
	}
	return nil
}
