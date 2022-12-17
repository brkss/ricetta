package db

import (
	"context"
	"testing"

	"github.com/brkss/vanillefraise2/utils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)



// RandomUser create random user and test db CreateUser operation 
func RandomUser(t *testing.T) User {
	arg := CreateUserParams{
		ID: uuid.New().String(),	
		Name: utils.RandomName(),
		Email: utils.RandomEmail(),
		Username: utils.RandomEmail(),
		Password: utils.RandomString(10),
	}
	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, user.ID, arg.ID)
	require.Equal(t, user.Name, arg.Name)
	require.Equal(t, user.Email, arg.Email)
	require.Equal(t, user.Password, arg.Password)

	return user
}

func TestCreateUser(t *testing.T){
	RandomUser(t);	
}
