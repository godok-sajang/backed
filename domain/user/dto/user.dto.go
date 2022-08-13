package dto

import (
	"database/sql"
	"time"

	"github.com/golang-jwt/jwt"
)

type UserInfoRequest struct {
	Nickname string  `json:"nickname"`
	Email    *string `json:"email"`
	Password *string `json:"password"`
	Birth    string  `json:"birth"`
	Gender   int     `json:"gender"`
}

type UserInfo struct {
	UserId   int
	Nickname string
	Email    *string
	Password *string
	Birth    time.Time
	Gender   int
}

type UserVerified struct {
	Email    string
	Password string
}

type JwtCustomClaims struct {
	UserId int  `json:"userId"`
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
	UserID   int     `json:"userID"`
	Nickname *string `json:"nickname"`
	Email    string  `json:"email"`
}
