package db

import (
	"context"
	"errors"
	"godok/domain/user/dto"
	"godok/util"

	eu "godok/util/errorutil"
)

var getUserInfoQuery = `
SELECT * FROM test_table where user_id = $1
`

func (q *Queries) GetUserInfo(ctx context.Context, userId int) (dto.UserInfo, error) {
	var user dto.UserInfo

	if err := q.db.QueryRowContext(ctx, getUserInfoQuery, userId).Scan(
		&user.UserID,
		&user.Nickname,
		&user.Email,
		&user.Password,
		&user.Birth); err != nil {
		return user, eu.InternalError(err)
	}
	return user, nil
}

var createUserInfoQuery = `
INSERT INTO test_table(nickname, email, password, birth, gender) VALUES($1, $2, $3, $4, $5)
RETURNING user_id, nickname, email
`

func (q *Queries) CreateUserInfo(ctx context.Context, req dto.UserInfo) (dto.UserInfo, error) {
	var (
		err      error
		userInfo dto.UserInfo
	)
	if err = q.db.QueryRowContext(ctx, createUserInfoQuery,
		req.Nickname,
		req.Email,
		req.Password,
		req.Birth,
		req.Gender,
	).Scan(
		&userInfo.UserID,
		&userInfo.Nickname,
		&userInfo.Email,
	); err != nil {
		err = errors.New(err.Error())
		return dto.UserInfo{}, err
	}
	return userInfo, nil
}

var checkAuthQuery = `
	select user_id, nickname, email from user_account where email=$1, password=$2 limit1
`

type checkAuthResponse struct {
	UserID   int    `json:"userId" bson:"user_id"`
	Nickname string `json:"nickname" bson:"nickname"`
	Email    string `json:"email" bson:"email"`
}

func (q *Queries) CheckAuth(ctx context.Context, req dto.UserSignInRequest) (checkAuthResponse, error) {
	var ret checkAuthResponse
	if err := q.db.QueryRowContext(ctx, checkAuthQuery,
		req.Email,
		util.HashPassword(req.Password),
	).Scan(
		&ret.UserID,
		&ret.Nickname,
		&ret.Email,
	); err != nil {
		return ret, eu.InternalError(err)
	}
	return ret, nil
}

var checkNicknameQuery = `
	select nickname from user_account where nickname=lower($1) limit 1
`

func (q *Queries) CheckNickname(ctx context.Context, nickname string) ([]string, error) {
	var (
		ret []string
		err error
	)
	rows, err := q.db.QueryContext(ctx, checkNicknameQuery, nickname)
	if err != nil {
		return nil, eu.InternalError(err)
	}
	defer rows.Close()
	for rows.Next() {
		var row string
		if err := rows.Scan(
			&row,
		); err != nil {
			return nil, eu.InternalError(err)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, eu.InternalError(err)
	}
	return ret, nil
}

var getUserInfoByRequestQuery = `
	select * from user_account where true and %v %v %v limit 1 
`

func (q *Queries) GetUserInfoByRequest(c context.Context, req dto.UserInfoRequest) ([]dto.UserInfo, error) {
	var ret []dto.UserInfo
	rows, err := q.db.QueryContext(c, getUserInfoByRequestQuery,
		req.GetQueryByNickname(),
		req.GetQueryByEmail(),
		req.GetQueryByPassword(),
	)
	if err != nil {
		return nil, eu.InternalError(err)
	}
	defer rows.Close()

	for rows.Next() {
		var row dto.UserInfo
		if err := rows.Scan(
			&row.UserID,
			&row.Nickname,
			&row.Email,
			&row.Password,
			&row.Birth,
			&row.Gender,
		); err != nil {
			return nil, eu.InternalError(err)
		}
		ret = append(ret, row)
	}
	if err := rows.Err(); err != nil {
		return nil, eu.InternalError(err)
	}
	return ret, nil

}
