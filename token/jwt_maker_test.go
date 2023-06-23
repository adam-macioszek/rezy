package token

import (
	"fmt"
	"testing"
	"time"

	"github.com/adam-macioszek/rezy/random"
	"github.com/stretchr/testify/require"
)

func TestJWTMaker(t *testing.T) {
	maker, err := NewJWTMaker(random.RandomString(32))
	require.NoError(t, err)
	username := "root"

	duration := time.Minute
	iat := time.Now()
	exp := iat.Add(duration)
	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, iat, payload.IssuedAt, time.Second)
	require.WithinDuration(t, exp, payload.ExpiredAt, time.Second)
}

func TestExpiredToken(t *testing.T) {
	maker, err := NewJWTMaker(random.RandomString(32))
	require.NoError(t, err)

	username := "root"
	token, err := maker.CreateToken(username, -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	expectedError := fmt.Errorf("token has invalid claims: %v", ErrExpiredToken.Error())
	require.EqualError(t, err, expectedError.Error())
	require.Nil(t, payload)
}
