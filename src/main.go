package main

import (
	"os"

	"./constant"
	"./ui"
)

func usage(exitValue int, message string) {

	var execName string = os.Args[0]

	if message != "" {
		println("MESSAGE: " + message)
	}
	println("USAGE:")
	println("\t" + execName + " NES_ROM_PATH")
	println("NES_ROM_PATH " + "the path of your nes game")
	os.Exit(exitValue)
}

func main() {
	var argc int = len(os.Args)

	if argc != 2 {
		usage(constant.ExitFailure, "not enought arguments")
	}
	if !ui.Start(os.Args[1]) {
		usage(constant.ExitFailure, "Execution error")
	}
}
