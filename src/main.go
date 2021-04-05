package main

import (
	"flag"
	"fmt"
	"log"
	"path"
)

var (
	version = "v1.0.0"
	baseURL = "https://api.harvestapp.com/v2/"

	userHome = getUserPath()

	tokenPath          = flag.String("token", path.Join(userHome, ".harvest/token"), "Path to file containing Harvest Personal Token")
	defaultAccountPath = flag.String("accountPath", path.Join(userHome, ".harvest/account"), "Path to file containing the Default Account ID to use")
	account            = flag.String("account", getDefaultAccount(defaultAccountPath), "Account ID")
)

func main() {
	printHeader()
	flag.Parse()

	if *account == "" {
		log.Fatal("No Account ID provided.")
	}

	userinfo := getUserInfo()
	fmt.Println(userinfo)
}
