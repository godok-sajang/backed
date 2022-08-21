package dto

import (
	"database/sql"
	"echo_sample/util"
	"net/mail"
	"time"
	"unicode"

	"github.com/golang-jwt/jwt"
)

type UserInfoRequest struct {
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Birth    string `json:"birth"`
	Gender   int    `json:"gender"`
}

func (uir UserInfoRequest) ValidateNickname() bool {
	if !(len(uir.Nickname) >= 3 && len(uir.Nickname) <= 15) {
		return false
	}
	return true
}

func (uir UserInfoRequest) ValidateEmail() bool {
	_, err := mail.ParseAddress(uir.Email)
	return err == nil
}

func (uir UserInfoRequest) ValidatePassword() bool {
	if len(uir.Password) < 8 || len(uir.Password) > 50 {
		return false
	}
	var number, letter, special bool
	for _, c := range uir.Password {
		switch {
		case unicode.IsNumber(c):
			number = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c) || c == ' ':
			special = true
		case unicode.IsLetter(c):
			letter = true
		}
	}
	return (number && letter) || (number && special) || (letter && special)
}

func (uir UserInfoRequest) ToDtoUserInfo() UserInfo {
	return UserInfo{
		UserID:   0,
		Nickname: uir.Nickname,
		Email:    uir.Email,
		Password: util.HashPassword(uir.Password),
		Birth:    uir.BirthToTimestamp(),
		Gender:   uir.Gender,
	}
}

func (uir UserInfoRequest) BirthToTimestamp() time.Time {
	birth, err := time.Parse(time.RFC3339, uir.Birth)
	if err != nil {
		birth, err = time.Parse("2006-01-02", uir.Birth)
	}
	return birth
}

func (uir UserInfoRequest) ValidateBirth() error {
	_, err := time.Parse(time.RFC3339, uir.Birth)
	if err != nil {
		_, err = time.Parse("2006-01-02", uir.Birth)
	}
	return err
}

type UserInfo struct {
	UserID   int       `json:"userId" bson:"user_id"`
	Nickname string    `json:"nickname" bson:"nickname"`
	Email    string    `json:"email" bson:"email"`
	Password string    `json:"password" bson:"password"`
	Birth    time.Time `json:"birth" bson:"birth"`
	Gender   int       `json:"gender" bson:"gender"`
}

type UserSignInRequest struct {
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
}

type JwtCustomClaims struct {
	UserID int  `json:"userID"`
	Admin  bool `json:"admin"`
	jwt.StandardClaims
}

type UserTokens struct {
	Success       bool           `json:"success"`
	UserID        sql.NullInt64  `json:"userID"`
	AccessToken   sql.NullString `json:"accessToken"`
	RefreshToken  sql.NullString `json:"refreshToken"`
	ExistUserInfo ExistUser      `json:"existUserInfo"`
}

type ExistUser struct {
	UserID   int    `json:"userID"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
}
