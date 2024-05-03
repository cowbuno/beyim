package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"runtime/debug"
	"time"

	"github.com/iancoleman/orderedmap"
	"github.com/sirupsen/logrus"
)

var logger = logrus.New()

func init() {
	logger.Formatter = new(CustomFormatter)
	logger.Out = os.Stdout

	level, err := logrus.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		level = logrus.InfoLevel
	}
	logger.SetLevel(level)

	logBuffer := make(chan *logrus.Entry, 5000)
	go processLogEntries(logBuffer)
	logger.Hooks.Add(&AsyncHook{logBuffer})
}

type CustomFormatter struct{}

func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	oMap := orderedmap.New()
	for _, field := range []string{"timestamp", "level", "method", "payload", "request_at", "duration_time", "result", "error", "message", "description"} {
		if val, ok := entry.Data[field]; ok {
			oMap.Set(field, val)
		}
	}

	serialized, err := json.Marshal(oMap)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal fields to JSON, %v", err)
	}
	return append(serialized, '\n'), nil
}

type AsyncHook struct {
	logBuffer chan<- *logrus.Entry
}

func (hook *AsyncHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (hook *AsyncHook) Fire(entry *logrus.Entry) error {
	hook.logBuffer <- entry
	return nil
}

func processLogEntries(entries chan *logrus.Entry) {
	for entry := range entries {
		logger.Out.Write([]byte(entry.Message))
	}
}

func main() {

	payload := `{
		"success": true,
		"message": "",
		"data": {
			"posts": [{
				"microtopicId": 22,
				"objectiveId": 22,
				"subject": "Физика",
				"microtopic": "Скалярные и векторные величины",
				"objective": "7.1.1.6 – различать скалярные и векторные величины и приводить примеры",
				"iconUrl": "https://cf-beyim-subject.beyim.ai/1--97944bef-fb50-456b-b802-abfddebcd1dc.png",
				"subjectId": 2,
				"id": "656e2ad13e8f14f57aa5d31c",
				"category": "video",
				"resources": ["https://cf-beyim-video.beyim.ai/687274d3-35c6-4e52-a4b5-885dbb3614b6/hls/04a6640e43b34de4978d50694820684c--skalyarnye_i_vektornye_velichiny.m3u8"],
				"contentId": "656e28b23e8f14f57aa5d316",
				"description": "{\"root\":{\"children\":[{\"children\":[{\"detail\":0,\"format\":0,\"mode\":\"normal\",\"style\":\"\",\"text\":\"Скалярные и векторные величины \",\"type\":\"text\",\"version\":1}],\"direction\":\"ltr\",\"format\":\"\",\"indent\":0,\"type\":\"paragraph\",\"version\":1}],\"direction\":\"ltr\",\"format\":\"\",\"indent\":0,\"type\":\"root\",\"version\":1}}",
				"thumbnail": "https://cf-beyim-thumbnail.beyim.ai/04a6640e43b34de4978d50694820684c--skalyarnye_i_vektornye_velichiny.mp4--80a605c345064ddf9bc1d36fa5e8c63e--thumbnail.jpg"
			}]
		}
	}`

	var data map[string]interface{}

	json.Unmarshal([]byte(payload), &data)
	url := "urlPt"
	req := Request{SomeString: "123345678912345678912345678143356754654632141324", URL: &url, Index: 1}
	AddFilter(&req, &req.SomeString)
	logPrint(
		"Method", "get",
		"payload1:request", req,
		// "payload2:dataset", data,
		"result1:result", Person{Id: 1},
		"result2:data", &[]int{1, 2, 3, 4, 5},
		"result3:err", errors.New("bad"),
	)

}

func logPrint(keyval ...any) {
	startTime := time.Now()

	logMessages := map[string]interface{}{}
	var keys = make([]string, 0, len(keyval)/2+2)
	var key string
	var value interface{}

	for i, msg := range keyval {
		if i%2 == 0 {
			key = fmt.Sprintf("%v", msg)
			if isError(key) {
				key = "err"
			} else if isMethod(key) {
				key = "method"
			} else {
				keys = append(keys, key)
			}
		} else {
			value = msg
			logMessages[key] = value
		}
	}
	otherFields := make(map[string]interface{})
	payload := make(map[string]interface{})
	result := make(map[string]interface{})
	err := getErrorValue(logMessages["err"])
	method := fmt.Sprintf("%v", logMessages["method"])

	for _, k := range keys {
		formatedKey := formatedName(k)
		if isResult(k) {
			if json, ok := structNameAndIDJson(logMessages[k]); ok {
				result[formatedKey] = json
			} else if json, ok := iterableNameAndLengthJson(logMessages[k]); ok {
				result[formatedKey] = json
			} else {
				result[formatedKey] = logMessages[k]
			}

		} else if isPayload(k) {
			if isStruct(logMessages[k]) {
				logMessages[k] = Filter(logMessages[k])
			}
			payload[formatedKey] = logMessages[k]
		} else {
			otherFields[formatedKey] = logMessages[k]
		}
	}
	fmt.Println("=============")

	logRequest(method, payload, startTime, result, err)
}

func logRequest(method string, payload interface{}, startTime time.Time, result interface{}, err error) {
	// duration := time.Since(startTime)
	logger.WithFields(logrus.Fields{
		"timestamp":  time.Now().Format(time.RFC3339),
		"level":      getLoggerLevel(err),
		"method":     method,
		"payload":    payload,
		"request_at": startTime.Format(time.RFC3339),
		// "duration_time": duration.String(),
		"result": result,
		"error":  formatError(err),
	}).Info("Processed request")
}

func getLoggerLevel(err error) string {
	if err != nil {
		return "error"
	}
	return "info"
}

func formatError(err error) interface{} {
	if err != nil {
		return map[string]string{
			"message":     err.Error(),
			"stack_trace": string(debug.Stack()),
		}
	}
	return nil
}

// func extractTextFromJSON(jsonStr string) (string, error) {
// 	var data struct {
// 		Root struct {
// 			Children []struct {
// 				Children []struct {
// 					Text string `json:"text"`
// 				} `json:"children"`
// 			} `json:"children"`
// 		} `json:"root"`
// 	}

// 	err := json.Unmarshal([]byte(jsonStr), &data)
// 	if err != nil {
// 		return "", err
// 	}

// 	if len(data.Root.Children) > 0 && len(data.Root.Children[0].Children) > 0 {
// 		return data.Root.Children[0].Children[0].Text, nil
// 	}
// 	return "", fmt.Errorf("no text found")
// }
