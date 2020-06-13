package mocks

import (
	"encoding/json"
	"testing"
)

// UnmarshalJSON JSON 文字列から map に変換
func UnmarshalJSON(t *testing.T, jsonString string) map[string]interface{} {
	t.Helper()
	var data map[string]interface{}
	err := json.Unmarshal([]byte(jsonString), &data)
	if err != nil {
		t.Fatalf(err.Error())
	}
	return data
}

// MarshalJSON map から JSON 文字列に変換
func MarshalJSON(t *testing.T, jsonMap map[string]interface{}) string {
	t.Helper()
	b, err := json.Marshal(jsonMap)
	if err != nil {
		t.Fatalf(err.Error())
	}
	return string(b)
}
