package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func sendRequest(url string, data interface{}, method string, token string) (string, error) {
	//
	b, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	body := bytes.NewBuffer([]byte(b))

	client := &http.Client{}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json;charset=utf-8")

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	result, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return "", err
	}

	return string(result), nil
}
