package api

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"

	"myst/pkg/crypto"
	"myst/pkg/regex"
	"myst/pkg/user"
)

// "$argon2id$v=19$m=262144,t=10,p=2$ny7MyNZJ5OMSWyWBIOGV4g$U32rqke4W3y3uBlM+bF/2MfBYZC3dm9Z8F6YquPoUtY"

var (
	jwtCookieLifetime = 604800
)

// LoginHandler handles login requests and throttles them
func LoginHandler(c *gin.Context) {
	var data struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	err := c.ShouldBind(&data)
	if err != nil {
		Error(c, 400, nil)
		return
	}
	// if !regex.Match("username", data.Username) {
	// 	Invalidate(
	// 		c,
	// 		"username",
	// 		"Usermame can only contain letters and numbers",
	// 	)
	// }
	//if !regex.Match("password", data.Password) {
	//Invalidate(
	//	c,
	//	"password",
	//	"Password can contain letters numbers and special characters",
	//)
	//}
	//if !DataValid(c) {
	//	return
	//}

	//debugHash

	u, err := user.Get("username", data.Username)
	if err == user.ErrNotFound {
		Error(c, 404, nil)
		return
	} else if err != nil {
		log.Error(err)
		Error(c, 500, nil)
		return
	}

	match, err := crypto.VerifyPassword(data.Password, u.PasswordHash)
	if err != nil {
		log.Error(err)
		Error(c, 500, nil)
		return
	} else if !match {
		Error(c, 401, nil)
		return
	}

	err = loginUser(c, u.ID)
	if err != nil {
		log.Error(err)
		Error(c, 500, nil)
		return
	}
}

// RegisterHandler creates a new user
func RegisterHandler(c *gin.Context) {
	var data struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	err := c.ShouldBind(&data)
	if err != nil {
		Error(c, http.StatusBadRequest, nil)
		return
	}

	_, err = user.Get("username", data.Username)
	if err == nil {
		Error(c, http.StatusForbidden, nil)
		return
	} else if err != user.ErrNotFound {
		log.Error(err)
		Error(c, http.StatusInternalServerError, nil)
		return
	}

	u, err := user.New(data.Username, data.Password)
	if err != nil {
		log.Error(err)
		Error(c, http.StatusInternalServerError, nil)
		return
	}

	err = loginUser(c, u.ID)
	if err != nil {
		log.Error(err)
		Error(c, http.StatusInternalServerError, nil)
		return
	}
}

func loginUser(c *gin.Context, userID string) error {
	exp := time.Now().Unix() + int64(jwtCookieLifetime)
	iat := time.Now().Unix()
	nbf := time.Now().Unix()

	// if password matches hash
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256, jwt.MapClaims{
			"exp": exp,
			"iat": iat,
			"nbf": nbf,
			"usr": userID,
		},
	)
	key, err := base64.StdEncoding.DecodeString(jwtSecretKey)
	if err != nil {
		return err
	}
	enc, err := token.SignedString(key)
	if err != nil {
		return err
	}

	c.SetCookie(
		"auth_token",
		enc,
		jwtCookieLifetime,
		"/",
		".localhost", // TODO: set to keystore name on prod, .localhost on devel
		false,        // TODO: true on production, false on development
		true,
	)

	Success(c, enc)

	return nil
}

// VerifyLogin authenticates a client
func VerifyLogin(c *gin.Context) (bool, error) {
	auth := c.GetHeader("authorization")
	if auth != "" {
		// Remove the "Bearer" prefix
		parts := strings.Split(auth, "Bearer")
		if len(parts) != 2 {
			return false, fmt.Errorf("invalid authorization header format")
		}
		auth = strings.TrimSpace(parts[1])
	} else {
		cookie, err := c.Cookie("auth_token")
		if err == nil {
			auth = cookie
		}
	}

	if auth == "" {
		return false, fmt.Errorf("authorization required")
	}

	// Validate token format
	match := regex.Match("jwt", auth)
	if !match {
		return false, fmt.Errorf("invalid JSON web token format")
	}

	// Check if authentication token is in the valid format
	token, err := jwt.Parse(
		auth, func(token *jwt.Token) (interface{}, error) {
			// Verify signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected token signing method: %v", token.Header["alg"])
			}
			key, err := base64.StdEncoding.DecodeString(jwtSecretKey)
			if err != nil {
				return nil, fmt.Errorf("jwt secret decode failed")
			}
			// return the secret when token is valid format
			return key, nil
		},
	)

	if err != nil {
		return false, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return false, fmt.Errorf("could not parse JWT claims")
	}

	id, ok := claims["usr"].(string)
	if !ok {
		return false, fmt.Errorf("invalid usr data")
	}

	u, err := user.Get("username", id)
	if err != nil {
		return false, fmt.Errorf("user not found")
	}

	err = claims.Valid()
	if err != nil {
		return false, err
	}

	c.Set("user_id", u.ID)
	return true, nil
}
