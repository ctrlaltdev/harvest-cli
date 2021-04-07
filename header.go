package main

import (
	"fmt"
	"math/rand"
	"time"
)

func printHeader() {
	headers := [][8]string{}

	h1 := [...]string{
		"  _______ _    _ _____ ",
		" |__   __| |  | |_   _|",
		"    | |  | |__| | | |  ",
		"    | |  |  __  | | |  ",
		"    | |  | |  | |_| |_ ",
		"    |_|  |_|  |_|_____|",
		"",
		"",
	}

	h2 := [...]string{
		"",
		" ████████╗██╗  ██╗██╗",
		" ╚══██╔══╝██║  ██║██║",
		"    ██║   ███████║██║",
		"    ██║   ██╔══██║██║",
		"    ██║   ██║  ██║██║",
		"    ╚═╝   ╚═╝  ╚═╝╚═╝",
		"",
	}

	h3 := [...]string{
		"",
		" _|_|_|_|_|  _|    _|  _|_|_|  ",
		"     _|      _|    _|    _|    ",
		"     _|      _|_|_|_|    _|    ",
		"     _|      _|    _|    _|    ",
		"     _|      _|    _|  _|_|_|  ",
		"",
		"",
	}

	headers = append(headers, h1)
	headers = append(headers, h2)
	headers = append(headers, h3)

	rand.Seed(time.Now().UnixNano())
	i := rand.Intn(len(headers)) // #nosec G404

	for _, l := range headers[i] {
		fmt.Println(l)
	}
}
