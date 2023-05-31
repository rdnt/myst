package config

import (
	"os"

	"github.com/joho/godotenv"
)

// Debug indicates if the process should show detailed debug information on the
// logs/terminal.
var Debug bool = true

func init() {
	_ = godotenv.Overload(".env.production")
	_ = godotenv.Overload(".env.development")
	if os.Getenv("MYST_DEBUG") == "true" {
		Debug = true
	}
}
