package main

import (
	"fmt"
	"reflect"
)

var filterList map[reflect.Type][]string

func init() {
	filterList = make(map[reflect.Type][]string)
}

func Filter(Struct interface{}) interface{} {
	val := reflect.ValueOf(Struct)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	typ := val.Type()
	newInstance := reflect.New(val.Type()).Elem()
	if fields, ok := filterList[typ]; ok {
		for i := 0; i < val.NumField(); i++ {
			field := val.Field(i)
			fieldName := typ.Field(i).Name
			newField := newInstance.Field(i)
			if contains(fields, fieldName) {
				switch field.Kind() {
				case reflect.String:
					if field.Len() > 20 {
						changedStr := fmt.Sprintf("%s...len(%d)", field.String()[:20], field.Len())
						newField.SetString(changedStr)
					} else {
						newField.Set(field)
					}
				case reflect.Slice, reflect.Array:
					FilterArray(field)
				default:
					newField.Set(field)
				}
			} else {
				newField.Set(field)
			}
		}
	} else {
		newInstance.Set(val)
	}

	return newInstance.Interface()
}

func FilterArray(arr interface{}) interface{} {
	sliceVal := reflect.ValueOf(arr)
	if sliceVal.Kind() != reflect.Slice {
		fmt.Println("FilterArray: provided data is not a slice")
		return nil
	}

	resultSlice := reflect.MakeSlice(sliceVal.Type(), 0, sliceVal.Len())

	for i := 0; i < sliceVal.Len(); i++ {
		element := sliceVal.Index(i)
		if !element.CanAddr() {
			newElem := reflect.New(element.Type()).Elem()
			newElem.Set(element)
			element = newElem.Addr()
		}
		filtered := Filter(element.Interface())
		resultSlice = reflect.Append(resultSlice, reflect.ValueOf(filtered))
	}

	return resultSlice.Interface()
}

func contains(arr []string, item string) bool {
	for _, s := range arr {
		if s == item {
			return true
		}
	}
	return false
}

func AddFilter(Struct interface{}, StructField ...interface{}) {

	nameOfStruct := reflect.ValueOf(Struct).Elem().Type()

	for r := range StructField {
		fieldVal := reflect.ValueOf(StructField[r])
		if fieldVal.Kind() != reflect.Ptr || fieldVal.IsNil() {
			fmt.Println("Error: StructField is not a pointer or is nil.")
			continue
		}
		s := reflect.ValueOf(Struct).Elem()
		f := reflect.ValueOf(StructField[r]).Elem()

		for i := 0; i < s.NumField(); i++ {
			valueField := s.Field(i)
			if valueField.Addr().Interface() == f.Addr().Interface() {
				filterList[nameOfStruct] = append(filterList[nameOfStruct], s.Type().Field(i).Name)
			}
		}
	}
}
