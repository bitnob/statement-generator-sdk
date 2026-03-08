package statements

import (
	"fmt"
	"sort"

	"github.com/shopspring/decimal"
)

// Calculator handles balance and total calculations
type Calculator struct{}

// NewCalculator creates a new calculator instance
func NewCalculator() *Calculator {
	return &Calculator{}
}

// CalculateBalances calculates running balances and totals for transactions
func (c *Calculator) CalculateBalances(transactions []Transaction, openingBalance decimal.Decimal) (*CalculationResult, error) {
	result := &CalculationResult{
		ClosingBalance: openingBalance,
		TotalCredits:   decimal.Zero,
		TotalDebits:    decimal.Zero,
		Balances:       make([]decimal.Decimal, len(transactions)),
	}

	runningBalance := openingBalance

	for i, txn := range transactions {
		// Add transaction amount to running balance
		runningBalance = runningBalance.Add(txn.Amount)
		result.Balances[i] = runningBalance

		// Calculate totals based on transaction type
		if txn.Type == Credit {
			// Credits should have positive amounts
			if txn.Amount.IsPositive() {
				result.TotalCredits = result.TotalCredits.Add(txn.Amount)
			} else if txn.Amount.IsZero() {
				// Zero amount is acceptable, just skip
			} else {
				// Negative credit amount - treat absolute value as credit
				result.TotalCredits = result.TotalCredits.Add(txn.Amount.Abs())
			}
		} else if txn.Type == Debit {
			// Debits should have negative amounts
			if txn.Amount.IsNegative() {
				result.TotalDebits = result.TotalDebits.Add(txn.Amount.Abs())
			} else if txn.Amount.IsZero() {
				// Zero amount is acceptable, just skip
			} else {
				// Positive debit amount - treat as debit
				result.TotalDebits = result.TotalDebits.Add(txn.Amount)
			}
		}
	}

	result.ClosingBalance = runningBalance

	return result, nil
}

// ValidateProvidedBalances validates that provided balances match calculations
func (c *Calculator) ValidateProvidedBalances(transactions []Transaction, openingBalance decimal.Decimal) error {
	runningBalance := openingBalance

	for i, txn := range transactions {
		runningBalance = runningBalance.Add(txn.Amount)

		if txn.Balance != nil {
			// Allow for small rounding differences (0.0000001)
			tolerance := decimal.NewFromFloat(0.0000001)
			diff := runningBalance.Sub(*txn.Balance).Abs()

			if diff.GreaterThan(tolerance) {
				return &CalculationError{
					Expected: runningBalance,
					Actual:   *txn.Balance,
					Message: fmt.Sprintf(
						"balance mismatch for transaction %s at index %d: expected %s, got %s",
						txn.ID, i, runningBalance.StringFixed(2), txn.Balance.StringFixed(2),
					),
				}
			}
		}
	}

	return nil
}

// SortTransactions sorts transactions by date (and preserves original order for same date)
func (c *Calculator) SortTransactions(transactions []Transaction) []Transaction {
	// Create a copy to avoid modifying the original slice
	sorted := make([]Transaction, len(transactions))
	copy(sorted, transactions)

	// Stable sort to preserve order of same-date transactions
	sort.SliceStable(sorted, func(i, j int) bool {
		return sorted[i].Date.Before(sorted[j].Date)
	})

	return sorted
}

// CalculateStatementTotals calculates all totals for a statement input
func (c *Calculator) CalculateStatementTotals(input StatementInput) (*CalculationResult, error) {
	// Sort transactions by date
	sortedTransactions := c.SortTransactions(input.Transactions)

	// Calculate balances and totals
	result, err := c.CalculateBalances(sortedTransactions, input.OpeningBalance)
	if err != nil {
		return nil, err
	}

	// Validate provided balances if any
	hasProvidedBalances := false
	for _, txn := range sortedTransactions {
		if txn.Balance != nil {
			hasProvidedBalances = true
			break
		}
	}

	if hasProvidedBalances {
		if err := c.ValidateProvidedBalances(sortedTransactions, input.OpeningBalance); err != nil {
			return nil, err
		}
	}

	return result, nil
}

// ReconcileBalance checks if the calculated closing balance matches an expected value
func (c *Calculator) ReconcileBalance(calculated, expected decimal.Decimal) error {
	tolerance := decimal.NewFromFloat(0.01) // Allow 1 cent difference
	diff := calculated.Sub(expected).Abs()

	if diff.GreaterThan(tolerance) {
		return &CalculationError{
			Expected: expected,
			Actual:   calculated,
			Message: fmt.Sprintf(
				"balance reconciliation failed: expected %s, calculated %s (difference: %s)",
				expected.StringFixed(2), calculated.StringFixed(2), diff.StringFixed(2),
			),
		}
	}

	return nil
}

// ApplyRunningBalances updates transactions with calculated running balances
func (c *Calculator) ApplyRunningBalances(transactions []Transaction, openingBalance decimal.Decimal) []Transaction {
	result := make([]Transaction, len(transactions))
	runningBalance := openingBalance

	for i, txn := range transactions {
		runningBalance = runningBalance.Add(txn.Amount)
		result[i] = txn

		// Only update balance if not already provided
		if txn.Balance == nil {
			balance := runningBalance
			result[i].Balance = &balance
		}
	}

	return result
}

// CalculatePeriodSummary calculates summary statistics for a period
func (c *Calculator) CalculatePeriodSummary(transactions []Transaction) map[string]interface{} {
	summary := map[string]interface{}{
		"transaction_count":     len(transactions),
		"average_credit":        decimal.Zero,
		"average_debit":         decimal.Zero,
		"largest_credit":        decimal.Zero,
		"largest_debit":         decimal.Zero,
		"credit_count":          0,
		"debit_count":           0,
	}

	creditCount := 0
	debitCount := 0
	creditSum := decimal.Zero
	debitSum := decimal.Zero
	largestCredit := decimal.Zero
	largestDebit := decimal.Zero

	for _, txn := range transactions {
		if txn.Type == Credit && txn.Amount.IsPositive() {
			creditCount++
			creditSum = creditSum.Add(txn.Amount)
			if txn.Amount.GreaterThan(largestCredit) {
				largestCredit = txn.Amount
			}
		} else if txn.Type == Debit {
			debitCount++
			absAmount := txn.Amount.Abs()
			debitSum = debitSum.Add(absAmount)
			if absAmount.GreaterThan(largestDebit) {
				largestDebit = absAmount
			}
		}
	}

	summary["credit_count"] = creditCount
	summary["debit_count"] = debitCount
	summary["largest_credit"] = largestCredit
	summary["largest_debit"] = largestDebit

	if creditCount > 0 {
		avgCredit := creditSum.Div(decimal.NewFromInt(int64(creditCount)))
		summary["average_credit"] = avgCredit
	}

	if debitCount > 0 {
		avgDebit := debitSum.Div(decimal.NewFromInt(int64(debitCount)))
		summary["average_debit"] = avgDebit
	}

	return summary
}