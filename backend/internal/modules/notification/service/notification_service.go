package service

import (
	"context"
	"fmt"

	"github.com/koperasi-gresik/backend/internal/modules/iam/model"
	iamRepository "github.com/koperasi-gresik/backend/internal/modules/iam/repository"
	notificationModel "github.com/koperasi-gresik/backend/internal/modules/notification/model"
	"github.com/koperasi-gresik/backend/internal/modules/notification/repository"
	"github.com/koperasi-gresik/backend/internal/shared/event"
	"github.com/koperasi-gresik/backend/internal/shared/notification"
	"github.com/koperasi-gresik/backend/internal/shared/worker"
	"github.com/sirupsen/logrus"
)

type NotificationService interface {
	Start(ctx context.Context)
}

type notificationService struct {
	repo       repository.NotificationRepository
	tokenRepo  iamRepository.TokenRepository
	pushAdp    notification.PushAdapter
	subscriber event.Subscriber
	taskDist   worker.TaskDistributor
}

func NewNotificationService(
	repo repository.NotificationRepository,
	tokenRepo iamRepository.TokenRepository,
	pushAdp notification.PushAdapter,
	sub event.Subscriber,
	taskDist worker.TaskDistributor,
) NotificationService {
	return &notificationService{
		repo:       repo,
		tokenRepo:  tokenRepo,
		pushAdp:    pushAdp,
		subscriber: sub,
		taskDist:   taskDist,
	}
}

func (s *notificationService) Start(ctx context.Context) {
	go s.subscriber.Consume(ctx, s.handleEvent)
}

func (s *notificationService) handleEvent(ev event.Event) error {
	switch ev.Type {
	case event.EventUserCreated:
		s.handleUserCreated(ev)
	case event.EventLoanApproved:
		s.handleLoanApproved(ev)
	case event.EventStockLow:
		s.handleStockLow(ev)
	}
	return nil
}

func (s *notificationService) handleUserCreated(ev event.Event) {
	// Assume payload contains the User model or data
	user, ok := ev.Payload.(*model.User)
	if !ok {
		return
	}

	// 1. Save In-App Notification
	notif := &notificationModel.Notification{
		UserID:  user.ID,
		Title:   "Welcome!",
		Message: fmt.Sprintf("Hi %s, welcome to Koperasi SaaS. Your account is ready.", user.Name),
		Type:    "success",
	}
	notif.OrganizationID = ev.OrganizationID
	_ = s.repo.Create(context.Background(), notif)

	// 2. Enqueue Push Notification (FCM Placeholder)
	s.handlePushNotification(user.ID, notif.Title, notif.Message)

	// 3. Enqueue Email Task
	emailBody := fmt.Sprintf(`
		<h3>Welcome to Koperasi SaaS</h3>
		<p>Hi %s,</p>
		<p>Your account has been created successfully.</p>
		<p>Thank you!</p>
	`, user.Name)

	payload := &worker.PayloadSendEmail{
		To:      user.Email,
		Subject: "Welcome to Koperasi SaaS",
		Body:    emailBody,
	}
	err := s.taskDist.DistributeTaskSendEmail(context.Background(), payload)
	if err != nil {
		logrus.Errorf("[NotificationService] Failed to enqueue welcome email: %v", err)
	}
}

func (s *notificationService) handleLoanApproved(ev event.Event) {
	// Payload would be Loan ID or struct. For demo, we assume we get a map with info.
	payload, ok := ev.Payload.(map[string]interface{})
	if !ok {
		return
	}

	loanNum := payload["loan_number"].(string)
	memberEmail := payload["member_email"].(string)

	// Determine UserID? (Loans are associated with Members, but Notifications are for Users)
	// For this SaaS, Members often correspond to a User account.
	// We'll assume the event provides a user_id if we have one, otherwise we might skip in-app for now.
	// But let's check if ev.AggregateID can be used or if payload has user_id.

	// 1. Save In-App Notification (if UserID is present)
	if userID, exists := payload["user_id"].(uint); exists {
		notif := &notificationModel.Notification{
			UserID:  userID,
			Title:   "Loan Approved",
			Message: fmt.Sprintf("Your loan application %s has been approved.", loanNum),
			Type:    "info",
		}
		notif.OrganizationID = ev.OrganizationID
		_ = s.repo.Create(context.Background(), notif)
		s.handlePushNotification(userID, notif.Title, notif.Message)
	}

	// 3. Enqueue Email Task
	emailBody := fmt.Sprintf(`
		<h3>Loan Approved</h3>
		<p>Great news! Your loan application <b>%s</b> has been approved.</p>
	`, loanNum)

	qPayload := &worker.PayloadSendEmail{
		To:      memberEmail,
		Subject: "Loan Approved: " + loanNum,
		Body:    emailBody,
	}
	if err := s.taskDist.DistributeTaskSendEmail(context.Background(), qPayload); err != nil {
		logrus.Errorf("Failed to enqueue email: %v", err)
	}
}

func (s *notificationService) handlePushNotification(userID uint, title, message string) {
	// 1. Get User's device tokens
	tokens, err := s.tokenRepo.GetByUserID(context.Background(), userID)
	if err != nil || len(tokens) == 0 {
		return
	}

	// 2. Deliver to each device
	for _, t := range tokens {
		_ = s.pushAdp.Send(context.Background(), t.DeviceToken, title, message)
	}

	logrus.Infof("[NotificationService] Sent Push to %d devices for User %d", len(tokens), userID)
}

func (s *notificationService) handleStockLow(evt event.Event) {
	payload := evt.Payload.(map[string]interface{})
	warehouseName := payload["warehouse_name"].(string)
	qty := payload["quantity"].(float64)
	
	title := "⚠️ Low Stock Alert"
	message := fmt.Sprintf("Stock in %s for Product #%d is low (%v). Please reorder soon.", 
		warehouseName, evt.AggregateID, qty)

	// Send to system notification stream
	notif := &notificationModel.Notification{
		UserID:  0, // System-wide
		Title:   title,
		Message: message,
		Type:    "warning",
		IsRead:  false,
	}
	notif.OrganizationID = evt.OrganizationID
	_ = s.repo.Create(context.Background(), notif)

	logrus.Infof("[NotificationService] Handled Stock Alert for Org %d: %s", evt.OrganizationID, message)
}
