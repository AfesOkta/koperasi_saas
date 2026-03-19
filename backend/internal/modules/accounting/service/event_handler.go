package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/koperasi-gresik/backend/internal/modules/accounting/dto"
	"github.com/koperasi-gresik/backend/internal/shared/event"
)

// AccountingEventHandler listens to domain events and creates corresponding journal entries.
type AccountingEventHandler struct {
	accService AccountingService
	subscriber event.Subscriber
}

// NewAccountingEventHandler creates a new event handler.
func NewAccountingEventHandler(accService AccountingService, subscriber event.Subscriber) *AccountingEventHandler {
	return &AccountingEventHandler{
		accService: accService,
		subscriber: subscriber,
	}
}

// Start begins listening for events. Blocking call, run in goroutine.
func (h *AccountingEventHandler) Start(ctx context.Context) {
	log.Println("🎧 Accounting EventHandler started listening for domain events...")
	h.subscriber.Consume(ctx, h.handleEvent)
}

func (h *AccountingEventHandler) handleEvent(evt event.Event) error {
	ctx := context.Background()

	switch evt.Type {
	case event.EventSavingsDeposited:
		return h.handleSavingsDeposited(ctx, evt)
	case event.EventSavingsWithdrawn:
		return h.handleSavingsWithdrawn(ctx, evt)
	case event.EventLoanDisbursed:
		return h.handleLoanDisbursed(ctx, evt)
	case event.EventLoanInstallmentPaid:
		return h.handleLoanInstallmentPaid(ctx, evt)
	default:
		// Ignore unhandled events
		return nil
	}
}

func (h *AccountingEventHandler) handleSavingsDeposited(ctx context.Context, evt event.Event) error {
	var payload event.SavingsTransactionPayload
	if err := h.decodePayload(evt.Payload, &payload); err != nil {
		return err
	}

	// Double Entry: 
	// Debit: Kas (1101)
	// Credit: Simpanan Anggota (2301 for MVP, or dynamic based on product_code)
	
	req := dto.JournalEntryCreateRequest{
		Description:     fmt.Sprintf("Savings Deposit: %s", payload.Description),
		Date:            h.currentDate(),
		SourceModule:    "savings",
		SourceReference: fmt.Sprintf("DEP-%d", evt.AggregateID),
		Lines: []dto.JournalEntryLineRequest{
			{AccountCode: "1101", Debit: payload.Amount, Credit: 0, Description: "Kas Masuk Simpanan"},
			{AccountCode: "2301", Debit: 0, Credit: payload.Amount, Description: "Simpanan Anggota"},
		},
	}

	idempotencyKey := fmt.Sprintf("je.savings.dep.%d", evt.AggregateID)
	_, err := h.accService.CreateJournalEntryIdempotent(ctx, evt.OrganizationID, idempotencyKey, req)
	if err != nil {
		log.Printf("❌ Failed to create journal entry for SavingsDeposited: %v", err)
		return err
	}

	log.Printf("✅ Journal entry created for SavingsDeposited (AggregateID: %d)", evt.AggregateID)
	return nil
}

func (h *AccountingEventHandler) handleSavingsWithdrawn(ctx context.Context, evt event.Event) error {
	var payload event.SavingsTransactionPayload
	if err := h.decodePayload(evt.Payload, &payload); err != nil {
		return err
	}

	// Double Entry: 
	// Debit: Simpanan Anggota (2301)
	// Credit: Kas (1101)
	
	req := dto.JournalEntryCreateRequest{
		Description:     fmt.Sprintf("Savings Withdrawal: %s", payload.Description),
		Date:            h.currentDate(),
		SourceModule:    "savings",
		SourceReference: fmt.Sprintf("WDL-%d", evt.AggregateID),
		Lines: []dto.JournalEntryLineRequest{
			{AccountCode: "2301", Debit: payload.Amount, Credit: 0, Description: "Penarikan Simpanan Anggota"},
			{AccountCode: "1101", Debit: 0, Credit: payload.Amount, Description: "Kas Keluar"},
		},
	}

	idempotencyKey := fmt.Sprintf("je.savings.wdl.%d", evt.AggregateID)
	_, err := h.accService.CreateJournalEntryIdempotent(ctx, evt.OrganizationID, idempotencyKey, req)
	if err != nil {
		log.Printf("❌ Failed to create journal entry for SavingsWithdrawn: %v", err)
		return err
	}

	return nil
}

func (h *AccountingEventHandler) handleLoanDisbursed(ctx context.Context, evt event.Event) error {
	var payload event.LoanTransactionPayload
	if err := h.decodePayload(evt.Payload, &payload); err != nil {
		return err
	}

	// Double Entry: 
	// Debit: Piutang Anggota (1201)
	// Credit: Kas (1101)
	
	req := dto.JournalEntryCreateRequest{
		Description:     fmt.Sprintf("Loan Disbursement: %s", payload.Description),
		Date:            h.currentDate(),
		SourceModule:    "loan",
		SourceReference: fmt.Sprintf("DISB-%d", evt.AggregateID),
		Lines: []dto.JournalEntryLineRequest{
			{AccountCode: "1201", Debit: payload.Amount, Credit: 0, Description: "Piutang Anggota"},
			{AccountCode: "1101", Debit: 0, Credit: payload.Amount, Description: "Kas Keluar"},
		},
	}

	idempotencyKey := fmt.Sprintf("je.loan.disb.%d", evt.AggregateID)
	_, err := h.accService.CreateJournalEntryIdempotent(ctx, evt.OrganizationID, idempotencyKey, req)
	if err != nil {
		log.Printf("❌ Failed to create journal entry for LoanDisbursed: %v", err)
		return err
	}

	return nil
}

func (h *AccountingEventHandler) handleLoanInstallmentPaid(ctx context.Context, evt event.Event) error {
	var payload event.LoanTransactionPayload
	if err := h.decodePayload(evt.Payload, &payload); err != nil {
		return err
	}

	// Double Entry: 
	// Debit: Kas (1101) -> Total Amount
	// Credit: Piutang Anggota (1201) -> Principal Part
	// Credit: Pendapatan Bunga (4102) -> Interest Part
	
	req := dto.JournalEntryCreateRequest{
		Description:     fmt.Sprintf("Loan Installment Payment: %s", payload.Description),
		Date:            h.currentDate(),
		SourceModule:    "loan",
		SourceReference: fmt.Sprintf("LPAY-%d", evt.AggregateID),
		Lines: []dto.JournalEntryLineRequest{
			{AccountCode: "1101", Debit: payload.Amount, Credit: 0, Description: "Kas Masuk"},
			{AccountCode: "1201", Debit: 0, Credit: payload.PrincipalPart, Description: "Angsuran Pokok"},
			{AccountCode: "4102", Debit: 0, Credit: payload.InterestPart, Description: "Pendapatan Bunga"},
		},
	}

	idempotencyKey := fmt.Sprintf("je.loan.pay.%d", evt.AggregateID)
	_, err := h.accService.CreateJournalEntryIdempotent(ctx, evt.OrganizationID, idempotencyKey, req)
	if err != nil {
		log.Printf("❌ Failed to create journal entry for LoanInstallmentPaid: %v", err)
		return err
	}

	return nil
}

// Helper funcs
func (h *AccountingEventHandler) decodePayload(payload interface{}, out interface{}) error {
	// Payload from Kafka normally unmarshals as map[string]interface{}, convert back to json bytes and then to struct
	bytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, out)
}

func (h *AccountingEventHandler) currentDate() string {
	return time.Now().Format("2006-01-02")
}
