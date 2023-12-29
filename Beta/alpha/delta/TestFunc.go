package main

import (
	"fmt"
	"syscall/js"
)

func TestUserInput(input string) string {
	return "User Input: " + input
}

func Main(this js.Value, args []js.Value) interface{} {
	if len(args) > 0 {
		input := args[0].String() // Extract input from JavaScript context
		return TestUserInput(input)
	}
	return "Error: No arguments provided."
}

func main() {
	fmt.Println("Setting Main function...")
	js.Global().Set("Main", js.FuncOf(Main))
	fmt.Println("Main function set.")

	// Prevent the Go program from exiting
	c := make(chan struct{}, 0)
	<-c
}
