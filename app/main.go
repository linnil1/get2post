package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nqd/flat"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

func retrieveAndDel(key string, queryParams *url.Values) (string, error) {
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

func PostJson(targetURL string, jsonData []byte) (string, error) {
	// Send POST request to URL
	resp, err := http.Post(targetURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", errors.New("Failed to send POST request")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New("Failed to read POST response")
	}
	return fmt.Sprintf("POST response: %s", body), nil
}

func main() {
	r := gin.Default()

	r.GET("/get2post", func(c *gin.Context) {
		// Get all querys
		queryParams := c.Request.URL.Query()

		// Get app secret from env
		appSecret := os.Getenv("APP_SECRET")
		if appSecret != "" {
			secret, _ := retrieveAndDel("secret", &queryParams)
			if appSecret != secret {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid secret"})
				return
			}
		}

		// Get URL from parameters
		targetURL, err := retrieveAndDel("url", &queryParams)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Extract Data
		queryParamsDict, err := ExtractDataFromParams(&queryParams)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// map -> json string
		jsonData, err := json.Marshal(queryParamsDict)
		if err != nil {
			err = errors.New("Failed to convert to JSON")
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Send a POST request with JSON data
		message, err := PostJson(targetURL, jsonData)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": message})
	})

	r.Run()
}
