package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/nqd/flat"
	"io"
	"net/http"
	"net/url"
	"os"
)

func RetrieveAndDel(key string, queryParams *url.Values) (string, error) {
	// Get key from parameters
	value := queryParams.Get(key)
	if value == "" {
		return "", errors.New(key + " not provided")
	}
	queryParams.Del(key)
	return value, nil
}

func ExtractDataFromParams(queryParams *url.Values) (map[string]interface{}, error) {
	// Transform parameters into map
	// Step1: url.Values -> map[string]string{}
	queryParamsShrink := make(map[string]interface{})
	for k, v := range *queryParams {
		queryParamsShrink[k] = v[0]
	}

	// step2: unflatten
	queryParamsDict, err := flat.Unflatten(queryParamsShrink, nil)
	if err != nil {
		return nil, errors.New("Failed to unflatten")
	}

	// step3: Extract 'data' key
	data, ok := queryParamsDict["data"]
	if !ok {
		return nil, errors.New("data not provided")
	}
	return data.(map[string]interface{}), nil
}

func PostJson(targetURL string, jsonData []byte) (int, string, error) {
	// Send POST request to URL
	resp, err := http.Post(targetURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return resp.StatusCode, "", errors.New("Failed to send POST request")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, "", errors.New("Failed to read POST response")
	}
	return resp.StatusCode, string(body), nil
}

func SetupRouter(appSecret string) *gin.Engine {
	// Main function
	r := gin.Default()

	r.GET("/get2post", func(c *gin.Context) {
		// Get all querys
		queryParams := c.Request.URL.Query()

		// Get app secret from env
		if appSecret != "" {
			secret, _ := RetrieveAndDel("secret", &queryParams)
			if appSecret != secret {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid secret"})
				return
			}
		}

		// Get URL from parameters
		targetURL, err := RetrieveAndDel("url", &queryParams)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Extract Data
		queryParamsDict, err := ExtractDataFromParams(&queryParams)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// map -> json string
		jsonData, err := json.Marshal(queryParamsDict)
		if err != nil {
			err = errors.New("Failed to convert to JSON")
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Send a POST request with JSON data
		status, message, err := PostJson(targetURL, jsonData)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": message, "status": status})
	})
	return r
}

func main() {
	appSecret := os.Getenv("APP_SECRET")
	r := SetupRouter(appSecret)
	_ = r.Run() // port: 8080
}
