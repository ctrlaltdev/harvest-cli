package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"time"
)

var (
	client = http.Client{Timeout: time.Duration(5 * time.Second)}
)

func newRequest(method string, path string, body *bytes.Buffer) *http.Request {
	req, err := http.NewRequest("GET", baseURL+path, body)
	check(err)

	req.Header.Set("Authorization", "Bearer "+*token)
	req.Header.Set("Harvest-Account-Id", *account)
	req.Header.Set("User-Agent", "THI "+version)

	return req
}

func getUserInfo() string {
	req := newRequest("GET", "users/me", bytes.NewBuffer(nil))

	res, err := client.Do(req)
	check(err)

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	check(err)

	return string(body)
}
