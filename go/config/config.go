package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

var Debug bool

func Load() error {
	Debug = false
	err1 := godotenv.Overload(".env.production")
	err2 := godotenv.Overload(".env.development")
	if Get("DEBUG") == "true" {
		Debug = true
	}
	if err1 != nil && err2 != nil {
		return fmt.Errorf("no environment files found")
	}
	return nil
}

func Get(key string) string {
	return os.Getenv("MIST_" + key)
}
