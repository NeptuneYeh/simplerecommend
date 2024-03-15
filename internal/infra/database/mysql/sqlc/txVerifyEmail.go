package db

import (
	"context"
	"database/sql"
	"time"
)

type VerifyEmailTxParams struct {
	Code string `json:"code"`
}

type VerifyEmailTxResult struct {
	UpdateResult int64 `json:"updateResult"`
}

func (store *SQLStore) VerifyEmailTx(ctx context.Context, arg VerifyEmailTxParams) (VerifyEmailTxResult, error) {
	var result VerifyEmailTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		// STEP-1 找 account_valid_temp 內有沒有 該 code
		verifyData, err := q.GetAccountValid(ctx, arg.Code)
		if err != nil {
			return err
		}

		arg1 := UpdateAccountValidParams{
			ID:      verifyData.AccountID,
			IsValid: true,
			VerifiedAt: sql.NullTime{
				Time:  time.Now(),
				Valid: true,
			},
		}

		// STEP-2 更新 code 對應到的 account 的 is_valid & verified_at
		updateResult, err := q.UpdateAccountValid(ctx, arg1)
		result.UpdateResult, _ = updateResult.RowsAffected()
		if err != nil {
			return err
		}

		// STEP-3 刪除 account_valid_temp 內和這個 account_id 相關的 temp 資料
		err = q.DeleteAccountValidByAccountID(ctx, verifyData.AccountID)
		if err != nil {
			return err
		}

		return err
	})

	return result, err
}
