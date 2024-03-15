package accountRequests

import "errors"

type VerifyEmailRequest struct {
	Code string `form:"code" binding:"required"`
}

var (
	ErrVerifyEmail = errors.New("failed verify email")
)
