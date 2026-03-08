package statements

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

// Validator handles input validation for statement generation
type Validator struct {
	currencyCodes map[string]bool
	postalRegex   *regexp.Regexp
}

// NewValidator creates a new validator instance
func NewValidator() *Validator {
	return &Validator{
		currencyCodes: initCurrencyCodes(),
		postalRegex:   regexp.MustCompile(`^[A-Z0-9\s-]{3,10}$`),
	}
}

// ValidateStatementInput validates the complete statement input
func (v *Validator) ValidateStatementInput(input StatementInput) error {
	errors := &ValidationErrors{}

	// Validate account
	if err := v.validateAccount(input.Account); err != nil {
		if ve, ok := err.(ValidationError); ok {
			errors.Add(ve.Field, ve.Message)
		}
	}

	// Validate period
	if input.PeriodStart.IsZero() {
		errors.Add("period_start", "period start date is required")
	}
	if input.PeriodEnd.IsZero() {
		errors.Add("period_end", "period end date is required")
	}
	if !input.PeriodEnd.After(input.PeriodStart) {
		errors.Add("period", "period end must be after period start")
	}

	// Validate transactions
	transactionIDs := make(map[string]bool)
	for i, txn := range input.Transactions {
		if err := v.validateTransaction(txn, input.PeriodStart, input.PeriodEnd, i); err != nil {
			if ve, ok := err.(ValidationError); ok {
				errors.Add(ve.Field, ve.Message)
			}
		}

		// Check for duplicate IDs
		if transactionIDs[txn.ID] {
			errors.Add(fmt.Sprintf("transactions[%d].id", i),
				fmt.Sprintf("duplicate transaction ID: %s", txn.ID))
		}
		transactionIDs[txn.ID] = true
	}

	// Validate institution if provided
	if input.Institution != nil {
		if err := v.validateInstitution(*input.Institution); err != nil {
			if ve, ok := err.(ValidationError); ok {
				errors.Add(ve.Field, ve.Message)
			}
		}
	}

	if errors.HasErrors() {
		return errors
	}

	return nil
}

// validateAccount validates account information
func (v *Validator) validateAccount(account Account) error {
	if account.Number == "" {
		return NewValidationError("account.number", "account number is required")
	}
	if account.HolderName == "" {
		return NewValidationError("account.holder_name", "holder name is required")
	}
	if account.Currency == "" {
		return NewValidationError("account.currency", "currency is required")
	}
	if err := v.ValidateCurrencyCode(account.Currency); err != nil {
		return NewValidationError("account.currency", err.Error())
	}
	if account.Address != nil {
		if err := v.ValidateAddress(account.Address); err != nil {
			return err
		}
	}
	return nil
}

// validateTransaction validates a single transaction
func (v *Validator) validateTransaction(txn Transaction, periodStart, periodEnd time.Time, index int) error {
	prefix := fmt.Sprintf("transactions[%d]", index)

	if txn.ID == "" {
		return NewValidationError(prefix+".id", "transaction ID is required")
	}
	if txn.Description == "" {
		return NewValidationError(prefix+".description", "transaction description is required")
	}
	if txn.Date.IsZero() {
		return NewValidationError(prefix+".date", "transaction date is required")
	}

	// Check if transaction is within period
	if txn.Date.Before(periodStart) || txn.Date.After(periodEnd) {
		return NewValidationError(prefix+".date",
			fmt.Sprintf("transaction date %s outside statement period", txn.Date.Format("2006-01-02")))
	}

	// Validate transaction type
	if txn.Type != Credit && txn.Type != Debit {
		return NewValidationError(prefix+".type",
			fmt.Sprintf("invalid transaction type: %s", txn.Type))
	}

	// Validate amount based on type
	if txn.Type == Credit && txn.Amount.IsNegative() {
		return NewValidationError(prefix+".amount",
			"credit transaction cannot have negative amount")
	}
	if txn.Type == Debit && txn.Amount.IsPositive() {
		return NewValidationError(prefix+".amount",
			"debit transaction must have negative amount")
	}

	return nil
}

// validateInstitution validates institution information
func (v *Validator) validateInstitution(inst Institution) error {
	if inst.Name == "" {
		return NewValidationError("institution.name", "institution name is required")
	}
	if inst.Address != nil {
		if err := v.ValidateAddress(inst.Address); err != nil {
			return err
		}
	}
	return nil
}

// ValidateCurrencyCode validates an ISO 4217 currency code
func (v *Validator) ValidateCurrencyCode(code string) error {
	if len(code) != 3 {
		return fmt.Errorf("invalid currency code: must be 3 characters")
	}

	upperCode := strings.ToUpper(code)
	if upperCode != code {
		return fmt.Errorf("invalid currency code: must be uppercase")
	}

	if !v.currencyCodes[code] {
		return fmt.Errorf("invalid currency code: %s is not a valid ISO 4217 code", code)
	}

	return nil
}

// ValidateAddress validates an address
func (v *Validator) ValidateAddress(address *Address) error {
	if address == nil {
		return nil
	}

	if address.Line1 == "" {
		return NewValidationError("address.line1", "address line 1 is required")
	}
	if len(address.Line1) > 200 {
		return NewValidationError("address.line1", "address line 1 too long (max 200 characters)")
	}

	if address.City == "" {
		return NewValidationError("address.city", "city is required")
	}

	if address.Country == "" {
		return NewValidationError("address.country", "country is required")
	}

	// Validate postal code format if provided
	if address.PostalCode != "" && len(address.PostalCode) > 20 {
		return NewValidationError("address.postal_code", "invalid postal code format")
	}

	return nil
}

// ValidateBalanceConsistency validates that provided balances are consistent
func (v *Validator) ValidateBalanceConsistency(transactions []Transaction, openingBalance decimal.Decimal) error {
	runningBalance := openingBalance

	for i, txn := range transactions {
		runningBalance = runningBalance.Add(txn.Amount)

		if txn.Balance != nil {
			// Allow small rounding differences (up to 0.01)
			diff := runningBalance.Sub(*txn.Balance).Abs()
			tolerance := decimal.NewFromFloat(0.01)

			if diff.GreaterThan(tolerance) {
				return NewCalculationError(fmt.Sprintf(
					"balance mismatch for transaction %s at index %d: expected %s, provided %s",
					txn.ID, i, runningBalance.String(), txn.Balance.String()))
			}
		}
	}

	return nil
}

// ValidateProvidedBalances checks if provided balances match calculated ones
func (v *Validator) ValidateProvidedBalances(transactions []Transaction, openingBalance decimal.Decimal) error {
	return v.ValidateBalanceConsistency(transactions, openingBalance)
}

// initCurrencyCodes initializes the set of valid ISO 4217 currency codes
func initCurrencyCodes() map[string]bool {
	codes := map[string]bool{
		// Major currencies
		"USD": true, "EUR": true, "GBP": true, "JPY": true, "CHF": true,
		"CAD": true, "AUD": true, "NZD": true, "CNY": true, "HKD": true,
		"SGD": true, "SEK": true, "NOK": true, "DKK": true, "ZAR": true,
		"MXN": true, "BRL": true, "RUB": true, "INR": true, "KRW": true,
		"TRY": true, "SAR": true, "AED": true, "PLN": true, "THB": true,
		"IDR": true, "MYR": true, "PHP": true, "CZK": true, "HUF": true,
		"CLP": true, "PEN": true, "COP": true, "ARS": true, "VND": true,

		// African currencies
		"NGN": true, "KES": true, "GHS": true, "EGP": true, "MAD": true,
		"TND": true, "ETB": true, "UGX": true, "TZS": true, "XOF": true,
		"XAF": true, "ZMW": true, "BWP": true, "MUR": true, "RWF": true,

		// Other currencies
		"ILS": true, "QAR": true, "OMR": true, "KWD": true, "BHD": true,
		"JOD": true, "LBP": true, "PKR": true, "LKR": true, "NPR": true,
		"BDT": true, "MMK": true, "KHR": true, "LAK": true, "ISK": true,
		"BGN": true, "HRK": true, "RON": true, "UAH": true, "BYN": true,
		"KZT": true, "UZS": true, "AMD": true, "GEL": true, "AZN": true,

		// Caribbean & Pacific
		"JMD": true, "TTD": true, "BBD": true, "BSD": true, "BZD": true,
		"XCD": true, "FJD": true, "PGK": true, "WST": true, "SBD": true,
		"TOP": true, "VUV": true,

		// Others
		"ALL": true, "AFN": true, "IRR": true, "IQD": true, "LYD": true,
		"SYP": true, "YER": true, "SOS": true, "SDG": true, "SSP": true,
		"MNT": true, "KGS": true, "TJS": true, "TMT": true, "MDL": true,
		"MKD": true, "BAM": true, "RSD": true, "DZD": true, "AOA": true,
		"NAD": true, "SZL": true, "LSL": true, "MGA": true, "MZN": true,
		"SCR": true, "MVR": true, "BTN": true, "BND": true, "TWD": true,
		"MOP": true, "KPW": true, "UYU": true, "PYG": true, "BOB": true,
		"VES": true, "SRD": true, "GYD": true, "HTG": true, "CUP": true,
		"DOP": true, "AWG": true, "ANG": true, "GTQ": true, "HNL": true,
		"NIO": true, "CRC": true, "PAB": true, "GMD": true, "GNF": true,
		"SLL": true, "LRD": true, "CVE": true, "STN": true, "BIF": true,
		"DJF": true, "ERN": true, "KMF": true, "CDF": true, "MRU": true,
		"GIP": true, "FKP": true, "SHP": true, "XPF": true,
	}
	return codes
}