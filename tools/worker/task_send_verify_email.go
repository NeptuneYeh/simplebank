package worker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	postgresdb "github.com/NeptuneYeh/simplebank/internal/infrastructure/database/postgres/sqlc"
	"github.com/NeptuneYeh/simplebank/tools/helper"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

type PayloadSendVerifyEmail struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Code     string `json:"code"`
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

	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).Str("queue", info.Queue).Int("max_retry", info.MaxRetry).Msg("enqueued task")
	return nil
}

func (processor *RedisTaskProcessor) ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error {
	var payload PayloadSendVerifyEmail
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", asynq.SkipRetry)
	}

	account, err := processor.store.GetUser(ctx, payload.Username)
	if err != nil {
		if errors.Is(err, postgresdb.ErrRecordNotFound) {
			return fmt.Errorf("account doesn't exist: %w", asynq.SkipRetry)
		}
		return fmt.Errorf("failed to get account: %w", err)
	}

	randomString, _ := helper.RandomString(32)
	verifyEmail, err := processor.store.CreateVerifyEmail(ctx, postgresdb.CreateVerifyEmailParams{
		Username:   account.Username,
		Email:      account.Email,
		SecretCode: randomString,
	})
	if err != nil {
		return fmt.Errorf("failed to create verify email: %w", err)
	}

	subject := "Welcome to Simple Bank"
	verifyUrl := fmt.Sprintf("http://localhost:8080/v1/verify_email?email_id=%d&secret_code=%s",
		verifyEmail.ID, verifyEmail.SecretCode)
	content := fmt.Sprintf(`Hello %s,<br/>
	Thank you for registering with us!<br/>
	Your code is: %s<br/>
	Please <a href="%s">click here</a> to verify your email address.<br/>
	`, account.FullName, payload.Code, verifyUrl)
	to := []string{payload.Email}
	err = processor.mailer.SendEmail(subject, content, to, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to send verify email: %w", err)
	}
	// send event to queue (prepare to send email to account)
	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).Str("email", account.Email).Msg("processed task")
	return nil
}
