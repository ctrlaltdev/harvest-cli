package main

import (
	"flag"
	"fmt"
	"path"
)

var (
	version = "v1.0.0"
	baseURL = "https://api.harvestapp.com/v2/"

	userHome = getUserPath()

	tokenPath   = flag.String("tokenPath", path.Join(userHome, ".harvest/token"), "Path to file containing Harvest Personal Token")
	accountPath = flag.String("accountPath", path.Join(userHome, ".harvest/account"), "Path to file containing the Default Account ID to use")

	token   = flag.String("token", getToken(tokenPath), "Harvest Personal Token")
	account = flag.String("account", getAccount(accountPath), "Account ID")
)

func main() {
	printHeader()
	flag.Parse()

	if *token == "" {
		fmt.Print("Harverst Personal Token: ")
		fmt.Scanln(token)
		saveToFile("token", *token)
	}

	if *account == "" {
		fmt.Print("Account ID: ")
		fmt.Scanln(account)
		saveToFile("account", *account)
	}

	userinfo := getUserInfo()
	fmt.Println(userinfo)
}
