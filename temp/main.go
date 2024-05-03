package main

import (
	"fmt"
	"reflect"
)

var filterList map[reflect.Type][]string

type Example struct {
	Apple string
	Pear  []int
}

func Filter(Struct interface{}) {
	val := reflect.ValueOf(Struct).Elem()
	typ := val.Type()

	if fields, ok := filterList[typ]; ok {
		for i := 0; i < val.NumField(); i++ {
			field := val.Field(i)
			fieldName := typ.Field(i).Name
			if contains(fields, fieldName) && field.CanSet() {
				switch field.Kind() {
				case reflect.String:
					if field.Len() > 20 {
						changedStr := fmt.Sprintf("%s...len(%d)", field.String()[:20], field.Len())
						field.SetString(changedStr)
					}
				case reflect.Slice, reflect.Array:
					originalLen := field.Len()
					if originalLen > 3 {
						newSlice := field.Slice(0, 3)
						field.Set(newSlice)
					}
				}
			}
		}
	}
}



func contains(arr []string, item string) bool {
	for _, s := range arr {
		if s == item {
			return true
		}
	}
	return false
}

func AddFilter(Struct interface{}, StructField ...interface{}) (fields map[reflect.Type][]string) {
	fields = make(map[reflect.Type][]string)

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
				fields[nameOfStruct] = append(fields[nameOfStruct], s.Type().Field(i).Name)
			}
		}
	}
	return fields
}

func main() {
	e := Example{
		Apple: "12345678901234567890123456789012345678901234567890",
		Pear:  []int{1, 2, 3, 4, 5, 6, 7},
	}

	filterList = AddFilter(&e, &e.Apple, &e.Pear)
	fmt.Println(filterList)
	Filter(&e)
	fmt.Println(e)
}
