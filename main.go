package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

var log = logrus.New()

type Request struct {
	SomeString string
	URL        *string
	Index      int
}

type Person struct {
	Id int
}

type OrderedFormatter struct {
	FieldOrder []string
}

func (f *OrderedFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	orderedData := make(map[string]interface{}, len(entry.Data)+1)
	for _, key := range f.FieldOrder {
		if value, ok := entry.Data[key]; ok {
			orderedData[key] = value
		}
	}

	for key, value := range entry.Data {
		if _, ok := orderedData[key]; !ok {
			orderedData[key] = value
		}
	}

	orderedData["msg"] = entry.Message

	serializedData, err := json.Marshal(orderedData)
	if err != nil {
		return nil, err
	}
	return append(serializedData, '\n'), nil
}

func init() {
	log.SetFormatter(&logrus.JSONFormatter{

		DisableTimestamp: true,
	})
	log.SetFormatter(&OrderedFormatter{
		FieldOrder: []string{"Method", "payload", "err", "result"},
	})
	log.SetOutput(os.Stdout)

}
func main() {
	url := "url"
	logPrintJson(
		"Method", "ajax",
		"payload1:request", &Request{SomeString: "someData", URL: &url, Index: 1},
		"payload2:dataset", &[]string{"1", "2", "3,a"},
		"result1:result", &Person{Id: 3},
		"result2:data", &[]int{1, 2, 3, 4, 5},
		"result3:err", errors.New("bad"),
	)

}

func logPrintJson(keyval ...any) {
	logMessages := map[string]interface{}{}
	var keys = make([]string, 0, len(keyval)/2+2)
	var key string
	var value interface{}

	for i, msg := range keyval {
		if i%2 == 0 {
			key = fmt.Sprintf("%v", msg)
			if isError(key) {
				key = "err"
			} else {
				keys = append(keys, key)
			}
		} else {
			value = msg
			logMessages[key] = value
		}
	}

	fields := logrus.Fields{}
	fields["err"] = ErrorToStr(logMessages["err"])

	payload := logrus.Fields{}
	result := logrus.Fields{}

	for _, k := range keys {
		formatedKey := formatedName(k)
		if json, ok := iterableNameAndLengthJson(logMessages[k]); ok && isResult(k) {
			result[formatedKey] = json
		} else if json, ok := structNameAndIDJson(logMessages[k]); ok && isResult(k) {
			result[formatedKey] = json
		} else if isPayload(k) {
			payload[formatedKey] = logMessages[k]
		} else {
			fields[formatedKey] = logMessages[k]
		}
	}
	fields["payload"] = payload
	fields["result"] = result
	log.WithFields(fields).Info("Logged data")
}
