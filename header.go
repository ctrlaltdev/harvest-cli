package main

import (
	"fmt"
	"strings"

	te "github.com/muesli/termenv"
)

func printHeader() {
	header := []string{
		"  _   _                           _   ",
		" | | | | __ _ _ ____   _____  ___| |_ ",
		" | |_| |/ _\\ | '__\\ \\ / / _ \\/ __| __|",
		" |  _  | (_| | |   \\ V /  __/\\__ \\ |_ ",
		" |_| |_|\\__,_|_|    \\_/ \\___||___/\\__|",
		"                                      ",
	}

	fmt.Printf("%s\n\n", te.String(strings.Join(header, "\n")).Foreground(color("202")).String())

}
