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
	foreignAmount uint64,
	foreignLedger Ledger,
) (uint64, error) {
	foreignScale := CurrencyAssetScales[foreignLedger]
	localScale := CurrencyAssetScales[c.options.LocalCurrencyLedger]

	rate, err := c.GetRate(foreignLedger)
	if err != nil {
		return 0, err
	}

	// Shift by the negative asset scale to get our 'standard form' decimal
	foreignDec := decimal.NewFromUint64(foreignAmount).Shift(-foreignScale)

	if direction == TradeSell {
		localDec := foreignDec.DivRound(rate, HfxPrecision)
		// Shift back to minor units and Ceil to nearest minor unit in our favour (see https://docs.tigerbeetle.com/single-page/#coding-recipes-currency-exchange)
		return uint64(localDec.Shift(localScale).Ceil().IntPart()), nil

	} else {
		localDec := foreignDec.Mul(rate)
		// Shift back to minor units and Floor in our favour
		return uint64(localDec.Shift(localScale).Floor().IntPart()), nil
	}
}

func (c *Core) ForeignFromLocal(
	direction TradeDirection,
	localAmount uint64,
	foreignLedger Ledger,
) (uint64, error) {
	localScale := CurrencyAssetScales[c.options.LocalCurrencyLedger]
	foreignScale := CurrencyAssetScales[foreignLedger]

	rate, err := c.GetRate(foreignLedger)
	if err != nil {
		return 0, err
	}

	localDec := decimal.NewFromUint64(localAmount).Shift(-localScale)

	if direction == TradeSell {
		foreignDec := localDec.Mul(rate)
		return uint64(foreignDec.Shift(foreignScale).Floor().IntPart()), nil
	} else {
		foreignDec := localDec.DivRound(rate, HfxPrecision)
		return uint64(foreignDec.Shift(foreignScale).Ceil().IntPart()), nil
	}
}
