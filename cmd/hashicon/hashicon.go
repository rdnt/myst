package main

import (
	cryptorand "crypto/rand"
	"encoding/base64"
	"fmt"
	"math/rand"

	"myst/hashicon"
)

var key1, _ = base64.StdEncoding.DecodeString("TN8aNobbU/IOUWcA3Z7W1jyhjB97QcKfFnctG212vHk=")
var key2, _ = base64.StdEncoding.DecodeString("whgJQU6U7OdqXzTij2wER4i262s2/1UzYc9tmwCI49U=")
var key3, _ = base64.StdEncoding.DecodeString("QVd4WN13FmujaLEeBQ47s/A7zW7D2dH5HIlQXWpTFlU=")
var key4, _ = base64.StdEncoding.DecodeString("SX8FN4PUuzURE0XEzyCgdkEGrTjWmeW1+UuRVNT0M10=")
var key5, _ = base64.StdEncoding.DecodeString("H6ayCca0tVi1hag+uo0saDwvuBV+MqiPaiiAUl4oxV0=")

var k128, _ = base64.StdEncoding.DecodeString("lbEvdPpqssAoXUUPJYenXw==")
var k256, _ = base64.StdEncoding.DecodeString("TN8aNobbU/IOUWcA3Z7W1jyhjB97QcKfFnctG212vHk=")
var k384, _ = base64.StdEncoding.DecodeString("N4WzXhbNlUuNo2yvkRg+pJo/5Nf+gzJk9I77nW/MuSRq2jUkpPEbscKlOrK6BBAL")
var k512, _ = base64.StdEncoding.DecodeString("IQei1hUwwYfGUdYOH1vyTwb1oM/q7CZN1Tn75aT02M595dMjM6b1qYzPdrP44QxQCPNti7wtCbLM6DTK0NE4/A==")
var k640, _ = base64.StdEncoding.DecodeString("OF/OWdPUKAd4EswFdavFGx/zcy291WlMFzniWdVTB4fsdSYKVbZz9tLL915/8QzFmulhbSfVUfTAgbFgo1amcJXtRGodnyRGw8UbpUYbYT4=")
var k768, _ = base64.StdEncoding.DecodeString("PM2GUw3OWjRpfcCT8mnIvIkjnDlQhN1KcTi3ZiN+WCVXCVbeq7h5K4gNe14iY6dWOgoHwtwQKQEPeIjxRsygw84XJtOPtKtlEwAnHII/YxDhFfbfdCpyJdrvz/WosBfj")
var k896, _ = base64.StdEncoding.DecodeString("+FhPRY6NkjgBy7qDoxQTcoUAbjukakMvBCyE3A+B12suwDWbMtAJP9VNiRDZ8g9/WprMDPuM1oSJM/h5iw62U/AnyiLdckAgngOv2uh6PD/1rbFFLVU87R5snbU8ZVIsZ98M4h/VLpV+fIZe+bFsAA==")
var k1024, _ = base64.StdEncoding.DecodeString("IA+SNxp/wRA13dQNh8+VbEB+zbav3d7UCg9fy3d3M8iYs5difC19xF8XSntXCOXJjyomT9KJ++IdpvLmDURwgU710QC669e2gxYtVUXgsR7UCNbNRNhH51AZ+wUE4Q4RCir24tVTRgE6vMww28WhPfp2MGkUsvcrD6DF7DMT5Cg=")

func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := cryptorand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func main() {

	//size := 256
	//width := int(math.Sqrt(float64(size) / 4))
	//height := int(math.Sqrt(float64(size) / 4))
	//fmt.Println(width, height)

	//bits := 256
	//b, err := GenerateRandomBytes(bits / 8)
	//if err != nil {
	//	panic(err)
	//}
	//h, err := hashicon.New(b)
	//if err != nil {
	//	panic(err)
	//}
	//err = h.Save("test")
	//if err != nil {
	//	panic(err)
	//}

	it := 120 // 120
	sizes := []int{
		256, /*512, 1024, 2048, 4096, 8192, 16384, 32768, 65536,*/
	}

	for i := 0; i < it; i++ {
		b, _ := GenerateRandomBytes(sizes[rand.Intn(len(sizes))] / 8)
		h, _ := hashicon.New(b)
		h.Save(fmt.Sprintf("test-%d", i))
	}

	//h128, _ := hashicon.New(k128)
	//h256, _ := hashicon.New(k256)
	//h384, _ := hashicon.New(k384)
	//h512, _ := hashicon.New(k512)
	//h640, _ := hashicon.New(k640)
	//h768, _ := hashicon.New(k768)
	//h896, _ := hashicon.New(k896)
	//h1024, _ := hashicon.New(k1024)
	//
	//h128.Save("h128")
	//h256.Save("h256")
	//h384.Save("h384")
	//h512.Save("h512")
	//h640.Save("h640")
	//h768.Save("h768")
	//h896.Save("h896")
	//h1024.Save("h1024")

	//os.MkdirAll("tmp/hashicons", os.ModePerm)
	//svg("key1", hash(key1))
	//svg("key2", hash(key2))
	//svg("key3", hash(key3))
	//svg("key4", hash(key4))
	//svg("key5", hash(key5))
}

// 4 bits per pixel so / 4 of the total amount of bytes, then / 2

//func hash(b []byte) []uint8 {
//	// 16 chars, 256 / 16 so 16 chars in HEX
//	pix := make([]uint8, 8*8)
//	for i, r := range b {
//		a := r & 0b00001111 / 4
//		b := r & 0b11110000 >> 4 / 4
//		pix[i*2] = a
//		pix[i*2+1] = b
//	}
//
//	//for i, p := range pix {
//	//	fmt.Printf("%s ", pixelChar(p))
//	//	if i%8 == 7 {
//	//		fmt.Println()
//	//	}
//	//}
//	//fmt.Println()
//	return pix
//}
//
//func pixelChar(p uint8) string {
//	switch p {
//	case 0:
//		return "."
//	case 1:
//		return "."
//	case 2:
//		return "░"
//	case 3:
//		return "█"
//	case 4:
//		return "█"
//	default:
//		return " "
//	}
//}
//
//func pixelColor(p uint8) string {
//	switch p {
//	default:
//		return "#000000"
//	case 1:
//		return "#3a8d99"
//	case 2:
//		return "#20e9b7"
//	case 3:
//		return "#20e9b7"
//	}
//}
//
//func getBit(b byte, i int) uint8 {
//	if i < 0 || i > 7 {
//		return 0
//	}
//	return b >> i & 1
//}
//
//func svg(name string, pix []float32) {
//	f, _ := os.Create(fmt.Sprintf("./tmp/hashicons/hashicon-%s.svg", name))
//	f.Write([]byte(`<svg width="256" height="256" version="1.1" xmlns="http://www.w3.org/2000/svg">`))
//	for i, p := range pix {
//		x := i % 8
//		y := i / 8
//		f.Write([]byte(fmt.Sprintf(`<rect x="%d" y="%d" width="32" height="32" fill="#20e9b7" opacity="%f" />`, x*32, y*32, p)))
//	}
//
//	f.Write([]byte(`</svg>`))
//	f.Close()
//}
