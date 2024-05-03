package main

import (
	"reflect"
	"regexp"
	"strings"
)

func isSlice(v interface{}) bool {

	val := reflect.ValueOf(v)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() == reflect.Slice || val.Kind() == reflect.Array {
		return true
	}
	return false
}
func isStruct(v interface{}) bool {
	val := reflect.ValueOf(v)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	return val.Kind() == reflect.Struct
}

func getErrorValue(val interface{}) error {
	if err, ok := val.(error); ok {
		return err
	}
	return nil
}

func isResult(key string) bool {
	re := regexp.MustCompile(`result`)
	return re.MatchString(key)
}

func isError(key string) bool {
	re := regexp.MustCompile(`err$`)
	return re.MatchString(key)
}

func isPayload(key string) bool {
	re := regexp.MustCompile(`payload\d+`)
	return re.MatchString(key)
}

func isMethod(key string) bool {
	re := regexp.MustCompile(`method`)
	return re.MatchString(strings.ToLower(key))
}

func formatedName(key string) string {
	index := getIndex(key, ':')
	if index == -1 {
		return key
	}
	return key[index+1:]

}

func getIndex(text string, char rune) int {
	for i, ch := range text {
		if char == ch {
			return i
		}
	}
	return -1
}

// =================
// for json
type TypeAndID struct {
	Type string
	ID   int
}

func structNameAndIDJson(v interface{}) (TypeAndID, bool) {
	val := reflect.ValueOf(v)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() == reflect.Struct {
		idField := val.FieldByName("ID")
		if !idField.IsValid() {
			idField = val.FieldByName("Id")
		}
		if idField.IsValid() {
			return TypeAndID{Type: val.Type().Name(), ID: idField.Interface().(int)}, true
		}
	}
	return TypeAndID{}, false
}

type TypeAndLength struct {
	Type   string
	Length int
}

func iterableNameAndLengthJson(v any) (TypeAndLength, bool) {
	val := reflect.ValueOf(v)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() == reflect.Slice || val.Kind() == reflect.Array {
		return TypeAndLength{Type: val.Type().String(), Length: val.Len()}, true
	}
	return TypeAndLength{}, false
}
