package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

func fromJSON(s string) (Log, error) {
	l := &Log{}

	if err := json.Unmarshal([]byte(s), &l.Fields); err != nil {
		return *l, err
	}

	if n, ok := l.Fields["level"].(string); ok {
		l.Level = string(n)
		delete(l.Fields, "level")
	}

	if n, ok := l.Fields["timestamp"].(string); ok {
		l.Timestamp = string(n)
		delete(l.Fields, "timestamp")
	}

	if n, ok := l.Fields["msg"].(string); ok {
		l.Msg = string(n)
		delete(l.Fields, "msg")
	}

	return *l, nil
}

func toString(value interface{}) string {
	switch value.(type) {
	case string:
		return value.(string)
	case bool:
		return strconv.FormatBool(value.(bool))
	case float64:
		return fmt.Sprintf("%f", value.(float64))
	case []string:
		return strings.Join(value.([]string), "\n")
	default:
		return value.(string)
	}
}

func formatLog(l Log) string {
	level := formatLevel(l.Level)
	msg := formatMsg(l.Msg)

	parts := []string{
		"  ",
		msg,
	}

	for k, val := range l.Fields {
		parts = append(parts, formatLabel(k)+": "+toString(val))
	}

	return level + strings.Join(parts[:], " ")
}
