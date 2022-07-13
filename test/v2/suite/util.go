package suite

import (
	"myst/pkg/crypto"
)

func randomStrings(count, size int) ([]string, error) {
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
