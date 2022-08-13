package middleware

import (
	"echo_sample/config"
	"encoding/base64"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
)

var TokenValidationMinutes = 60 * 60 * 2
var RefreshValidationMinutes = 60 * 60 * 24 * 365

func CreateToken(userid int64, duration int) (string, error) {
	var err error
	//Creating Access Token
	secret, _ := base64.StdEncoding.DecodeString(config.Config("SECRET"))
	atClaims := jwt.MapClaims{}
	atClaims, err = createJwt(strconv.FormatInt(userid, 10), time.Duration(duration), atClaims)
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString(secret)
	if err != nil {
		return "", err
	}
	return token, nil
}

/* Private Functions */

func createJwt(subject string, expiration time.Duration, atclaims jwt.MapClaims) (jwt.MapClaims, error) {
	if subject != "" {
		atclaims["id"] = subject
	}
	atclaims["iat"] = time.Now().Unix()
	atclaims["exp"] = time.Now().Add(time.Second * expiration).Unix()

	return atclaims, nil
}
