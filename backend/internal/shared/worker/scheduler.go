package worker

import (
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/koperasi-gresik/backend/config"
	"github.com/sirupsen/logrus"
)

type Scheduler struct {
	scheduler *asynq.Scheduler
}

func NewScheduler(redisCfg config.RedisConfig) *Scheduler {
	redisOpt := asynq.RedisClientOpt{
		Addr:     redisCfg.Host + ":" + redisCfg.Port,
		Password: redisCfg.Password,
		DB:       redisCfg.DB,
	}

	return &Scheduler{
		scheduler: asynq.NewScheduler(redisOpt, &asynq.SchedulerOpts{
			Location: nil, // Use system timezone
		}),
	}
}

func (s *Scheduler) RegisterTask(cronSpec, taskType string, payload interface{}) error {
	b, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload for %s: %v", taskType, err)
	}

	task := asynq.NewTask(taskType, b)
	entryID, err := s.scheduler.Register(cronSpec, task)
	if err != nil {
		return fmt.Errorf("failed to register task %s: %v", taskType, err)
	}

	logrus.Infof("[Asynq] Registered periodic task: %s (cron: %s, entryID: %s)", taskType, cronSpec, entryID)
	return nil
}

func (s *Scheduler) Start() error {
	logrus.Info("[Asynq] Starting Scheduler...")
	return s.scheduler.Start()
}

func (s *Scheduler) Stop() {
	logrus.Info("[Asynq] Stopping Scheduler...")
	s.scheduler.Shutdown()
}
