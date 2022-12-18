package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestHashingPassword(t *testing.T) {
	// testing hashing password
	password := RandomString(10)
	hash, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hash)

	// testing comparing hash with password !
	err = VerifyPassword(password, hash)
	require.NoError(t, err)
}

func TestInvalidPassword(t *testing.T) {

	password1 := RandomString(10)
	password2 := RandomString(10)
	hash, err := HashPassword(password1)
	require.NoError(t, err)
	require.NotEmpty(t, hash)

	// testing comparing false password
	err = VerifyPassword(password2, hash)
	require.Error(t, err)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())
}
