package core

// Tigerbeetle account codes
const (
	AccountCodeBranchLiquidity uint16 = 1000
	AccountCodeBranchFees      uint16 = 1001
	AccountCodeBranchOvers     uint16 = 2000
	AccountCodeBranchShorts    uint16 = 2001
	AccountCodeBranchControl   uint16 = 9000

	AccountCodeCustomer uint16 = 3000
)

// Currency/ledger codes. Starting with a basic set from ISO 4217.
const (
	LedgerGBP uint32 = 826 // Great British Pound
	LedgerUSD uint32 = 840 // United States Dollar
	LedgerEUR uint32 = 978 // Euro

	LedgerJPY uint32 = 392 // Japanese Yen
	LedgerCAD uint32 = 124 // Canadian Dollar
	LedgerAUD uint32 = 036 // Australian Dollar
	LedgerCHF uint32 = 756 // Swiss Franc
	LedgerCNY uint32 = 156 // Chinese Yuan Renminbi
	LedgerHKD uint32 = 344 // Hong Kong Dollar
	LedgerNZD uint32 = 554 // New Zealand Dollar
	LedgerSEK uint32 = 752 // Swedish Krona
	LedgerNOK uint32 = 578 // Norwegian Krone
	LedgerDKK uint32 = 208 // Danish Krone
	LedgerSGD uint32 = 702 // Singapore Dollar
	LedgerINR uint32 = 356 // Indian Rupee
	LedgerMXN uint32 = 484 // Mexican Peso
	LedgerBRL uint32 = 986 // Brazilian Real
	LedgerZAR uint32 = 710 // South African Rand
	LedgerRUB uint32 = 643 // Russian Ruble
	LedgerKRW uint32 = 410 // South Korean Won
	LedgerTRY uint32 = 949 // Turkish Lira
	LedgerPLN uint32 = 985 // Polish Zloty
	LedgerTHB uint32 = 764 // Thai Baht
	LedgerIDR uint32 = 360 // Indonesian Rupiah
	LedgerMYR uint32 = 458 // Malaysian Ringgit
	LedgerPHP uint32 = 608 // Philippine Peso
	LedgerVND uint32 = 704 // Vietnamese Dong
	LedgerEGP uint32 = 818 // Egyptian Pound
	LedgerNGN uint32 = 566 // Nigerian Naira
	LedgerKES uint32 = 404 // Kenyan Shilling
	LedgerUAH uint32 = 980 // Ukrainian Hryvnia
	LedgerCLP uint32 = 152 // Chilean Peso
	LedgerCOP uint32 = 170 // Colombian Peso
	LedgerPEN uint32 = 604 // Peruvian Sol
	LedgerARS uint32 = 032 // Argentine Peso
	LedgerSAR uint32 = 682 // Saudi Riyal
	LedgerAED uint32 = 784 // UAE Dirham
	LedgerKWD uint32 = 414 // Kuwaiti Dinar
	LedgerQAR uint32 = 634 // Qatari Rial
)

// See https://docs.tigerbeetle.com/coding/data-modeling/.
var CurrencyAssetScales = map[uint32]uint8{
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
