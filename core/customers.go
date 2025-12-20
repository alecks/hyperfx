package core

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
	tbTypes "github.com/tigerbeetle/tigerbeetle-go/pkg/types"
)

type Customer struct {
	Id        uuid.UUID
	FullName  string
	CreatedAt time.Time
	Address   string
	Postcode  string

	IsBlocked     bool
	BlockedReason bool
}

type CreateCustomerData struct {
	FullName string
	Address  string
	Postcode string
}

// CreateCustomer inserts a customer into the PG database. It does not create a TB account.
func (c *Core) CreateCustomer(ctx context.Context, data CreateCustomerData) (*Customer, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	var createdAt time.Time
	if err := c.pgc.QueryRow(ctx, "INSERT INTO customers (id, full_name, address, postcode) VALUES ($1, $2, $3, $4) RETURNING created_at", id, data.FullName, data.Address, data.Postcode).
		Scan(&createdAt); err != nil {
		return nil, err
	}

	return &Customer{
		Id:        id,
		FullName:  data.FullName,
		CreatedAt: createdAt,
		Address:   data.Address,
		Postcode:  data.Postcode,
	}, nil
}

func (c *Core) GetCustomerById(ctx context.Context, id uuid.UUID) (*Customer, error) {
	cust := &Customer{}
	err := c.pgc.QueryRow(ctx, "SELECT id, full_name, created_at, address, postcode, is_blocked, COALESCE(blocked_reason, '') FROM customers WHERE id = $1", id).
		Scan(&cust.Id, &cust.FullName, &cust.CreatedAt, &cust.Address, &cust.Postcode, &cust.IsBlocked, &cust.BlockedReason)

	return cust, err
}

type CustomerLedgerAccount struct {
	CustomerId  uuid.UUID
	Ledger      uint32
	TbAccountId tbTypes.Uint128
}

// AddCustomerLedgerAccount stores a customer's TB account ID for a certain ledger in PG.
func (c *Core) AddCustomerLedgerAccount(ctx context.Context, data CustomerLedgerAccount) error {
	// TODO: maybe add a helper function to make this look less awful
	tbAccountUuid, err := uuid.FromBytes(data.TbAccountId[:])
	if err != nil {
		return err
	}

	_, err = c.pgc.Exec(
		ctx,
		"INSERT INTO customer_ledger_accounts (customer_id, ledger_id, tb_account_id) VALUES ($1, $2, $3)",
		data.CustomerId,
		data.Ledger,
		tbAccountUuid,
	)
	return err
}
