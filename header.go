package main

import (
	"fmt"
)

func printHeader() {

	headers := []string{
		"",
		" _   _                           _   ",
		"| | | | __ _ _ ____   _____  ___| |_ ",
		"| |_| |/ _` | '__\\ \\ / / _ \\/ __| __|",
		"|  _  | (_| | |   \\ V /  __/\\__ \\ |_ ",
		"|_| |_|\\__,_|_|    \\_/ \\___||___/\\__|",
		"",
		"",
	}

	for _, l := range headers {
		fmt.Println(l)
	}
}
