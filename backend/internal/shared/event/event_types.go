package event

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
