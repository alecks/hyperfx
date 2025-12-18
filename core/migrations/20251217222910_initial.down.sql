DROP TRIGGER IF EXISTS trg_block_unusual_customer ON fx_trades;
DROP FUNCTION IF EXISTS fn_block_unusual_customers();

DROP TABLE IF EXISTS midday_counts;
DROP TABLE IF EXISTS ledger_adjustments;
DROP TABLE IF EXISTS sod_openings;
DROP TABLE IF EXISTS eod_closings;
DROP TABLE IF EXISTS fx_trades;
DROP TABLE IF EXISTS customer_ledger_accounts;

DROP TABLE IF EXISTS customers;
DROP TABLE IF EXISTS operators;

