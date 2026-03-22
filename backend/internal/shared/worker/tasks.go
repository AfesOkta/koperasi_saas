package worker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/sirupsen/logrus"
)

const (
	TypeSendEmail = "email:send"
	TypeRunEOD    = "closing:eod"
	TypeRunEOM    = "closing:eom"
)

type PayloadSendEmail struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"` // HTML string
}

type PayloadRunEOD struct {
	Date string `json:"date"` // YYYY-MM-DD
}

type PayloadRunEOM struct {
	Month int `json:"month"`
	Year  int `json:"year"`
}

// ClosingService interface placeholder to avoid circular imports if needed, 
// but we'll use a functional approach or interface injection in server.go


func HandleSendEmailTask(mailer *Mailer) asynq.HandlerFunc {
	return func(ctx context.Context, t *asynq.Task) error {
		var p PayloadSendEmail
		if err := json.Unmarshal(t.Payload(), &p); err != nil {
			return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
		}

		err := mailer.SendEmail(p.To, p.Subject, p.Body)
		if err != nil {
			logrus.Errorf("[Asynq] Failed to send email to %s: %v", p.To, err)
			return err // Return err allows Asynq to auto-retry via exponential backoff
		}

		logrus.Infof("[Asynq] Successfully sent email to %s", p.To)
		return nil
	}
}

func HandleRunEODTask(svc interface{ ProcessAllTenantsEOD(ctx context.Context, date string) error }) asynq.HandlerFunc {
	return func(ctx context.Context, t *asynq.Task) error {
		var p PayloadRunEOD
		if err := json.Unmarshal(t.Payload(), &p); err != nil {
			return err
		}
		return svc.ProcessAllTenantsEOD(ctx, p.Date)
	}
}

func HandleRunEOMTask(svc interface{ ProcessAllTenantsEOM(ctx context.Context, month, year int) error }) asynq.HandlerFunc {
	return func(ctx context.Context, t *asynq.Task) error {
		var p PayloadRunEOM
		if err := json.Unmarshal(t.Payload(), &p); err != nil {
			return err
		}
		return svc.ProcessAllTenantsEOM(ctx, p.Month, p.Year)
	}
}
