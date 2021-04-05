package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var (
	client = http.Client{Timeout: time.Duration(5 * time.Second)}
)

func newRequest(method string, path string, body *bytes.Buffer) *http.Request {
	token := getToken()

	req, err := http.NewRequest("GET", baseURL+path, body)
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Harvest-Account-Id", *account)
	req.Header.Set("User-Agent", "THI "+version)

	return req
}

func getUserInfo() string {
	req := newRequest("GET", "users/me", bytes.NewBuffer(nil))

	res, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}

	return string(body)
}
