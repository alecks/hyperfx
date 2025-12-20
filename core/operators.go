package core

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Operator struct {
	Id           uuid.UUID
	Username     string
	PasswordHash []byte
	CreatedAt    time.Time
	IsActive     bool
}

type CreateOperatorData struct {
	Username string
	Password []byte
}

// CreateOperator inserts a new operator into the database, returning the newly created operator.
func (c *Core) CreateOperator(ctx context.Context, data CreateOperatorData) (*Operator, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	passHash, err := bcrypt.GenerateFromPassword(data.Password, bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	var createdAt time.Time
	var isActive bool
	if err := c.pgc.QueryRow(ctx, "INSERT INTO operators (id, username, password_hash) VALUES ($1, $2, $3) RETURNING created_at, is_active", id, data.Username, passHash).
		Scan(&createdAt, &isActive); err != nil {
		return nil, err
	}

	return &Operator{
		Id:           id,
		Username:     data.Username,
		PasswordHash: passHash,
		CreatedAt:    createdAt,
		IsActive:     isActive,
	}, nil
}

// GetOperator queries the database for an operator with a matching id OR username.
func (c *Core) GetOperator(ctx context.Context, id uuid.UUID, username string) (*Operator, error) {
	op := &Operator{}
	if err := c.pgc.QueryRow(ctx, "SELECT id, username, password_hash, created_at, is_active FROM operators WHERE id = $1 OR username = $2", id, username).
		Scan(&op.Id, &op.Username, &op.PasswordHash, &op.CreatedAt, &op.IsActive); err != nil {
		return nil, err
	}

	return op, nil
}

// SetOperatorActivated sets the is_active field of an operator to the given value.
func (c *Core) SetOperatorActivated(ctx context.Context, id uuid.UUID, isActive bool) error {
	if _, err := c.pgc.Exec(
		ctx,
		"UPDATE operators SET is_active = $1 WHERE id = $2",
		isActive,
		id,
	); err != nil {
		return err
	}

	return nil
}

// VerifyOperatorPassword queries the database for an operator's password hash and checks against the password given.
// Returns nil on success, error otherwise.
func (c *Core) VerifyOperatorPassword(ctx context.Context, id uuid.UUID, password []byte) error {
	var hash []byte
	if err := c.pgc.QueryRow(ctx, "SELECT password_hash FROM operators WHERE id = $1", id).
		Scan(&hash); err != nil {
		return err
	}

	return bcrypt.CompareHashAndPassword(hash, password)
}

// SetOperatorPassword hashes a new password for an operator and updates the database.
// You may want to call VerifyOperatorPassword before this.
func (c *Core) SetOperatorPassword(ctx context.Context, id uuid.UUID, password []byte) error {
	hash, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = c.pgc.Exec(ctx, "UPDATE operators SET password_hash = $1 WHERE id = $2", hash, id)
	if err != nil {
		return err
	}

	return nil
}
