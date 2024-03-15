package worker

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/NeptuneYeh/simplerecommend/init/logger"
	"github.com/hibiken/asynq"
)

type PayloadSendVerifyEmail struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Code  string `json:"code"`
}

const TaskSendVerifyEmail = "task:send_verify_email"

func (distributor *RedisTaskDistributor) DistributeTaskSendVerifyEmail(ctx context.Context, payload *PayloadSendVerifyEmail, opts ...asynq.Option) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal task payload: %w", err)
	}

	task := asynq.NewTask(TaskSendVerifyEmail, jsonPayload, opts...)
	info, err := distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	logger.MyLog.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).Str("queue", info.Queue).Int("max_retry", info.MaxRetry).Msg("enqueued task")
	return nil
}

func (processor *RedisTaskProcessor) ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error {
	var payload PayloadSendVerifyEmail
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", asynq.SkipRetry)
	}

	account, err := processor.store.GetAccountByEmail(ctx, payload.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("account doesn't exist: %w", asynq.SkipRetry)
		}
		return fmt.Errorf("failed to get account: %w", err)
	}
	subject := "Welcome to Simple Recommend"
	verifyUrl := fmt.Sprintf("http://localhost:8080/account/verify/email?code=%s", payload.Code)
	content := fmt.Sprintf(`Hello %s,<br/>
	Thank you for registering with us!<br/>
	Your code is: %s<br/>
	Please <a href="%s">click here</a> to verify your email address.<br/>
	`, payload.Name, payload.Code, verifyUrl)
	to := []string{payload.Email}
	err = processor.mailer.SendEmail(subject, content, to, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to send verify email: %w", err)
	}
	// send event to queue (prepare to send email to account)
	logger.MyLog.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).Str("email", account.Email).Msg("processed task")
	return nil
}
