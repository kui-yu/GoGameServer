package main

import (
	"bytes"
	"encoding/json"

	"io/ioutil"
	"net/http"
)

func SendRequest(url string, data interface{}, method string, token string) (string, error) {
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

	if method == "GET" {
		req.Header.Set("Content-Type", "text/plain")
	} else {
		req.Header.Set("Content-Type", "application/json;charset=utf-8")
	}

	if token != "" {
		req.Header.Set("token", token)
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

	DebugLog("SendRequest req:%s res%s\n", string(b), string(result))

	return string(result), nil
}

func SendRequestByString(url string, data string, method string, token string) (string, error) {
	body := bytes.NewBuffer([]byte(data))

	client := &http.Client{}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return "", err
	}

	if method == "GET" {
		req.Header.Set("Content-Type", "text/plain")
	} else {
		req.Header.Set("Content-Type", "application/json;charset=utf-8")
	}

	if token != "" {
		req.Header.Set("token", token)
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
