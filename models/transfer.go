package models

import "time"

// TransferStatus represents the status of a transfer
type TransferStatus string

const (
	StatusPending    TransferStatus = "pending"
	StatusProcessing TransferStatus = "processing"
	StatusCompleted  TransferStatus = "completed"
	StatusFailed     TransferStatus = "failed"
	StatusCancelled  TransferStatus = "cancelled"
	StatusReversed   TransferStatus = "reversed"
)

// Transfer represents a points transfer between users
type Transfer struct {
	IdemKey       string          `json:"idemKey"`                  // Idempotency-Key (primary lookup)
	TransferID    *int            `json:"transferId,omitempty"`     // Internal ID (optional in response)
	FromUserID    int             `json:"fromUserId"`               // Sender user ID
	ToUserID      int             `json:"toUserId"`                 // Receiver user ID
	Amount        int             `json:"amount"`                   // Points amount
	Status        TransferStatus  `json:"status"`                   // Transfer status
	Note          *string         `json:"note,omitempty"`           // Optional note
	CreatedAt     time.Time       `json:"createdAt"`                // Created timestamp
	UpdatedAt     time.Time       `json:"updatedAt"`                // Updated timestamp
	CompletedAt   *time.Time      `json:"completedAt,omitempty"`    // Completed timestamp
	FailReason    *string         `json:"failReason,omitempty"`     // Failure reason if failed
}

// TransferCreateRequest represents the request to create a transfer
type TransferCreateRequest struct {
	FromUserID int     `json:"fromUserId" validate:"required,min=1"`
	ToUserID   int     `json:"toUserId" validate:"required,min=1"`
	Amount     int     `json:"amount" validate:"required,min=1"`
	Note       *string `json:"note,omitempty"`
}

// TransferCreateResponse wraps the created transfer
type TransferCreateResponse struct {
	Transfer Transfer `json:"transfer"`
}

// TransferGetResponse wraps a single transfer
type TransferGetResponse struct {
	Transfer Transfer `json:"transfer"`
}

// TransferListResponse wraps paginated transfer list
type TransferListResponse struct {
	Data     []Transfer `json:"data"`
	Page     int        `json:"page"`
	PageSize int        `json:"pageSize"`
	Total    int        `json:"total"`
}

// EventType represents the type of ledger event
type EventType string

const (
	EventTransferOut EventType = "transfer_out"
	EventTransferIn  EventType = "transfer_in"
	EventAdjust      EventType = "adjust"
	EventEarn        EventType = "earn"
	EventRedeem      EventType = "redeem"
)

// PointLedger represents a point transaction in the ledger
type PointLedger struct {
	ID           int        `json:"id"`
	UserID       int        `json:"userId"`
	Change       int        `json:"change"`         // Positive for receiving, negative for sending
	BalanceAfter int        `json:"balanceAfter"`   // Balance after this transaction
	EventType    EventType  `json:"eventType"`      // Type of event
	TransferID   *int       `json:"transferId,omitempty"` // Reference to transfer ID
	Reference    *string    `json:"reference,omitempty"`  // Additional reference
	Metadata     *string    `json:"metadata,omitempty"`   // JSON metadata
	CreatedAt    time.Time  `json:"createdAt"`
}
