package main

import (
	"fmt"
	"reflect"
	"regexp"
)

func isResult(key string) bool {
	re := regexp.MustCompile(`result\d+`)
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

// =====================
// for text

func iterableNameAndLength(v any) (string, int, bool) {
	val := reflect.ValueOf(v)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() == reflect.Slice || val.Kind() == reflect.Array {
		return val.Type().String(), val.Len(), true
	}
	return "", 0, false
}

func structNameAndID(v interface{}) (string, interface{}, bool) {
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
			return val.Type().Name(), idField.Interface(), true
		}
	}
	return "", nil, false
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

func ErrorToStr(value interface{}) string {
	val := reflect.ValueOf(value)

	if !val.IsValid() {
		return "nil"
	}

	if val.Kind() == reflect.Ptr && !val.IsNil() {
		val = val.Elem()
	}

	if err, ok := val.Interface().(error); ok {
		return err.Error()
	}
	str := fmt.Sprintf("%v", val)[1:]
	return str[:len(str)-1]
}
