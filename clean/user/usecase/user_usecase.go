package usecase

import (
	"context"
	"database/sql"
	"godok/clean/domain"
	"godok/domain/user/dto"
	userdto "godok/domain/user/dto"
	"godok/middleware"
	eu "godok/util/errorutil"
	"time"
)

type userUsecase struct {
	userRepo       domain.UserRepo
	contextTimeout time.Duration
}

func New(c context.Context, userRepo domain.UserRepo, contextTimeout time.Duration) domain.UserUsecase {
	return &userUsecase{
		userRepo:       userRepo,
		contextTimeout: contextTimeout,
	}
}

func (us *userUsecase) CreateUserInfo(c context.Context, req domain.UserInfoRequest) (dto.UserTokens, error) {
	ok, err := CheckNicknameDuplicated(c, us, req.Nickname)
	if err != nil || !ok {
		return dto.UserTokens{}, err
	}

	ok, err = CheckEmailDuplicated(c, us, req.Email)
	if err != nil || !ok {
		return dto.UserTokens{}, err
	}

	userInfo, err := us.userRepo.CreateUserInfo(c, req.ToDtoUserInfo())
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

func (us *userUsecase) SignIn(ctx context.Context, req domain.UserSignInRequest) (dto.UserTokens, error) {
	userInfo, err := us.userRepo.GetUserInfoByRequest(ctx, domain.UserInfoRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return dto.UserTokens{}, err
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

func CheckNicknameDuplicated(c context.Context, us *userUsecase, nickname string) (bool, error) {
	users, err := us.userRepo.GetUserInfoByRequest(c, domain.UserInfoRequest{Nickname: nickname})
	if err != nil {
		return true, err
	}
	if users.UserID != 0 {
		return true, eu.New().WithCustomCode("nickname duplicated")
	}
	return false, nil
}

func CheckEmailDuplicated(c context.Context, us *userUsecase, email string) (bool, error) {
	users, err := us.userRepo.GetUserInfoByRequest(c, domain.UserInfoRequest{Email: email})
	if err != nil {
		return true, err
	}
	if users.UserID != 0 {
		return true, eu.New().WithCustomCode("email duplicated")
	}
	return false, nil
}
