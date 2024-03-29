// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package db

import (
	"context"
	"database/sql"
)

type Querier interface {
	CreateAccount(ctx context.Context, arg CreateAccountParams) (sql.Result, error)
	CreateAccountValid(ctx context.Context, arg CreateAccountValidParams) (sql.Result, error)
	CreateProduct(ctx context.Context, arg CreateProductParams) (sql.Result, error)
	CreateRecommendProduct(ctx context.Context, productID int64) (sql.Result, error)
	DeleteAccount(ctx context.Context, id int64) error
	DeleteAccountValidByAccountID(ctx context.Context, accountID int64) error
	DeleteProduct(ctx context.Context, id int64) error
	DeleteRecommendProduct(ctx context.Context, id int64) error
	GetAccount(ctx context.Context, id int64) (Account, error)
	GetAccountByEmail(ctx context.Context, email string) (Account, error)
	GetAccountValid(ctx context.Context, code string) (AccountValidTemp, error)
	GetProduct(ctx context.Context, id int64) (Product, error)
	ListAccounts(ctx context.Context, arg ListAccountsParams) ([]Account, error)
	ListProducts(ctx context.Context, arg ListProductsParams) ([]Product, error)
	ListRecommendProducts(ctx context.Context) ([]ListRecommendProductsRow, error)
	UpdateAccount(ctx context.Context, arg UpdateAccountParams) (sql.Result, error)
	UpdateAccountValid(ctx context.Context, arg UpdateAccountValidParams) (sql.Result, error)
	UpdateProduct(ctx context.Context, arg UpdateProductParams) (sql.Result, error)
}

var _ Querier = (*Queries)(nil)
