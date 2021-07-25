package api

import (
	"encoding/base64"

	keystore2 "myst/cmd/server/keystore"

	"myst/user"
	"myst/userkey"

	"github.com/gin-gonic/gin"
	"github.com/sanity-io/litter"
)

var (
	testKey1, _     = base64.StdEncoding.DecodeString("ECoq1t6TVSpEkc27hmGKkYsE3m46Xxo9buRTPyu08lTh229cS77LGhwkFpNws4Ugw+MsRVSEvJ4Eb7djVacYH9gIpY7zpNTKABGaz8UTjCXlRqR/7QZVXDr0atNUDr43")
	testKeystore, _ = base64.StdEncoding.DecodeString("8ZgbIT5UyF9R2P7vYrZ6u5k6IB0iuPfyMV27GBj25ikGoBEycte8LaqMBgrIYIFG15Z/8zRVAUmyDj3zd1WwbIDsrkOfT0hUECOvc65j9VI=")
)

func InitData(c *gin.Context) {
	u1, err := user.New("rdnt", "1234567890")
	if err != nil {
		panic(err)
	}

	u2, err := user.New("test", "0987654321")
	if err != nil {
		panic(err)
	}

	s1, err := keystore2.New("default", testKeystore)
	if err != nil {
		panic(err)
	}

	k1, err := userkey.New(u1.ID, s1.ID, testKey1)
	if err != nil {
		panic(err)
	}

	s2, err := keystore2.New("secondary", testKeystore)
	if err != nil {
		panic(err)
	}

	k2, err := userkey.New(u2.ID, s2.ID, testKey1)
	if err != nil {
		panic(err)
	}

	litter.Dump(u1, u2, s1, k1, s2, k2)
}
