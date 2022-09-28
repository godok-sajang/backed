package domain

import "context"

type UserUsecase interface {
}

type UserRepo interface {
	CheckAuth(context.Context, UserSignInRequest) (CheckAuthResponse, error)
	GetUserInfoByRequest(context.Context, UserInfoRequest) (UserInfo, error)
	CreateUserInfo(context.Context, UserInfo) (UserInfo, error)
	GetUserInfo(context.Context, int) (UserInfo, error)
}
