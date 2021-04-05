package main

import (
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func getUserPath() string {
	user, err := user.Current()
	check(err)

	return user.HomeDir
}

func getToken(path *string) string {
	_, err := os.Stat(*path)

	if os.IsNotExist(err) {
		return ""
	}

	data, err := ioutil.ReadFile(*tokenPath)
	check(err)

	return strings.Replace(string(data), "\n", "", -1)
}

func getAccount(path *string) string {
	_, err := os.Stat(*path)

	if os.IsNotExist(err) {
		return ""
	}

	data, err := ioutil.ReadFile(*path)
	check(err)

	return strings.Replace(string(data), "\n", "", -1)
}

func saveToFile(filename string, data string) {
	filePath := path.Join(userHome, ".harvest", filename)
	_, err := os.Stat(filePath)

	if os.IsNotExist(err) {
		err := ioutil.WriteFile(filePath, []byte(data), 0600)
		check(err)
	}
}
