package regex

import (
	"fmt"
	"regexp"
)

var regexes = map[string]interface{}{
	// JWT token: BASE64_URL_ENC DOT BASE64_URL_ENC DOT BASE64_URL_ENC
	"jwt": `^[a-zA-Z0-9-_]+?\.[a-zA-Z0-9-_]+?\.([a-zA-Z0-9-_]+)+$`,
	// Username: starts with letter
	//           only alphanumeric characters
	//           3-32 characters
	"username": `^[A-Za-z][A-Za-z0-9]{2,32}$`,
	// Password: at least 1 letter (lowercase or uppercase)
	//           at least 1 number
	//           at least 1 special character (TODO @sht remove)
	//           8-64 characters
	"master_password": []string{
		`[A-Za-z]+`,
		"[0-9]+",
		`^.{8,64}$`,
	},
	"uuid": `^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$`,
}

// Match returns true if the given string matches the pattern(s) of the
// specified key in the regexes map. Returns an error if a pattern is invalid
func Match(key string, input string) bool {
	if pattern, ok := regexes[key]; ok {
		switch pattern.(type) {
		case string:
			// Pattern is string, check if pattern matches input
			rex, _ := regexp.Compile(pattern.(string))
			return rex.MatchString(input)
		case []string:
			// Pattern is array of strings, check if input matches all of them
			for _, pat := range pattern.([]string) {
				rex, _ := regexp.Compile(pat)
				match := rex.MatchString(input)
				if !match {
					return false
				}
			}
			return true
		default:
			return false
		}
	}
	// Pattern doesn't exist in regexes map
	return false
}

func Load() error {
	for _, pattern := range regexes {
		switch pattern.(type) {
		case string:
			// Pattern is string, check if pattern matches input
			_, err := regexp.Compile(pattern.(string))
			if err != nil {
				return err
			}
		case []string:
			// Pattern is array of strings, check if input matches all of them
			for _, pat := range pattern.([]string) {
				_, err := regexp.Compile(pat)
				if err != nil {
					return err
				}
			}
		default:
			return fmt.Errorf(
				"Invalid pattern format (expected string or array of strings)",
			)
		}
	}
	return nil
}
