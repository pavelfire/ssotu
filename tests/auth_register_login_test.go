package tests

import (
	"sso/tests/suite"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/golang-jwt/jwt"
	ssov1 "github.com/pavelfire/protostu/gen/go/sso"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	emptyAppId     = 0
	appId          = 1
	appSecret      = "test-secret-key"
	passDefaultLen = 10
)

func TestAuthRegisterLogin(t *testing.T) {
	ctx, st := suite.New(t)
	email := gofakeit.Email()
	password := randomFakePassword()

	respReg, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: password,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, respReg.GetUserId())

	respLogin, err := st.AuthClient.Login(ctx, &ssov1.LoginRequest{
		Email:    email,
		Password: password,
		AppId:    appId,
	})
	require.NoError(t, err)

	token := respLogin.GetToken()
	assert.NotEmpty(t, token)

	tokenParsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(appSecret), nil
	})
	require.NoError(t, err)
	// assert.NotEmpty(t, tokenParsed)

	claims, ok := tokenParsed.Claims.(jwt.MapClaims)
	assert.True(t, ok)

	assert.Equal(t, email, claims["email"].(string))
	assert.

}

func randomFakePassword() string {
	return gofakeit.Password(true, true, true, true, true, passDefaultLen)
}
