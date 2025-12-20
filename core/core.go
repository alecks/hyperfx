package core

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/gofrs/uuid"

	"github.com/golang-migrate/migrate/v4"
	migratePgx "github.com/golang-migrate/migrate/v4/database/pgx"
	"github.com/golang-migrate/migrate/v4/source/iofs"

	pgxuuid "github.com/jackc/pgx-gofrs-uuid"
	pgxdecimal "github.com/jackc/pgx-shopspring-decimal"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"

	tb "github.com/tigerbeetle/tigerbeetle-go"
	tbTypes "github.com/tigerbeetle/tigerbeetle-go/pkg/types"
)

const (
	FILENAME_HFX_DIR   = ".hyperfx"
	FILENAME_NAMESPACE = "namespace"
)

//go:embed migrations/*.sql
var migrationsFs embed.FS

type Core struct {
	tbc       tb.Client
	pgc       *pgxpool.Pool
	options   Options
	namespace uuid.UUID

	ids    *knownIds
	Logger *slog.Logger
}

type Options struct {
	TbAddresses []string
	TbClusterId tbTypes.Uint128

	PgUrl string

	HfxDir              string
	LocalCurrencyLedger uint32
}

type knownIds struct {
	// Map from currency code to account ID
	liquidity map[uint32]tbTypes.Uint128
	overs     map[uint32]tbTypes.Uint128
	shorts    map[uint32]tbTypes.Uint128
	control   map[uint32]tbTypes.Uint128
	fees      tbTypes.Uint128 // Only local currency
}

// New creates a new HyperFX instance and creates connections to PG and TB.
// You must call Close to close database connections and shut down HyperFX, unless an error is returned.
func New(ctx context.Context, logger *slog.Logger, options Options) (*Core, error) {
	// NOTE: this aims to never have a broken Core struct at any point, hence the verbosity

	if options.TbAddresses == nil || options.PgUrl == "" ||
		options.LocalCurrencyLedger == 0 {
		return nil, errors.New("core: TbAddresses, PgUrl and LocalCurrencyLedger are required")
	}

	if options.HfxDir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf(
				"core: failed to get user home dir: %w. Specify Options.HfxDir to override use of home directory",
				err,
			)
		}
		options.HfxDir = filepath.Join(homeDir, FILENAME_HFX_DIR)
	}

	namespace, err := loadNamespace(options, logger)
	if err != nil {
		return nil, err
	}

	tbc, pgc, err := connectDatabases(ctx, options, logger)
	if err != nil {
		return nil, err
	}

	if err := migratePg(ctx, pgc, logger); err != nil {
		tbc.Close()
		pgc.Close()
		return nil, err
	}

	ids, err := initSystemAccounts(options, namespace, tbc, logger)
	if err != nil {
		tbc.Close()
		pgc.Close()
		return nil, err
	}

	return &Core{
		tbc:       tbc,
		pgc:       pgc,
		options:   options,
		namespace: namespace,
		ids:       ids,
		Logger:    logger,
	}, nil
}

// Close shuts down HyperFX gracefully, including closing DB connections.
func (c *Core) Close() {
	c.tbc.Close()
	c.pgc.Close()
}

// loadNamespace either gets the existing namespace from a file or generates a new one and writes it.
func loadNamespace(options Options, logger *slog.Logger) (uuid.UUID, error) {
	// TODO: not sure if the core should really do this, maybe leave it to the user
	var namespace uuid.UUID
	namespaceFilepath := filepath.Join(options.HfxDir, FILENAME_NAMESPACE)
	namespaceBytes, err := os.ReadFile(namespaceFilepath)

	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			os.Mkdir(options.HfxDir, 0o700)
			namespace, err = uuid.NewV4()
			if err != nil {
				return uuid.Nil, err
			}
			if err := os.WriteFile(namespaceFilepath, namespace.Bytes(), 0o600); err != nil {
				return uuid.Nil, fmt.Errorf(
					"core: failed to write namespace file (check permissions?): %w",
					err,
				)
			}

			logger.Warn(
				"namespace file did not exist, new namespace generated and written",
				"filepath",
				namespaceFilepath,
				"namespace",
				namespace.String(),
			)
		} else {
			return uuid.Nil, fmt.Errorf(
				"core: failed to read namespace file (check permissions?): %w",
				err,
			)
		}
	} else {
		namespace, err = uuid.FromBytes(namespaceBytes)
		if err != nil {
			return uuid.Nil, err
		}

		logger.Info(
			"existing namespace file found",
			"filepath",
			namespaceFilepath,
			"namespace",
			namespace.String(),
		)
	}

	return namespace, nil
}

// connectDatabases creates TB and PG clients and tests connection. This function will close database connections itself on error.
func connectDatabases(
	ctx context.Context,
	options Options,
	logger *slog.Logger,
) (tb.Client, *pgxpool.Pool, error) {
	logger.Debug(
		"creating TB client",
		"cluster_id",
		options.TbClusterId,
		"addresses",
		options.TbAddresses,
	)
	tbc, err := tb.NewClient(options.TbClusterId, options.TbAddresses)
	if err != nil {
		return nil, nil, fmt.Errorf("core: failed to create TB client: %w", err)
	}

	logger.Debug(
		"performing TB Nop request",
		"cluster_id",
		options.TbClusterId,
		"addresses",
		options.TbAddresses,
	)
	if err := tbc.Nop(); err != nil {
		tbc.Close()
		return nil, nil, fmt.Errorf("core: failed to send Nop to TB: %w", err)
	}
	logger.Info(
		"TB connection is UP",
		"cluster_id",
		options.TbClusterId,
		"addresses",
		options.TbAddresses,
	)

	logger.Debug("creating PG pool", "url", options.PgUrl)
	pgConf, err := pgxpool.ParseConfig(options.PgUrl)
	if err != nil {
		tbc.Close()
		return nil, nil, fmt.Errorf("core: failed to parse PG url: %w", err)
	}
	// Allow pgx to be used with gofrs UUID and shopspring decimals.
	pgConf.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		typeMap := conn.TypeMap()
		pgxuuid.Register(typeMap)
		pgxdecimal.Register(typeMap)
		return nil
	}

	pgc, err := pgxpool.NewWithConfig(ctx, pgConf)
	if err != nil {
		tbc.Close()
		return nil, nil, fmt.Errorf("core: failed to create PG connection pool: %w", err)
	}

	logger.Debug("pinging PG", "url", options.PgUrl)
	if err := pgc.Ping(ctx); err != nil {
		tbc.Close()
		pgc.Close()
		return nil, nil, fmt.Errorf("core: failed to ping PG: %w", err)
	}

	logger.Info("PG connection is UP", "url", options.PgUrl)
	return tbc, pgc, nil
}

// migratePg runs migrations on an existing postgres pool.
func migratePg(ctx context.Context, pgc *pgxpool.Pool, logger *slog.Logger) error {
	db := stdlib.OpenDBFromPool(pgc)

	dbDriver, err := migratePgx.WithInstance(db, &migratePgx.Config{})
	if err != nil {
		db.Close()
		return fmt.Errorf("core: failed to create pgx driver for PG migrations: %w", err)
	}

	sourceDriver, err := iofs.New(migrationsFs, "migrations")
	if err != nil {
		dbDriver.Close()
		return fmt.Errorf("core: failed to create iofs driver for PG migrations: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", sourceDriver, "pgx5", dbDriver)
	if err != nil {
		dbDriver.Close()
		sourceDriver.Close()
		return fmt.Errorf("core: failed to create PG migrate instance: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			logger.Debug("PG migrations: no change, up to date")
		} else {
			return fmt.Errorf("core: PG migrations failed: %w", err)
		}
	}
	logger.Info("PG migrations complete, PG ready for operation")
	return nil
}

// initSystemAccounts gets Tigerbeetle system account IDs and ensures system accounts exist.
func initSystemAccounts(
	options Options,
	namespace uuid.UUID,
	tbc tb.Client,
	logger *slog.Logger,
) (*knownIds, error) {
	ids := &knownIds{
		liquidity: map[uint32]tbTypes.Uint128{},
		overs:     map[uint32]tbTypes.Uint128{},
		shorts:    map[uint32]tbTypes.Uint128{},
		control:   map[uint32]tbTypes.Uint128{},
	}
	accountCreationBatch := []tbTypes.Account{}

	// LIQUIDITY, DISCREPANCY AND CONTROL ACCOUNTS
	for currCode := range CurrencyAssetScales {
		liqKey := fmt.Sprintf("branch_liquidity_%d", currCode)
		liqId := idWithNamespace(namespace, liqKey)

		ids.liquidity[currCode] = liqId
		accountCreationBatch = append(accountCreationBatch, tbTypes.Account{
			ID:     liqId,
			Ledger: currCode,
			Code:   AccountCodeBranchLiquidity,
			Flags: tbTypes.AccountFlags{
				History: true,
			}.ToUint16(),
		})

		oversKey := fmt.Sprintf("branch_overs_%d", currCode)
		oversId := idWithNamespace(namespace, oversKey)

		ids.overs[currCode] = oversId
		accountCreationBatch = append(accountCreationBatch, tbTypes.Account{
			ID:     oversId,
			Ledger: currCode,
			Code:   AccountCodeBranchOvers,
			Flags: tbTypes.AccountFlags{
				DebitsMustNotExceedCredits: true,
				History:                    true,
			}.ToUint16(),
		})

		shortsKey := fmt.Sprintf("branch_shorts_%d", currCode)
		shortsId := idWithNamespace(namespace, shortsKey)

		ids.shorts[currCode] = shortsId
		accountCreationBatch = append(accountCreationBatch, tbTypes.Account{
			ID:     shortsId,
			Ledger: currCode,
			Code:   AccountCodeBranchShorts,
			Flags: tbTypes.AccountFlags{
				CreditsMustNotExceedDebits: true,
				History:                    true,
			}.ToUint16(),
		})

		controlKey := fmt.Sprintf("branch_control_%d", currCode)
		controlId := idWithNamespace(namespace, controlKey)

		ids.control[currCode] = controlId
		accountCreationBatch = append(accountCreationBatch, tbTypes.Account{
			ID:     controlId,
			Ledger: currCode,
			Code:   AccountCodeBranchControl,
			Flags: tbTypes.AccountFlags{
				History: true,
			}.ToUint16(),
		})
	}

	// FEES ACCOUNT
	feesId := idWithNamespace(namespace, "branch_fees")
	accountCreationBatch = append(accountCreationBatch, tbTypes.Account{
		ID:     feesId,
		Ledger: options.LocalCurrencyLedger,
		Code:   AccountCodeBranchFees,
		Flags: tbTypes.AccountFlags{
			DebitsMustNotExceedCredits: true,
			History:                    true,
		}.ToUint16(),
	})

	var exists, failures int
	accountErrors, err := tbc.CreateAccounts(accountCreationBatch)
	if err != nil {
		return nil, fmt.Errorf("core: failed to send create accounts request to TB: %w", err)
	}
	for _, e := range accountErrors {
		account := accountCreationBatch[e.Index]
		switch e.Result {
		case tbTypes.AccountOK:
			// NOTE: AccountOK shouldn't be in the list of errors, but just to be sure.
		case tbTypes.AccountExists:
			// NOTE: there are other exists errors, but these would only occur if we changed the code, flags or user_data, or mismatched the ledger.
			exists++
			logger.Debug(
				"account creation: exists",
				"id",
				account.ID,
				"ledger",
				account.Ledger,
				"code",
				account.Code,
			)
		default:
			failures++
			logger.Error(
				"account creation: "+e.Result.String(),
				"id",
				account.ID,
				"ledger",
				account.Ledger,
				"code",
				account.Code,
			)
		}
	}

	logger.Info(
		"account creation requests complete, TB ready for operation",
		"total_requests",
		len(accountCreationBatch),
		"exists_occurences",
		exists,
		"failure_occurences",
		failures,
	)

	return ids, nil
}

// idWithNamespace generates a UUID v5 using the given namespace and a 'name' string.
// Can be used to create deterministic IDs for frequently used accounts, while still adhering to https://docs.tigerbeetle.com/reference/account/#id.
func idWithNamespace(namespace uuid.UUID, name string) tbTypes.Uint128 {
	idBytes := [16]byte(uuid.NewV5(namespace, name).Bytes())
	return tbTypes.BytesToUint128(idBytes)
}

// uuidToTb converts a uuid.UUID to a tbTypes.Uint128.
func uuidToTb(u uuid.UUID) tbTypes.Uint128 {
	return tbTypes.BytesToUint128([16]byte(u.Bytes()))
}

func tbToUuid(i tbTypes.Uint128) (uuid.UUID, error) {
	return uuid.FromBytes(i[:])
}
