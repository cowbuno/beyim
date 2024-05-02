package main

import (
	"bytes"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"testing"
)

type MockLogger struct {
	Buffer *bytes.Buffer
}

func (m *MockLogger) Info(args ...any) {
	m.Buffer.WriteString(fmt.Sprintln(args...))
}

func (m *MockLogger) Error(args ...any) {
	m.Buffer.WriteString(fmt.Sprintln(args...))
}

func TestLogPrint(t *testing.T) {
	var buf bytes.Buffer
	log = slog.New(
		slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelInfo}),
	)
	tests := []struct {
		name      string
		keyvals   []any
		wantLines []string
	}{
		{
			name: "simple string",
			keyvals: []any{
				"Name", "simple_string",
				"payload1:request", "request",
				"result1:result", "result",
				"result2:err", nil,
			},
			wantLines: []string{
				"\"payload1\":\"request\",\"err\":null,\"result1\":\"result\"",
				"\"level\":\"INFO\"",
			},
		}, {
			name: "err is bad",
			keyvals: []any{
				"Name", "err_is_bad",
				"payload1:request", "request",
				"result1:result", "result",
				"result2:err", errors.New("bad"),
			},
			wantLines: []string{
				"\"payload1\":\"request\",\"err\":\"bad\",\"result1\":\"result\"",
				"\"level\":\"ERROR\"",
			},
		}, {
			name: "take id from Person",
			keyvals: []any{
				"Name", "take_id_from_Person",
				"payload1:request", "request",
				"result1:result", Person{Id: 3},
				"result2:err", nil,
			},
			wantLines: []string{
				"\"payload1\":\"request\",\"err\":null,\"result1\":\"Person.id(3)\"",
				"\"level\":\"INFO\"",
			},
		}, {
			name: "take id from Person pointers",
			keyvals: []any{
				"Name", "take_id_from_Person_pointers",
				"payload1:request", "request",
				"result1:result", &Person{Id: 3},
				"result2:err", nil,
			},
			wantLines: []string{
				"\"payload1\":\"request\",\"err\":null,\"result1\":\"Person.id(3)\"",
				"\"level\":\"INFO\"",
			},
		}, {
			name: "take len of array",
			keyvals: []any{
				"Name", "take_len_of_data",
				"payload1:request", "request",
				"result1:result", &Person{Id: 3},
				"result2:data", []int{1, 2, 3, 4, 5},
				"result2:err", nil,
			},
			wantLines: []string{
				"\"payload1\":\"request\",\"err\":null,\"result1\":\"Person.id(3)\",\"result2\":\"[]int.len(5)\"",
				"\"level\":\"INFO\"",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf.Reset()
			logPrint(tt.keyvals...)
			for _, wantLine := range tt.wantLines {
				if !strings.Contains(buf.String(), wantLine) {
					t.Errorf("logPrint() = %v, want %v", buf.String(), wantLine)
				}
			}
		})
	}
}
