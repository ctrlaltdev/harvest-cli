package main

import (
	"github.com/ctrlaltdev/harvest-cli/cmd"
	"github.com/muesli/termenv"
)

type Colors struct {
	primary   termenv.Color
	secondary termenv.Color
}

var (
	version = "v1.2.0"

	color  = termenv.ColorProfile().Color
	colors = Colors{primary: color("202")}
)

func main() {
	printHeader()

	cmd.Execute()
}
