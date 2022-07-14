package suite

import (
	"math/rand"
	"time"

	"myst/pkg/crypto"
)

func randomStringsold(count, size int) ([]string, error) {
	ids := make([]string, count)
	var err error

	for i := 0; i < count; i++ {
		ids[i], err = crypto.GenerateRandomString(size)
		if err != nil {
			return nil, err
		}
	}

	return ids, nil
}

var src = rand.NewSource(time.Now().UnixNano())

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func random(size ...int) string {
	n := 32
	if len(size) > 0 {
		n = size[0]
	}

	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}
