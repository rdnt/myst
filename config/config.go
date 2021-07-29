package config

import (
	"github.com/joho/godotenv"
	"os"
)

// Debug indicates if the process should show detailed debug information on the
// logs/terminal.
var Debug bool

func init() {
	Debug = false
	_ = godotenv.Overload(".env.production")
	_ = godotenv.Overload(".env.development")
	if os.Getenv("MYST_DEBUG") == "true" {
		Debug = true
	}
}
