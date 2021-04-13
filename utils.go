package main

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"os/user"
	"path"
	"regexp"
	"strings"
)

func Check(e error) {
	if e != nil {
		panic(e)
	}
}

func GetUserPath() string {
	user, err := user.Current()
	Check(err)

	return user.HomeDir
}

type Param struct {
	Name  string
	Value string
}

func CreateGETParams(params []Param) string {
	var list []string

	for _, e := range params {
		list = append(list, fmt.Sprintf("%s=%s", url.QueryEscape(e.Name), url.QueryEscape(e.Value)))
	}

	return fmt.Sprintf("?%s", strings.Join(list, "&"))
}

func CreateFolderIfNotExists(folder string) {
	folderPath := path.Join(GetUserPath(), folder)
	_, err := os.Stat(folderPath)

	if os.IsNotExist(err) {
		err := os.Mkdir(folderPath, 0700)
		Check(err)
	}
}

func SaveToFile(folder string, filename string, data []byte) {
	folderPath := path.Join(GetUserPath(), folder)
	filePath := path.Join(folderPath, filename)

	_, err := os.Stat(folderPath)

	if os.IsNotExist(err) {
		CreateFolderIfNotExists(folder)
	}

	err = ioutil.WriteFile(filePath, data, 0600)
	Check(err)
}

func SafeFileName(s string, r string) string {
	re := regexp.MustCompile(`[^0-9A-Za-z]`)
	safe := re.ReplaceAllString(s, r)

	return safe
}

func ParamContainsNested(p []Param, key string, val string) bool {
	for _, e := range p {
		if (key == "Name" && e.Name == val) || (key == "Value" && e.Value == val) {
			return true
		}
	}

	return false
}
