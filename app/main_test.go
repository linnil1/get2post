package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func ExtractJson(t *testing.T, bodyByte []byte) map[string]interface{} {
	// extract the real json input from response
	data := make(map[string]string)
	err := json.Unmarshal(bodyByte, &data)
	assert.Nil(t, err)
	data1 := make(map[string]interface{})
	err = json.Unmarshal([]byte(data["message"]), &data1)
	assert.Nil(t, err)
	out, ok := data1["json"].(map[string]interface{})
	assert.Equal(t, ok, true)
	return out
}

func TestBasic(t *testing.T) {
	tests := []struct {
		password string
		url      string
		status   int
		ansJson  map[string]interface{}
	}{
		{
			// base case
			"",
			"&data.key1=value1",
			200,
			map[string]interface{}{"key1": "value1"},
		},
		{
			// no data
			"",
			"",
			500,
			nil,
		},
		{
			// password no provided
			"password",
			"&data.key1=value1",
			500,
			nil,
		},
		{
			// password error
			"password",
			"&secret=pass&data.key1=value1",
			500,
			nil,
		},
		{
			// success & json is complex
			"password",
			"&secret=password&data.key1=value1&data.key2.key21=value21",
			200,
			map[string]interface{}{"key1": "value1", "key2": map[string]string{"key21": "value21"}},
		},
	}

	for _, test := range tests {
		router := setupRouter(test.password)
		record := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/get2post?url=https://httpbun.org/post"+test.url, nil)
		router.ServeHTTP(record, req)
		assert.Equal(t, test.status, record.Code)

		if test.status == 200 {
			ans, _ := json.Marshal(test.ansJson)
			out, err := json.Marshal(ExtractJson(t, record.Body.Bytes()))
			assert.Nil(t, err)
			assert.Equal(t, string(ans), string(out))
		}
	}
}
