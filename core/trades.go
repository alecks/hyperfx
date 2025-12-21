package core

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/shopspring/decimal"
	tbTypes "github.com/tigerbeetle/tigerbeetle-go/pkg/types"
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
	DebitLedger  Ledger
	DebitAmount  uint64
	CreditLedger Ledger
	CreditAmount uint64
}

// TODO: implement
func (c *Core) GetRate(ledger Ledger) (decimal.Decimal, error) {
	return decimal.New(119, -2), nil
}

func (c *Core) LocalFromForeign(
	direction TradeDirection,
	foreignAmount decimal.Decimal,
	foreignLedger Ledger,
) (decimal.Decimal, error) {
	localScale := CurrencyAssetScales[c.options.LocalCurrencyLedger]

	rate, err := c.GetRate(foreignLedger)
	if err != nil {
		return decimal.Zero, err
	}

	if direction == TradeSell {
		// Ceil to nearest minor unit in our favour (see https://docs.tigerbeetle.com/single-page/#coding-recipes-currency-exchange)
		return foreignAmount.DivRound(rate, HfxPrecision).RoundCeil(localScale), nil
	} else {
		// Floor in our favour
		return foreignAmount.Mul(rate).RoundFloor(localScale), nil
	}
}

func (c *Core) ForeignFromLocal(
	direction TradeDirection,
	localAmount decimal.Decimal,
	foreignLedger Ledger,
) (decimal.Decimal, error) {
	foreignScale := CurrencyAssetScales[foreignLedger]

	rate, err := c.GetRate(foreignLedger)
	if err != nil {
		return decimal.Zero, err
	}

	if direction == TradeSell {
		return localAmount.Mul(rate).RoundFloor(foreignScale), nil
	} else {
		return localAmount.DivRound(rate, HfxPrecision).RoundCeil(foreignScale), nil
	}
}
