package asynq

import (
	"context"
	"github.com/NeptuneYeh/simplerecommend/init/config"
	db "github.com/NeptuneYeh/simplerecommend/internal/infra/database/mysql/sqlc"
	"github.com/NeptuneYeh/simplerecommend/pkg/mail"
	"github.com/NeptuneYeh/simplerecommend/worker"
	"github.com/hibiken/asynq"
	"log"
	"time"
)

var MyAsynq *Module

type Module struct {
	RedisOpt        *asynq.RedisClientOpt
	TaskDistributor worker.TaskDistributor
	TaskProcessor   worker.TaskProcessor
}

func NewModule() *Module {
	redisOpt := &asynq.RedisClientOpt{
		Addr: config.MyConfig.RedisAddress,
	}

	taskDistributor := worker.NewRedisTaskDistributor(*redisOpt)

	asynqModule := &Module{
		TaskDistributor: taskDistributor,
		RedisOpt:        redisOpt,
	}
	MyAsynq = asynqModule

	return asynqModule
}

func (module *Module) Run(redisOpt *asynq.RedisClientOpt, store db.Store) {
	mailer := mail.NewGmailSender(config.MyConfig.EmailSenderName, config.MyConfig.EmailSenderAddress, config.MyConfig.EmailSenderPassword)
	module.TaskProcessor = worker.NewRedisTaskProcessor(*redisOpt, store, mailer)
	go func() {
		err := module.TaskProcessor.Start()
		if err != nil {
			log.Fatalf("Failed to start task processor: %w", err)
		}
	}()
}

func (module *Module) Shutdown() error {
	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	module.TaskProcessor.Shutdown()
	return nil
}
