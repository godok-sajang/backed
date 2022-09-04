package dto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetQueryByNickname(t *testing.T) {
	testCase := []struct {
		req    UserInfoRequest
		result string
	}{
		{
			req: UserInfoRequest{
				Nickname: "",
			},
			result: "",
		},
		{
			req: UserInfoRequest{
				Nickname: "something",
			},
			result: "and nickname=something",
		},
	}
	for _, test := range testCase {
		ret := test.req.GetQueryByNickname()
		assert.Equal(t, test.result, ret)
	}
}

func TestGetQueryByEmail(t *testing.T) {
	testCase := []struct {
		req    UserInfoRequest
		result string
	}{
		{
			req: UserInfoRequest{
				Email: "",
			},
			result: "",
		},
		{
			req: UserInfoRequest{
				Email: "alstn5038@gmail.com",
			},
			result: "and email='alstn5038@gmail.com'",
		},
	}
	for _, test := range testCase {
		ret := test.req.GetQueryByEmail()
		assert.Equal(t, test.result, ret)
	}
}

func TestGetQueryByPassword(t *testing.T) {
	testCase := []struct {
		req    UserInfoRequest
		result string
	}{
		{
			req: UserInfoRequest{
				Password: "",
			},
			result: "",
		},
		{
			req: UserInfoRequest{
				Password: "123123",
			},
			result: "and password='123123'",
		},
	}
	for _, test := range testCase {
		ret := test.req.GetQueryByPassword()
		assert.Equal(t, test.result, ret)
	}
}
