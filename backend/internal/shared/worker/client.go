package worker

import (
	"context"
	"encoding/json"

	"github.com/hibiken/asynq"
	"github.com/koperasi-gresik/backend/config"
	"github.com/sirupsen/logrus"
)

type TaskDistributor interface {
	DistributeTaskSendEmail(ctx context.Context, payload *PayloadSendEmail, opts ...asynq.Option) error
}

type RedisTaskDistributor struct {
	client *asynq.Client
}

func NewRedisTaskDistributor(redisCfg config.RedisConfig) TaskDistributor {
	redisOpt := asynq.RedisClientOpt{
		Addr:     redisCfg.Host + ":" + redisCfg.Port,
		Password: redisCfg.Password,
		DB:       redisCfg.DB,
	}
	return &RedisTaskDistributor{
		client: asynq.NewClient(redisOpt),
	}
}

func (d *RedisTaskDistributor) DistributeTaskSendEmail(ctx context.Context, payload *PayloadSendEmail, opts ...asynq.Option) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	task := asynq.NewTask(TypeSendEmail, jsonPayload, opts...)
	// Distribute queue based on opts or default
	info, err := d.client.EnqueueContext(ctx, task)
	if err != nil {
		return err
	}
	logrus.Infof("[Asynq] Enqueued task: type=%s, queue=%s, id=%s", task.Type(), info.Queue, info.ID)
	return nil
}
