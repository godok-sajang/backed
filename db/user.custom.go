package db

import (
	"context"
	"database/sql"
	"echo_sample/domain/user/dto"
	"errors"
)

var getUserInfoQuery = `
SELECT * FROM test_table where user_id = $1
`

func (q *Queries) GetUserInfo(ctx context.Context, userId int) (user dto.UserInfo, err error) {
	err = q.db.QueryRowContext(ctx, getUserInfoQuery, userId).Scan(
		&user.UserId,
		&user.Nickname,
		&user.Email,
		&user.Password,
		&user.Birth)
	if err != nil {
		return
	}
	return
}

var createUserInfoQuery = `
INSERT INTO test_table(nickname, email, password, birth, gender) VALUES($1, $2, $3, $4, $5)
RETURNING user_id, nickname, email
`

func (q *Queries) CreateUserInfo(ctx context.Context, req dto.UserInfo) (dto.ExistUser, error) {
	var (
		err      error
		userInfo dto.ExistUser
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
		return dto.ExistUser{}, err
	}
	return userInfo, nil
}

var checkAuthQuery = `
	select user_id where email=$1, password=$2 limit1
`

func (q *Queries) CheckAuth(ctx context.Context, req dto.UserVerified) (userId int, err error) {
	if err = q.db.QueryRowContext(ctx, checkAuthQuery, req.Email, req.Password).Scan(&userId); err != nil {
		return 0, errors.New(err.Error())
	}
	return
}

var checkNicknameQuery = `
	select nickname from test_table where nickname=lower($1) limit 1
`

func (q *Queries) CheckNickname(ctx context.Context, nickname string) (sql.NullString, error) {
	var (
		ret sql.NullString
		err error
	)
	err = q.db.QueryRowContext(ctx, checkNicknameQuery, nickname).Scan(&ret)
	return ret, err
}

var checkEmailQuery = `
	select email from test_table where email=lower($1) limit 1
`

func (q *Queries) CheckEmail(ctx context.Context, email string) (sql.NullString, error) {
	var (
		ret sql.NullString
		err error
	)
	err = q.db.QueryRowContext(ctx, checkEmailQuery, email).Scan(&ret)
	return ret, err
}