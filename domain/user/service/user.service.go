package service

import (
	"context"
	"database/sql"
	"echo_sample/domain/user/dto"
	userdto "echo_sample/domain/user/dto"
	"echo_sample/middleware"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/pkg/errors"
)

type UserService struct{}

func (receiver *UserService) GetUserInfo(ctx context.Context, userId int) (user dto.UserInfo, err error) {
	user, err = dao.GetUserInfo(ctx, userId)
	if err != nil {
		return dto.UserInfo{}, err
	}
	return user, nil
}

func (receiver *UserService) CreateUserInfo(c context.Context, req dto.UserInfo) (ret dto.UserTokens, err error) {

	existNickname, err := dao.CheckNickname(c, req.Nickname)
	if err != nil && err != sql.ErrNoRows {
		err = errors.New(err.Error())
		return
	}
	if err != sql.ErrNoRows || existNickname.Valid {
		err = errors.New("nickname duplicated")
		return
	}
	existEmail, err := dao.CheckEmail(c, req.Nickname)
	if err != nil && err != sql.ErrNoRows {
		err = errors.New(err.Error())
		return
	}
	if err != sql.ErrNoRows || existEmail.Valid {
		err = errors.New("email duplicated")
		return
	}

	userInfo, err := dao.CreateUserInfo(c, req)
	if err != nil {

		return
	}
	accessToken, err := middleware.CreateToken(int64(userInfo.UserID), middleware.TokenValidationMinutes)
	refreshToken, err := middleware.CreateToken(int64(userInfo.UserID), middleware.RefreshValidationMinutes)

	userID := int64(userInfo.UserID)
	nickname := userInfo.Nickname
	email := userInfo.Email

	ret = dto.UserTokens{
		Success:      true,
		UserID:       sql.NullInt64{Int64: userID, Valid: userID != 0},
		AccessToken:  sql.NullString{String: accessToken, Valid: accessToken != ""},
		RefreshToken: sql.NullString{String: refreshToken, Valid: refreshToken != ""},
		ExistUserInfo: userdto.ExistUser{
			UserID:   userInfo.UserID,
			Nickname: nickname,
			Email:    email,
		},
	}
	return
}

func (receiver *UserService) SignIn(ctx context.Context, req dto.UserVerified) (tokenString string, err error) {
	userId, err := dao.CheckAuth(ctx, req)
	if err != nil {
		errors.Wrap(err, "SignIn failed")
	}

	// Set custom claims
	claims := &userdto.JwtCustomClaims{
		UserId: userId,
		Admin:  false,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token
	tokenString, err = token.SignedString([]byte("secret"))
	if err != nil {
		err = errors.Wrap(err, "SignedString")
		return
	}

	return tokenString, nil
}
