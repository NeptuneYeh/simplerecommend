package token

import (
	"fmt"
	"github.com/NeptuneYeh/simplerecommend/pkg/helper"
	myToken "github.com/NeptuneYeh/simplerecommend/pkg/token"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
	"math/rand"
	"testing"
	"time"
)

func TestJWTMaker(t *testing.T) {
	randomString, err := helper.RandomString(32)
	require.NoError(t, err)
	maker, err := myToken.NewJWTMaker(randomString)
	require.NoError(t, err)

	randomNumber := rand.Intn(10000)
	username := "test_" + fmt.Sprintf("%04d", randomNumber)

	duration := time.Minute
	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

func TestExpiredJWTToken(t *testing.T) {
	randomString, err := helper.RandomString(32)
	require.NoError(t, err)
	maker, err := myToken.NewJWTMaker(randomString)
	require.NoError(t, err)

	randomNumber := rand.Intn(10000)
	username := "test_" + fmt.Sprintf("%04d", randomNumber)

	duration := -time.Minute
	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, myToken.ErrExpiredToken.Error())
	require.Nil(t, payload)
}

func TestInvalidJWTTokenAlgNone(t *testing.T) {
	randomNumber := rand.Intn(10000)
	username := "test_" + fmt.Sprintf("%04d", randomNumber)
	payload, err := myToken.NewPayload(username, time.Minute)
	require.NoError(t, err)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	badToken, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	randomString, err := helper.RandomString(32)
	require.NoError(t, err)
	maker, err := myToken.NewJWTMaker(randomString)
	require.NoError(t, err)

	payload, err = maker.VerifyToken(badToken)
	require.Error(t, err)
	require.EqualError(t, err, myToken.ErrInvalidToken.Error())
	require.Nil(t, payload)
}
