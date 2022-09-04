package service

import (
	"context"
	"database/sql"
	"godok/domain/user/dto"
	userdto "godok/domain/user/dto"
	"godok/middleware"
	eu "godok/util/errorutil"

	"github.com/pkg/errors"
)

type UserService struct{}

func (receiver *UserService) GetUserInfo(ctx context.Context, userId int) (dto.UserInfo, error) {
	user, err := dao.GetUserInfo(ctx, userId)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (receiver *UserService) CreateUserInfo(c context.Context, req dto.UserInfoRequest) (dto.UserTokens, error) {
	ok, err := CheckNicknameDuplicated(c, req.Nickname)
	if err != nil || !ok {
		return dto.UserTokens{}, err
	}

	ok, err = CheckEmailDuplicated(c, req.Email)
	if err != nil || !ok {
		return dto.UserTokens{}, err
	}

	userInfo, err := dao.CreateUserInfo(c, req.ToDtoUserInfo())
	if err != nil {
		return dto.UserTokens{}, err
	}

	accessToken, err := middleware.CreateToken(int64(userInfo.UserID), middleware.TokenValidationMinutes)
	refreshToken, err := middleware.CreateToken(int64(userInfo.UserID), middleware.RefreshValidationMinutes)

	userID := int64(userInfo.UserID)
	nickname := userInfo.Nickname
	email := userInfo.Email

	ret := dto.UserTokens{
		Success:      true,
		UserID:       sql.NullInt64{Int64: userID, Valid: userID != 0},
		AccessToken:  sql.NullString{String: accessToken, Valid: accessToken != ""},
		RefreshToken: sql.NullString{String: refreshToken, Valid: refreshToken != ""},
		ExistUserInfo: userdto.ExistUser{
			UserID:   userID,
			Nickname: nickname,
			Email:    email,
		},
	}
	return ret, nil
}

func (receiver *UserService) SignIn(ctx context.Context, req dto.UserSignInRequest) (dto.UserTokens, error) {
	userInfo, err := dao.CheckAuth(ctx, req)
	if err != nil {
		errors.Wrap(err, "CheckAuth failed")
	}

	if userInfo.UserID == 0 {
		return dto.UserTokens{}, eu.New().WithCustomCode("UnAuthorized")
	}

	accessToken, err := middleware.CreateToken(int64(userInfo.UserID), middleware.TokenValidationMinutes)
	refreshToken, err := middleware.CreateToken(int64(userInfo.UserID), middleware.RefreshValidationMinutes)

	ret := dto.UserTokens{
		Success:      true,
		UserID:       sql.NullInt64{Int64: int64(userInfo.UserID), Valid: true},
		AccessToken:  sql.NullString{String: accessToken, Valid: accessToken != ""},
		RefreshToken: sql.NullString{String: refreshToken, Valid: refreshToken != ""},
		ExistUserInfo: userdto.ExistUser{
			UserID:   int64(userInfo.UserID),
			Nickname: userInfo.Nickname,
			Email:    userInfo.Email,
		},
	}

	return ret, nil
}

func CheckNicknameDuplicated(c context.Context, nickname string) (bool, error) {
	users, err := dao.GetUserInfoByRequest(c, dto.UserInfoRequest{Nickname: nickname})
	if err != nil {
		return true, err
	}
	if len(users) != 0 {
		return true, eu.New().WithCustomCode("nickname duplicated")
	}
	return false, nil
}

func CheckEmailDuplicated(c context.Context, email string) (bool, error) {
	users, err := dao.GetUserInfoByRequest(c, dto.UserInfoRequest{Email: email})
	if err != nil {
		return true, err
	}
	if len(users) != 0 {
		return true, eu.New().WithCustomCode("email duplicated")
	}
	return false, nil
}
