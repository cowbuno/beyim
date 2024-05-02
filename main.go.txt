package main

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
)

var log *slog.Logger

type Request struct {
	SomeString string
	URL        string
	Index      int
}

type Person struct {
	Id int
}

func main() {

	log = slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
	)
	logPrintJson(
		"Method", "ajax",
		"payload1:request", &Request{SomeString: "someData", URL: "url", Index: 1},
		"payload2:dataset", &[]string{"1", "2", "3,a"},
		"result1:result", &Person{Id: 3},
		"result2:data", &[]int{1, 2, 3, 4, 5},
		"result3:err", errors.New("bad"),
	)

}

func logPrint(keyval ...any) {
	logMessages := map[string]interface{}{}
	var keys = make([]string, 0, len(keyval)/2+2)
	var errKey string
	var key string
	var value interface{}
	method := log.Info

	for i, msg := range keyval {
		if i%2 == 0 {
			key = fmt.Sprintf("%v", msg)
			if isError(key) {
				errKey = "err"
				key = "err"
			} else {
				keys = append(keys, key)
			}
		} else {
			value = msg
			logMessages[key] = value
		}
	}

	if logMessages[errKey] != nil {
		method = log.Error
	}

	fields := make([]any, 0, len(keys))

	isPathErr := -1
	for _, k := range keys {
		if isPathErr == -1 && isPayload(k) {
			isPathErr = 0
		} else if isPathErr == 0 && !isPayload(k) {
			isPathErr = 1
		}

		if isPathErr == 1 {
			fields = append(fields, slog.Any(errKey, logMessages[errKey]))
			isPathErr = 2
		}
		formatedKey := formatedName(k)
		if nameOfType, length, ok := iterableNameAndLength(logMessages[k]); ok && isResult(k) {
			fields = append(fields, slog.String(formatedKey, fmt.Sprintf("%s.len(%d)", nameOfType, length)))
		} else if nameOfType, id, ok := structNameAndID(logMessages[k]); ok && isResult(k) {
			fields = append(fields, slog.Any(formatedKey, fmt.Sprintf("%s.id(%v)", nameOfType, id)))
		} else {
			fields = append(fields, formatedKey, logMessages[k])

		}
	}
	method("", fields...)
}

func logPrintJson(keyval ...any) {
	logMessages := map[string]interface{}{}
	var keys = make([]string, 0, len(keyval)/2+2)
	var errKey string
	var key string
	var value interface{}
	method := log.Info

	for i, msg := range keyval {
		if i%2 == 0 {
			key = fmt.Sprintf("%v", msg)
			if isError(key) {
				errKey = "err"
				key = "err"
			} else {
				keys = append(keys, key)
			}
		} else {
			value = msg
			logMessages[key] = value
		}
	}

	if logMessages[errKey] != nil {
		method = log.Error
	}

	fields := make([]any, 0, len(keys))

	isPathErr := -1
	for _, k := range keys {
		if isPathErr == -1 && isPayload(k) {
			isPathErr = 0
		} else if isPathErr == 0 && !isPayload(k) {
			isPathErr = 1
		}

		if isPathErr == 1 {
			fields = append(fields, slog.Any(errKey, logMessages[errKey]))
			isPathErr = 2
		}
		formatedKey := formatedName(k)
		if json, ok := iterableNameAndLengthJson(logMessages[k]); ok && isResult(k) {
			fields = append(fields, formatedKey, json)
		} else if json, ok := structNameAndIDJson(logMessages[k]); ok && isResult(k) {
			fields = append(fields, formatedKey, json)
		} else {
			fields = append(fields, formatedKey, logMessages[k])

		}
	}
	method("", fields...)
}
