package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
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

func ParamsToJson(queryParams *url.Values) ([]byte, error) {
	// Transform parameters into json
	queryParamsShrink := map[string]string{}
	for k, v := range *queryParams {
		queryParamsShrink[k] = v[0]
	}
	jsonData, err := json.Marshal(queryParamsShrink)
	if err != nil {
		return []byte{}, errors.New("Failed to convert to JSON")
	}
	return jsonData, nil
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

		// To json
		jsonData, err := ParamsToJson(&queryParams)
		if err != nil {
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
