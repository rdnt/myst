// +build js,wasm

package main

//go:generate cp $GOROOT/misc/wasm/wasm_exec.js .

import (
	"encoding/hex"
	"myst/tmp"
	"syscall/js"
)

func main() {
	js.Global().Set("keystore", js.ValueOf(map[string]interface{}{
		"_decrypt": js.FuncOf(Decrypt),
		"_encrypt": js.FuncOf(Encrypt),
	}))
	// js.Global().Set("decryptKeystore", js.FuncOf(Decrypt))
	// js.Global().Set("encryptKeystore", js.FuncOf(Encrypt))
	select {}
}

func Decrypt(this js.Value, args []js.Value) interface{} {
	if len(args) == 3 {
		raw := []byte(args[0].String())
		pass := args[1].String()
		cb := args[2]
		b := make([]byte, hex.DecodedLen(len(raw)))
		_, err := hex.Decode(b, raw)
		if err != nil {
			cb.Invoke(nil, err.Error())
			return nil
		}
		json, err := tmp.Decrypt(b, []byte(pass))
		if err != nil {
			cb.Invoke(nil, err.Error())
			return nil
		}
		cb.Invoke(string(json), nil)
		return nil
	}
	return false
}

func Encrypt(this js.Value, args []js.Value) interface{} {
	if len(args) == 3 {
		json := args[0].String()
		pass := args[1].String()
		cb := args[2]
		b, err := tmp.Encrypt([]byte(json), []byte(pass))
		if err != nil {
			cb.Invoke(nil, err.Error())
			return nil
		}
		h := make([]byte, hex.EncodedLen(len(b)))
		hex.Encode(h, b)
		cb.Invoke(string(h), nil)
		return nil
	}
	return false
}
