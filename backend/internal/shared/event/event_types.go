package event

import "context"

// Event type constants for the Kafka event bus.
const (
	// IAM Events
	EventUserCreated = "user.created"
	EventUserUpdated = "user.updated"

	// Member Events
	EventMemberRegistered = "member.registered"
	EventMemberUpdated    = "member.updated"

	// Savings Events
	EventSavingsDeposited = "savings.deposited"
	EventSavingsWithdrawn = "savings.withdrawn"

	// Loan Events
	EventLoanApplied         = "loan.applied"
	EventLoanApproved        = "loan.approved"
	EventLoanRejected        = "loan.rejected"
	EventLoanDisbursed       = "loan.disbursed"
	EventLoanInstallmentPaid = "loan.installment.paid"
	EventLoanSettled         = "loan.settled"

	// Sales / POS Events
	EventSaleCompleted = "sale.completed"
	EventSaleReturned  = "sale.returned"

	// Inventory Events
	EventInventoryAdjusted = "inventory.adjusted"
	EventStockTransferred  = "inventory.stock.transferred"
	EventStockLow          = "inventory.stock.low"

	// Purchasing Events
	EventPurchaseOrderCreated = "purchase_order.created"
	EventGoodsReceived        = "purchase_order.goods_received"

	// Accounting Events
	EventJournalCreated = "journal.created"
)

// Event represents a domain event published to Kafka.
type Event struct {
	Type           string      `json:"type"`
	AggregateID    uint        `json:"aggregate_id"`
	OrganizationID uint        `json:"organization_id"`
	Payload        interface{} `json:"payload"`
	Timestamp      int64       `json:"timestamp"`
}

// SavingsTransactionPayload defines the payload for savings events.
type SavingsTransactionPayload struct {
	MemberID    uint    `json:"member_id"`
	Amount      float64 `json:"amount"`
	ProductCode string  `json:"product_code"`
	Description string  `json:"description"`
}

// LoanTransactionPayload defines the payload for loan events.
type LoanTransactionPayload struct {
	MemberID       uint    `json:"member_id"`
	Amount         float64 `json:"amount"`
	PrincipalPart  float64 `json:"principal_part"`
	InterestPart   float64 `json:"interest_part"`
	Description    string  `json:"description"`
}

// Publisher interface for publishing domain events.
type Publisher interface {
	Publish(ctx context.Context, evt Event) error
	Close() error
}

// Subscriber interface for consuming domain events.
type Subscriber interface {
	Consume(ctx context.Context, handler func(Event) error)
	Close() error
}
