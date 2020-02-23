// +build js,wasm

package main

//go:generate cp $GOROOT/misc/wasm/wasm_exec.js .

import (
	"fmt"
	"syscall/js"
)


func main() {
	// document := js.Global().Get("document")
	// canvas := document.Call("getElementById", "canvas")
	// fmt.Println("from wasm:")
	// fmt.Println(document)

	js.Global().Set("test", js.FuncOf(test))

	select {}
}

func test(this js.Value, inputs []js.Value) interface{}  {
	fmt.Println("test invoked")
	return nil
}
