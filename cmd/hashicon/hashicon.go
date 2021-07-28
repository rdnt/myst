package main

import (
	cryptorand "crypto/rand"
	"encoding/base64"
	"fmt"
	"math/rand"

	"myst/hashicon"
)

var k256, _ = base64.StdEncoding.DecodeString("TN8aNobbU/IOUWcA3Z7W1jyhjB97QcKfFnctG212vHk=")

func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := cryptorand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func main() {
	testBulk()
	test256()
}

func testBulk() {
	count := 64
	bitSizes := []int{
		256, 512, 1024, 2048, 4096, 8192, 16384, 32768, 65536,
	}

	for i := 0; i < count; i++ {
		b, _ := GenerateRandomBytes(bitSizes[rand.Intn(len(bitSizes))] / 8)
		h, _ := hashicon.New(b)
		name := fmt.Sprintf("tmp/hashicons/test-%d-%d.svg", i, len(b)*8)
		h.Export(name)
		fmt.Println(name)
	}
}

func test256() {
	name := "tmp/hashicons/test.svg"
	h, _ := hashicon.New(k256)
	h.Export(name)
	fmt.Println(name)
}
