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
	// step1: get the json string from get2post response
	appJson := make(map[string]interface{})
	err := json.Unmarshal(bodyByte, &appJson)
	assert.Nil(t, err)
	messageStr, ok := appJson["message"].(string)
	assert.Equal(t, ok, true)

	// step2: json string to json and get what post data is
	httpbinJson := make(map[string]interface{})
	err = json.Unmarshal([]byte(messageStr), &httpbinJson)
	assert.Nil(t, err)
	data, ok := httpbinJson["json"].(map[string]interface{})
	assert.Equal(t, ok, true)
	return data
}

func TestBasic(t *testing.T) {
	// Testing under basic usage
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
			400,
			nil,
		},
		{
			// password no provided
			"password",
			"&data.key1=value1",
			400,
			nil,
		},
		{
			// password error
			"password",
			"&secret=pass&data.key1=value1",
			400,
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

	// request for each case
	for _, test := range tests {
		router := SetupRouter(test.password)
		record := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/get2post?url=https://httpbun.org/post"+test.url, nil)
		router.ServeHTTP(record, req)
		assert.Equal(t, test.status, record.Code)

		// if 200, compare the json with answer
		if test.status == 200 {
			ans, _ := json.Marshal(test.ansJson)
			out, err := json.Marshal(ExtractJson(t, record.Body.Bytes()))
			assert.Nil(t, err)
			assert.Equal(t, string(ans), string(out))
		}
	}
}
