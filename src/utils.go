package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"strings"
)

func getUserPath() string {
	user, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	return user.HomeDir
}

func getToken() string {
	data, err := ioutil.ReadFile(*tokenPath)
	if err != nil {
		log.Fatal(err)
	}

	return strings.Replace(string(data), "\n", "", -1)
}

func getDefaultAccount(path *string) string {
	_, err := os.Stat(*path)

	if os.IsNotExist(err) {
		return ""
	}

	data, err := ioutil.ReadFile(*path)
	if err != nil {
		log.Fatal(err)
	}

	return strings.Replace(string(data), "\n", "", -1)
}
