package test

// import (
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
// 	"godok/config"
// 	"godok/domain/user/dto"
// 	userdto "godok/domain/user/dto"
// 	"godok/domain/user/router"
// 	userService "godok/domain/user/service"
// 	"strings"

// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/labstack/echo/v4"
// 	"github.com/stretchr/testify/assert"
// )

// func TestInit(t *testing.T) {
// 	config.Init()
// 	config.EnvFile = "../env/test.env"
// 	userService.Init()
// }

// func TestGetUserInfo(t *testing.T) {
// 	var testCases = []struct {
// 		name        string
// 		input       string
// 		expectBody  dto.UserInfo
// 		expectCode  int
// 		expectError *string
// 	}{
// 		{
// 			name:  "유저 id가 1인 경우",
// 			input: "1",
// 			expectBody: dto.UserInfo{
// 				UserId:   1,
// 				UserName: "신민수",
// 				Email:    "alstn5038@gmail.com",
// 				Password: "123123",
// 				Age:      21,
// 			},
// 			expectCode:  http.StatusOK,
// 			expectError: nil,
// 		},
// 		{
// 			name:        "유저가 없는 경우",
// 			input:       "1000000",
// 			expectBody:  dto.UserInfo{},
// 			expectCode:  http.StatusNotFound,
// 			expectError: nil,
// 		},
// 		{
// 			name:        "id에 문자열을 입력한 경우",
// 			input:       "문자열",
// 			expectBody:  dto.UserInfo{},
// 			expectCode:  http.StatusBadRequest,
// 			expectError: nil,
// 		},
// 	}
// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			e := echo.New()

// 			req := httptest.NewRequest(echo.GET, "/user/info/"+tc.input, nil)
// 			rec := httptest.NewRecorder()

// 			e.GET("/user/info/:id", router.GetUserInfo)
// 			e.ServeHTTP(rec, req)

// 			var user dto.UserInfo
// 			_ = json.Unmarshal(rec.Body.Bytes(), &user)

// 			assert.Equal(t, tc.expectBody, user)
// 			assert.Equal(t, tc.expectCode, rec.Code)
// 		})
// 	}
// }

// func TestSignUp(t *testing.T) {
// 	var testCases = []struct {
// 		title       string
// 		userInfo    dto.UserInfo
// 		expectCode  int
// 		expectError *string
// 	}{
// 		{
// 			title: "입력값이 제대로 기입",
// 			userInfo: dto.UserInfo{
// 				UserName: "철수",
// 				Email:    "chulsoo@gmail.com",
// 				Password: "chulsoo",
// 				Age:      20,
// 			},
// 			expectCode:  http.StatusOK,
// 			expectError: nil,
// 		},
// 		{
// 			title: "입력값이 불충분하게 기입",
// 			userInfo: dto.UserInfo{
// 				UserName: "철수",
// 				Email:    "",
// 				Password: "chulsoo",
// 				Age:      10,
// 			},
// 			expectCode:  http.StatusBadRequest,
// 			expectError: nil,
// 		},
// 	}
// 	for _, tc := range testCases {
// 		t.Run(tc.title, func(t *testing.T) {
// 			e := echo.New()
// 			sendObj := fmt.Sprintf(`%v`, tc.userInfo)
// 			jsonStr := []byte(sendObj)
// 			req := httptest.NewRequest(echo.POST, "/user/info/", bytes.NewBuffer(jsonStr))
// 			rec := httptest.NewRecorder()

// 			e.POST("/user/info/", router.UserSignUp)
// 			e.ServeHTTP(rec, req)

// 			assert.Equal(t, tc.expectCode, rec.Code)
// 		})
// 	}
// }

// func TestSignIn(t *testing.T) {
// 	var testCases = []struct {
// 		title    string
// 		email    string
// 		password string
// 		result   int
// 	}{
// 		{
// 			title:    "login success",
// 			email:    "aa@gmail.com",
// 			password: "123",
// 			result:   http.StatusOK,
// 		},
// 		{
// 			title:    "login fail",
// 			email:    "invalid@gmail.com",
// 			password: "123",
// 			result:   http.StatusUnauthorized,
// 		},
// 	}
// 	for _, tc := range testCases {
// 		t.Run(tc.title, func(t *testing.T) {
// 			e := echo.New()
// 			sendObj := userdto.UserVerified{
// 				Email:    tc.email,
// 				Password: tc.password,
// 			}
// 			json, err := json.Marshal(sendObj)
// 			if err != nil {
// 				fmt.Println("marshaling failed")
// 			}
// 			req := httptest.NewRequest(echo.POST, "/user/signIn", strings.NewReader(string(json)))
// 			rec := httptest.NewRecorder()

// 			e.POST("/user/signIn", router.UserSignIn)
// 			e.ServeHTTP(rec, req)

// 			assert.Equal(t, tc.result, rec.Code)
// 		})
// 	}

// }
