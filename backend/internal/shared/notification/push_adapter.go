package notification

import (
	"context"
	"github.com/sirupsen/logrus"
)

type PushAdapter interface {
	Send(ctx context.Context, token, title, message string) error
}

type fcmAdapter struct {
	// apiKey or credentials would go here
}

func NewFCMAdapter() PushAdapter {
	return &fcmAdapter{}
}

func (a *fcmAdapter) Send(ctx context.Context, token, title, message string) error {
	// This is where Firebase Admin SDK would be called:
	// client.Send(ctx, &messaging.Message{ Token: token, Notification: ... })
	
	logrus.Infof("[FCM] Delivering Push to %s: %s - %s", token, title, message)
	return nil
}
