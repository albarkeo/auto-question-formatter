package main

import (
	"fmt"
	"syscall/js"
)

func main() {
	js.Global().Set("hello", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		fmt.Println("Hello, WebAssembly!")
		return nil
	}))

	js.Global().Set("printMessage", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		message := args[0].String()
		fmt.Println(message)
		return nil
	}))

	// Prevent the Go program from exiting
	c := make(chan struct{}, 0)
	<-c
}
