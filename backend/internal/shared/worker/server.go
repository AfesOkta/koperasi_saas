package worker

import (
	"context"

	"github.com/hibiken/asynq"
	"github.com/koperasi-gresik/backend/config"
	"github.com/sirupsen/logrus"
)

// Server encapsulates the background task consumer logic.
type Server struct {
	server *asynq.Server
	mux    *asynq.ServeMux
	mailer *Mailer
}

func NewServer(redisCfg config.RedisConfig, smtpCfg config.SMTPConfig) *Server {
	redisOpt := asynq.RedisClientOpt{
		Addr:     redisCfg.Host + ":" + redisCfg.Port,
		Password: redisCfg.Password,
		DB:       redisCfg.DB, // Might use a different DB for workers in prod
	}

	srv := asynq.NewServer(
		redisOpt,
		asynq.Config{
			Concurrency: 10, // Max concurrent workers
			Queues: map[string]int{
				"critical": 6, // 60% of workers process critical queue
				"default":  3,
				"low":      1,
			},
			ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
				logrus.Errorf("[Asynq] Error processing task %s: %v", task.Type(), err)
			}),
		},
	)

	mux := asynq.NewServeMux()
	mailer := NewMailer(smtpCfg)

	// Register Task Handlers
	mux.HandleFunc(TypeSendEmail, HandleSendEmailTask(mailer))

	return &Server{
		server: srv,
		mux:    mux,
		mailer: mailer,
	}
}

func (s *Server) RegisterClosingHandlers(svc interface {
	ProcessAllTenantsEOD(ctx context.Context, date string) error
	ProcessAllTenantsEOM(ctx context.Context, month, year int) error
}) {
	s.mux.HandleFunc(TypeRunEOD, HandleRunEODTask(svc))
	s.mux.HandleFunc(TypeRunEOM, HandleRunEOMTask(svc))
}

func (s *Server) Start() error {
	logrus.Info("[Asynq] Starting Email Worker Server...")
	return s.server.Start(s.mux)
}

func (s *Server) Stop() {
	logrus.Info("[Asynq] Shutting down Email Worker Server...")
	s.server.Shutdown()
}
