package core

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/shopspring/decimal"
	tbTypes "github.com/tigerbeetle/tigerbeetle-go/pkg/types"
)

type TradeDirection string

const (
	TradeBuy  TradeDirection = "BUY"
	TradeSell TradeDirection = "SELL"
)

type FxTrade struct {
	// Stored in PG
	TbPendingId   tbTypes.Uint128
	CustomerId    uuid.UUID
	OperatorId    uuid.UUID
	ExchangeRate  decimal.Decimal
	Direction     TradeDirection
	CreatedAt     time.Time
	Notes         string
	IsUnusual     bool
	UnusualReason string

	// Stored in TB
	DebitLedger  uint32
	DebitAmount  uint64
	CreditLedger uint32
	CreditAmount uint64
}
