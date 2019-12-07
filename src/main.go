package main

// func main() {
// 	var a uint = 60 /* 60 = 0011 1100 */
// 	var b uint = 13 /* 13 = 0000 1101 */
// 	var c uint = 0

// 	c = a & b /* 12 = 0000 1100 */
// 	fmt.Printf("Line 1 - Value of c is %d\n", c)

// 	c |= b /* 61 = 0011 1101 */
// 	c &= 3
// 	fmt.Printf("Line 2 - Value of c is %d\n", c)

// 	c = a ^ b /* 49 = 0011 0001 */
// 	fmt.Printf("Line 3 - Value of c is %d\n", c)

// 	c = a << 2 /* 240 = 1111 0000 */
// 	fmt.Printf("Line 4 - Value of c is %d\n", c)

// 	c = a >> 2 /* 15 = 0000 1111 */
// 	fmt.Printf("Line 5 - Value of c is %d\n", c)
// }

import (
	"os"

	"./constant"
)

func usage(exitValue int, message string) {

	var execName string = os.Args[0]

	if message != "" {
		println("MESSAGE: " + message)
	}
	println("USAGE:")
	println("\t" + execName + " NES_ROM_PATH")
	os.Exit(exitValue)
}

func main() {
	var argc int = len(os.Args)

	if argc != 2 {
		usage(constant.ExitFailure, "not enought arguments")
	}
}
