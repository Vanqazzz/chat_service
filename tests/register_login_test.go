package tests

import (
	"chat_service/tests/suite"
	"testing"
	"time"

	protos "github.com/Vanqazzz/protos/gen/go/chat_service/auth"
	"github.com/brianvoe/gofakeit"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	emptyAppID = 0
	appID      = 1
	appSecret  = "test-secret"

	PassDefaultLen = 10
)

func TestRegisterLogin(t *testing.T) {

	ctx, st := suite.New(t)

	email := gofakeit.Email()
	pass := gofakeit.Password(true, true, true, true, false, PassDefaultLen)

	respReg, err := st.AuthClient.Register(ctx, &protos.RegisterRequest{
		Email:    email,
		Password: pass,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, respReg.GetUserId())

	respLogin, err := st.AuthClient.Login(ctx, &protos.LoginRequest{
		Email:    email,
		Password: pass,
		AppId:    appID})
	require.NoError(t, err)

	token := respLogin.GetToken()
	require.NotEmpty(t, token)

	loginTime := time.Now()

	tokenParsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(appSecret), nil
	})
	require.NoError(t, err)

	claims, ok := tokenParsed.Claims.(jwt.MapClaims)
	require.True(t, ok)

	assert.Equal(t, respReg.GetUserId(), int64(claims["uid"].(float64)))
	assert.Equal(t, email, claims["email"].(string))
	assert.Equal(t, appID, int(claims["app_id"].(float64)))

	const deltaSeconds = 1

	assert.InDelta(t, loginTime.Add(st.Cfg.TokenTTL).Unix(), claims["exp"].(float64), deltaSeconds)
}

func TestRegister_Duplicate(t *testing.T) {

	ctx, st := suite.New(t)

	email := gofakeit.Email()
	pass := gofakeit.Password(true, true, true, true, false, PassDefaultLen)

	respReg, err := st.AuthClient.Register(ctx, &protos.RegisterRequest{Email: email, Password: pass})
	require.NoError(t, err)
	require.NotEmpty(t, respReg.GetUserId())

	respReg, err = st.AuthClient.Register(ctx, &protos.RegisterRequest{Email: email, Password: pass})
	require.Error(t, err)
	assert.Empty(t, respReg.GetUserId())
	assert.ErrorContains(t, err, "user already exists")

}

func TestLogin_Fail(t *testing.T) {
	ctx, st := suite.New(t)

	email := gofakeit.Email()
	pass := gofakeit.Password(true, true, true, true, false, PassDefaultLen)

	respReg, err := st.AuthClient.Register(ctx, &protos.RegisterRequest{Email: email, Password: pass})
	require.NoError(t, err)
	assert.NotEmpty(t, respReg.GetUserId())

	pass = gofakeit.Password(true, true, true, true, false, PassDefaultLen)

	logResp, err := st.AuthClient.Login(ctx, &protos.LoginRequest{Email: email, Password: pass, AppId: appID})
	require.Error(t, err)
	assert.Empty(t, logResp.GetToken())
	assert.ErrorContains(t, err, "incorrect email or password")
}

func TestRegister_FailCases(t *testing.T) {
	ctx, st := suite.New(t)

	tests := []struct {
		name        string
		email       string
		password    string
		appID       int32
		expectedErr string
	}{{
		name:        "Register with empty password",
		email:       gofakeit.Email(),
		password:    "",
		appID:       1,
		expectedErr: "password is required",
	}, {
		name:        "Register with empty email",
		email:       "",
		password:    gofakeit.Password(true, true, true, true, false, PassDefaultLen),
		appID:       1,
		expectedErr: "email is required",
	}, {
		name:        "Register with empty email and password",
		email:       "",
		password:    "",
		appID:       1,
		expectedErr: "email and password is empty",
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := st.AuthClient.Register(ctx, &protos.RegisterRequest{
				Email:    tt.email,
				Password: tt.password,
			})
			require.Error(t, err)
			require.Contains(t, err.Error(), tt.expectedErr)
		})
	}
}

func TestLogin_FailCases(t *testing.T) {
	ctx, st := suite.New(t)

	tests := []struct {
		name        string
		email       string
		password    string
		appID       int32
		exceptedErr string
	}{
		{
			name:        "Login with empty email",
			email:       "",
			password:    gofakeit.Password(true, true, true, true, false, PassDefaultLen),
			appID:       1,
			exceptedErr: "email is required",
		}, {
			name:        "Login with empty password",
			email:       gofakeit.Email(),
			password:    "",
			appID:       1,
			exceptedErr: "password is required",
		}, {
			name:        "Login with empty password and email",
			email:       "",
			password:    "",
			appID:       1,
			exceptedErr: "email and password is required",
		}, {
			name:        "Login with incorrect password",
			email:       gofakeit.Email(),
			password:    gofakeit.Password(true, true, true, true, false, PassDefaultLen),
			appID:       1,
			exceptedErr: "invalid email or password",
		}, {
			name:        "Login without appID",
			email:       gofakeit.Email(),
			password:    gofakeit.Password(true, true, true, true, false, PassDefaultLen),
			appID:       emptyAppID,
			exceptedErr: "app id is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			_, err := st.AuthClient.Login(ctx, &protos.LoginRequest{Email: tt.email, Password: tt.password, AppId: tt.appID})
			require.Error(t, err)
			require.Contains(t, err, tt.exceptedErr)
		})
	}

}
