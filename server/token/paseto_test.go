package token

import (
	"testing"
	"time"

	"github.com/brkss/vanillefraise2/utils"
	"github.com/stretchr/testify/require"
)

func TestPaseto(t *testing.T) {

	username := utils.RandomName()
	duration := time.Minute

	expireAt := time.Now().Add(duration)
	issuedAt := time.Now()

	maker, err := NewPasetoMaker(utils.RandomString(32))
	require.NoError(t, err)

	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, payload.Username, username)
	require.WithinDuration(t, payload.ExpiredAt, expireAt, time.Second)
	require.WithinDuration(t, payload.IssuedAt, issuedAt, time.Second)
}

func TestExpiredToken(t *testing.T) {
	maker, err := NewPasetoMaker(utils.RandomString(32))
	require.NoError(t, err)
	require.NotEmpty(t, maker)

	token, err := maker.CreateToken(utils.RandomName(), -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.Empty(t, payload)
	require.EqualError(t, err, ErrExpiredToken.Error())

}
