package repository

import (
	"context"
	"database/sql"
	"errors"
	"godok/clean/domain"
	"godok/util"

	eu "godok/util/errorutil"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(dialect, dsn string, idleConn, maxConn int) (domain.UserRepo, error) {
	db, err := sql.Open(dialect, dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(idleConn)
	db.SetMaxOpenConns(maxConn)

	return &userRepository{db}, nil
}

var getUserInfoQuery = `
SELECT * FROM test_table where user_id = $1
`

func (q *userRepository) GetUserInfo(ctx context.Context, userId int) (domain.UserInfo, error) {
	var user domain.UserInfo

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

func (q *userRepository) CreateUserInfo(ctx context.Context, req domain.UserInfo) (domain.UserInfo, error) {
	var (
		err      error
		userInfo domain.UserInfo
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
		return domain.UserInfo{}, err
	}
	return userInfo, nil
}

var checkAuthQuery = `
	select user_id, nickname, email from user_account where email=$1, password=$2 limit1
`

func (q *userRepository) CheckAuth(ctx context.Context, req domain.UserSignInRequest) (domain.CheckAuthResponse, error) {
	var ret domain.CheckAuthResponse
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

var getUserInfoByRequestQuery = `
	select * from user_account where is_delete=false and %v %v %v limit 1 
`

func (q *userRepository) GetUserInfoByRequest(c context.Context, req domain.UserInfoRequest) (domain.UserInfo, error) {
	var ret domain.UserInfo
	rows, err := q.db.QueryContext(c, getUserInfoByRequestQuery,
		req.GetQueryByNickname(),
		req.GetQueryByEmail(),
		req.GetQueryByPassword(),
	)
	if err != nil {
		return domain.UserInfo{}, eu.InternalError(err)
	}
	defer rows.Close()

	for rows.Next() {
		var row domain.UserInfo
		if err := rows.Scan(
			&row.UserID,
			&row.Nickname,
			&row.Email,
			&row.Password,
			&row.Birth,
			&row.Gender,
		); err != nil {
			return domain.UserInfo{}, eu.InternalError(err)
		}
	}
	if err := rows.Err(); err != nil {
		return domain.UserInfo{}, eu.InternalError(err)
	}
	return ret, nil

}
