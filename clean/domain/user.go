package domain

import (
	"database/sql"
	"fmt"
	"godok/util"
	"net/mail"
	"time"
	"unicode"

	"github.com/golang-jwt/jwt"
)

type CheckAuthResponse struct {
	UserID   int    `json:"userId" bson:"user_id"`
	Nickname string `json:"nickname" bson:"nickname"`
	Email    string `json:"email" bson:"email"`
}

type UserInfoRequest struct {
	ID       int
	Nickname string `json:"nickname" form:"nickname"`
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
	Birth    string `json:"birth" form:"birth"`
	Gender   int    `json:"gender" form:"gender"`
}

func (uir UserInfoRequest) GetQueryByNickname() string {
	query := ""

	if uir.Nickname != "" {
		query += fmt.Sprintf("and nickname=%v", uir.Nickname)
	}

	return query
}

func (uir UserInfoRequest) GetQueryByEmail() string {
	query := ""

	if uir.Email != "" {
		query += fmt.Sprintf("and email='%v'", uir.Email)
	}

	return query
}

func (uir UserInfoRequest) GetQueryByPassword() string {
	query := ""

	if uir.Password != "" {
		query += fmt.Sprintf("and password='%v'", uir.Password)
	}

	return query
}

func (uir UserInfoRequest) ValidateNickname() bool {
	if !(len(uir.Nickname) >= 3 && len(uir.Nickname) <= 15) {
		return false
	}
	fmt.Println("check")
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
	UserID   int64  `json:"userID"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
}
