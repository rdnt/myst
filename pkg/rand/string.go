package rand

import (
	"crypto/rand"
)

func String(length int) (string, error) {
	result := make([]byte, length)

	_, err := rand.Read(result)
	if err != nil {
		return "", err
	}

	for i := 0; i < length; i++ {
		result[i] &= 0x7F

		for result[i] < 32 || result[i] == 127 {
			_, err = rand.Read(result[i : i+1])
			if err != nil {
				return "", err
			}

			result[i] &= 0x7F
		}
	}
	
	return string(result), nil
}
