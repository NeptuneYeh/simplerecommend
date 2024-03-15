package db

import (
	"context"
	"database/sql"
	"fmt"
	db "github.com/NeptuneYeh/simplerecommend/internal/infra/database/mysql/sqlc"
	"github.com/stretchr/testify/require"
	"math/rand"
	"testing"
	"time"
)

func createTestAccount(t *testing.T) db.Account {
	randomNumber := rand.Intn(10000)
	name := "test_" + fmt.Sprintf("%04d", randomNumber)
	email := name + "@yopmail.com"

	arg := db.CreateAccountParams{
		Email:    email,
		Password: "12345678",
		Name:     name,
		IsValid:  true,
	}

	result, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	lastID, err := result.LastInsertId()
	require.NoError(t, err)
	getAccount, err := testQueries.GetAccount(context.Background(), lastID)

	require.Equal(t, arg.Name, getAccount.Name)
	require.Equal(t, arg.Password, getAccount.Password)
	require.Equal(t, arg.Email, getAccount.Email)

	require.NotZero(t, getAccount.ID)
	require.NotZero(t, getAccount.CreatedAt)

	return getAccount
}

func TestCreateAccount(t *testing.T) {
	createTestAccount(t)
}

func TestGetAccount(t *testing.T) {
	account1 := createTestAccount(t)
	account1Get, err := testQueries.GetAccount(context.Background(), account1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, account1Get)

	require.Equal(t, account1.ID, account1Get.ID)
	require.Equal(t, account1.Name, account1Get.Name)
	require.Equal(t, account1.Password, account1Get.Password)
	require.Equal(t, account1.Email, account1Get.Email)
	require.WithinDuration(t, account1.CreatedAt, account1Get.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	account1 := createTestAccount(t)
	time.Sleep(1 * time.Second)
	arg := db.UpdateAccountParams{
		ID:        account1.ID,
		Name:      "TEST_UPDATED",
		Password:  account1.Password,
		IsValid:   account1.IsValid,
		UpdatedAt: time.Now(),
	}

	_, err := testQueries.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)
	getAccount, err := testQueries.GetAccount(context.Background(), account1.ID)

	require.Equal(t, account1.ID, getAccount.ID)
	require.Equal(t, arg.Name, getAccount.Name)
	require.Equal(t, account1.Password, getAccount.Password)
	require.Equal(t, account1.Email, getAccount.Email)
	require.NotEqual(t, account1.UpdatedAt, getAccount.UpdatedAt)
}

func TestListAccount(t *testing.T) {
	for i := 0; i < 10; i++ {
		createTestAccount(t)
	}

	arg := db.ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, accounts, 5)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}

func TestDeleteAccount(t *testing.T) {
	account1 := createTestAccount(t)

	err := testQueries.DeleteAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	_, err = testQueries.GetAccount(context.Background(), account1.ID)
	require.EqualError(t, err, sql.ErrNoRows.Error())
}
