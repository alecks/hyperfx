package core

// HfxPrecision represents the decimal precision that should be used throughout the core.
// TODO: make this configurable.
const HfxPrecision = 9

// TradeDirection represents whether WE are buying a foreign currency from the customer or selling to the customer.
// This is purely stored in PG for easy querying.
type TradeDirection string

const (
	TradeBuy  TradeDirection = "BUY"
	TradeSell TradeDirection = "SELL"
)

// AccountCode represents a valid TB Account.code field (uint16).
type AccountCode uint16

const (
	AccountCodeBranchLiquidity AccountCode = 1000
	AccountCodeBranchFees      AccountCode = 1001
	AccountCodeBranchOvers     AccountCode = 2000
	AccountCodeBranchShorts    AccountCode = 2001
	AccountCodeBranchControl   AccountCode = 9000

	AccountCodeCustomer AccountCode = 3000
)

// Ledger represents a valid currency for the TB Account.ledger field (uint32).
type Ledger uint32

// ISO 4217
// TODO: maybe make this dynamic
const (
	LedgerGBP Ledger = 826 // Great British Pound
	LedgerUSD Ledger = 840 // United States Dollar
	LedgerEUR Ledger = 978 // Euro

	LedgerJPY Ledger = 392  // Japanese Yen
	LedgerCAD Ledger = 124  // Canadian Dollar
	LedgerAUD Ledger = 0o36 // Australian Dollar
	LedgerCHF Ledger = 756  // Swiss Franc
	LedgerCNY Ledger = 156  // Chinese Yuan Renminbi
	LedgerHKD Ledger = 344  // Hong Kong Dollar
	LedgerNZD Ledger = 554  // New Zealand Dollar
	LedgerSEK Ledger = 752  // Swedish Krona
	LedgerNOK Ledger = 578  // Norwegian Krone
	LedgerDKK Ledger = 208  // Danish Krone
	LedgerSGD Ledger = 702  // Singapore Dollar
	LedgerINR Ledger = 356  // Indian Rupee
	LedgerMXN Ledger = 484  // Mexican Peso
	LedgerBRL Ledger = 986  // Brazilian Real
	LedgerZAR Ledger = 710  // South African Rand
	LedgerRUB Ledger = 643  // Russian Ruble
	LedgerKRW Ledger = 410  // South Korean Won
	LedgerTRY Ledger = 949  // Turkish Lira
	LedgerPLN Ledger = 985  // Polish Zloty
	LedgerTHB Ledger = 764  // Thai Baht
	LedgerIDR Ledger = 360  // Indonesian Rupiah
	LedgerMYR Ledger = 458  // Malaysian Ringgit
	LedgerPHP Ledger = 608  // Philippine Peso
	LedgerVND Ledger = 704  // Vietnamese Dong
	LedgerEGP Ledger = 818  // Egyptian Pound
	LedgerNGN Ledger = 566  // Nigerian Naira
	LedgerKES Ledger = 404  // Kenyan Shilling
	LedgerUAH Ledger = 980  // Ukrainian Hryvnia
	LedgerCLP Ledger = 152  // Chilean Peso
	LedgerCOP Ledger = 170  // Colombian Peso
	LedgerPEN Ledger = 604  // Peruvian Sol
	LedgerARS Ledger = 0o32 // Argentine Peso
	LedgerSAR Ledger = 682  // Saudi Riyal
	LedgerAED Ledger = 784  // UAE Dirham
	LedgerKWD Ledger = 414  // Kuwaiti Dinar
	LedgerQAR Ledger = 634  // Qatari Rial
)

// CurrencyAssetScales represents how much the minor unit of a currency (stored in TB) needs to be scaled up to get to the usual 'display' unit.
// See https://docs.tigerbeetle.com/coding/data-modeling/.
var CurrencyAssetScales = map[Ledger]int32{
	LedgerGBP: 2, // 1 GBP = 100 pence
	LedgerUSD: 2, // 1 USD = 100 cents
	LedgerEUR: 2, // 1 EUR = 100 cents

	LedgerJPY: 0, // no minor unit
	LedgerCAD: 2, // 1 CAD = 100 cents
	LedgerAUD: 2, // 1 AUD = 100 cents
	LedgerCHF: 2, // 1 CHF = 100 rappen
	LedgerCNY: 2, // 1 CNY = 100 fen
	LedgerHKD: 2, // 1 HKD = 100 cents
	LedgerNZD: 2, // 1 NZD = 100 cents
	LedgerSEK: 2, // 1 SEK = 100 ore
	LedgerNOK: 2, // 1 NOK = 100 ore
	LedgerDKK: 2, // 1 DKK = 100 ore
	LedgerSGD: 2, // 1 SGD = 100 cents
	LedgerINR: 2, // 1 INR = 100 paise
	LedgerMXN: 2, // 1 MXN = 100 centavos
	LedgerBRL: 2, // 1 BRL = 100 centavos
	LedgerZAR: 2, // 1 ZAR = 100 cents
	LedgerRUB: 2, // 1 RUB = 100 kopecks
	LedgerKRW: 0, // no minor unit
	LedgerTRY: 2, // 1 TRY = 100 kurus
	LedgerPLN: 2, // 1 PLN = 100 groszy
	LedgerTHB: 2, // 1 THB = 100 satang
	LedgerIDR: 2, // 1 IDR = 100 sen
	LedgerMYR: 2, // 1 MYR = 100 sen
	LedgerPHP: 2, // 1 PHP = 100 centavos
	LedgerVND: 0, // no minor unit
	LedgerEGP: 2, // 1 EGP = 100 piastres
	LedgerNGN: 2, // 1 NGN = 100 kobo
	LedgerKES: 2, // 1 KES = 100 cents
	LedgerUAH: 2, // 1 UAH = 100 kopiykas
	LedgerCLP: 0, // no minor unit
	LedgerCOP: 2, // 1 COP = 100 centavos
	LedgerPEN: 2, // 1 PEN = 100 centimos
	LedgerARS: 2, // 1 ARS = 100 centavos
	LedgerSAR: 2, // 1 SAR = 100 halalas
	LedgerAED: 2, // 1 AED = 100 fils
	LedgerKWD: 3, // 1 KWD = 1000 fils
	LedgerQAR: 2, // 1 QAR = 100 dirhams
}
