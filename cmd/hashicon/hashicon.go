package main

import (
	"encoding/base64"
	"fmt"
	"os"
)

var key1, _ = base64.StdEncoding.DecodeString("TN8aNobbU/IOUWcA3Z7W1jyhjB97QcKfFnctG212vHk=")
var key2, _ = base64.StdEncoding.DecodeString("whgJQU6U7OdqXzTij2wER4i262s2/1UzYc9tmwCI49U=")
var key3, _ = base64.StdEncoding.DecodeString("QVd4WN13FmujaLEeBQ47s/A7zW7D2dH5HIlQXWpTFlU=")
var key4, _ = base64.StdEncoding.DecodeString("SX8FN4PUuzURE0XEzyCgdkEGrTjWmeW1+UuRVNT0M10=")
var key5, _ = base64.StdEncoding.DecodeString("H6ayCca0tVi1hag+uo0saDwvuBV+MqiPaiiAUl4oxV0=")

func main() {
	os.MkdirAll("tmp/hashicons", os.ModePerm)
	svg("key1", hash(key1))
	svg("key2", hash(key2))
	svg("key3", hash(key3))
	svg("key4", hash(key4))
	svg("key5", hash(key5))
}

func hash(b []byte) []uint8 {
	// 16 chars, 256 / 16 so 16 chars in HEX
	pix := make([]uint8, 8*8)
	for i, r := range b {
		a := r & 0b00001111 / 4
		b := r & 0b11110000 >> 4 / 4
		pix[i*2] = a
		pix[i*2+1] = b
	}

	//for i, p := range pix {
	//	fmt.Printf("%s ", pixelChar(p))
	//	if i%8 == 7 {
	//		fmt.Println()
	//	}
	//}
	//fmt.Println()
	return pix
}

func pixelChar(p uint8) string {
	switch p {
	case 0:
		return "."
	case 1:
		return "."
	case 2:
		return "░"
	case 3:
		return "█"
	case 4:
		return "█"
	default:
		return " "
	}
}

func pixelColor(p uint8) string {
	switch p {
	default:
		return "#000000"
	case 1:
		return "#3a8d99"
	case 2:
		return "#20e9b7"
	case 3:
		return "#20e9b7"
	}
}

func getBit(b byte, i int) uint8 {
	if i < 0 || i > 7 {
		return 0
	}
	return b >> (8 - i - 1) & 1
}

func svg(name string, pix []uint8) {
	f, _ := os.Create(fmt.Sprintf("./tmp/hashicons/hashicon-%s.svg", name))
	f.Write([]byte(`<svg width="256" height="256" version="1.1" xmlns="http://www.w3.org/2000/svg">`))
	for i, p := range pix {
		x := i % 8
		y := i / 8
		f.Write([]byte(fmt.Sprintf(`<rect x="%d" y="%d" width="32" height="32" fill="%s" opacity="%f" />`, x*32, y*32, pixelColor(p), float32(p)/3)))
	}

	f.Write([]byte(`</svg>`))
	f.Close()
}
