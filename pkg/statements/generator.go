package statements

import (
	"time"

	"github.com/shopspring/decimal"
)

// StatementGenerator is the main class for generating statements
type StatementGenerator struct {
	config     *Config
	validator  *Validator
	calculator *Calculator
}

// New creates a new StatementGenerator with options
func New(opts ...Option) *StatementGenerator {
	// Use default configuration with minimalist design
	config := DefaultConfig()

	// Apply options to override defaults
	for _, opt := range opts {
		opt(config)
	}

	return &StatementGenerator{
		config:     config,
		validator:  NewValidator(),
		calculator: NewCalculator(),
	}
}

// Generate generates a statement from input
func (g *StatementGenerator) Generate(input StatementInput) (*Statement, error) {
	// Apply institution from config if not provided
	if input.Institution == nil && g.config.Institution != nil {
		input.Institution = g.config.Institution
	}

	// Validate input
	if err := g.validator.ValidateStatementInput(input); err != nil {
		return nil, err
	}

	// Calculate totals
	result, err := g.calculator.CalculateStatementTotals(input)
	if err != nil {
		return nil, err
	}

	// Create statement
	statement := &Statement{
		Input:              input,
		ClosingBalance:     result.ClosingBalance,
		TotalCredits:       result.TotalCredits,
		TotalDebits:        result.TotalDebits,
		TransactionCount:   len(input.Transactions),
		GeneratedAt:        time.Now().In(g.config.TimeZone),
		calculatedBalances: result.Balances,
		generator:          g,
	}

	return statement, nil
}

// SetHTMLTemplate sets a custom HTML template for the generator
func (g *StatementGenerator) SetHTMLTemplate(template string) {
	g.config.HTMLTemplate = template
}

// ToPDF generates PDF from the statement
func (s *Statement) ToPDF() ([]byte, error) {
	// Implemented in statement_renderers.go to avoid circular dependency
	return s.renderPDF()
}

// ToCSV generates CSV from the statement
func (s *Statement) ToCSV() string {
	// Implemented in statement_renderers.go to avoid circular dependency
	return s.renderCSV()
}

// ToHTML generates HTML from the statement
func (s *Statement) ToHTML() string {
	// Implemented in statement_renderers.go to avoid circular dependency
	return s.renderHTML()
}

// QuickPDF generates a PDF statement quickly from transactions
func QuickPDF(transactions []Transaction, openingBalance decimal.Decimal) ([]byte, error) {
	// Create minimal input
	input := StatementInput{
		Account: Account{
			Number:     "****0000",
			HolderName: "Account Holder",
			Currency:   "USD",
		},
		Transactions:   transactions,
		PeriodStart:    findEarliestDate(transactions),
		PeriodEnd:      findLatestDate(transactions),
		OpeningBalance: openingBalance,
	}

	// Generate statement
	generator := New()
	statement, err := generator.Generate(input)
	if err != nil {
		return nil, err
	}

	return statement.ToPDF()
}

// QuickCSV generates a CSV statement quickly from transactions
func QuickCSV(transactions []Transaction, openingBalance decimal.Decimal) (string, error) {
	// Create minimal input
	input := StatementInput{
		Account: Account{
			Number:     "****0000",
			HolderName: "Account Holder",
			Currency:   "USD",
		},
		Transactions:   transactions,
		PeriodStart:    findEarliestDate(transactions),
		PeriodEnd:      findLatestDate(transactions),
		OpeningBalance: openingBalance,
	}

	// Generate statement
	generator := New()
	statement, err := generator.Generate(input)
	if err != nil {
		return "", err
	}

	return statement.ToCSV(), nil
}

// findEarliestDate finds the earliest date in transactions
func findEarliestDate(transactions []Transaction) time.Time {
	if len(transactions) == 0 {
		return time.Now().AddDate(0, -1, 0) // Default to 1 month ago
	}

	earliest := transactions[0].Date
	for _, txn := range transactions[1:] {
		if txn.Date.Before(earliest) {
			earliest = txn.Date
		}
	}
	return earliest
}

// findLatestDate finds the latest date in transactions
func findLatestDate(transactions []Transaction) time.Time {
	if len(transactions) == 0 {
		return time.Now() // Default to now
	}

	latest := transactions[0].Date
	for _, txn := range transactions[1:] {
		if txn.Date.After(latest) {
			latest = txn.Date
		}
	}
	return latest
}

// GetTransactionWithBalance returns a transaction with its calculated balance
func (s *Statement) GetTransactionWithBalance(index int) *Transaction {
	if index < 0 || index >= len(s.Input.Transactions) {
		return nil
	}

	txn := s.Input.Transactions[index]
	if index < len(s.calculatedBalances) {
		balance := s.calculatedBalances[index]
		txn.Balance = &balance
	}

	return &txn
}

// GetSortedTransactions returns transactions sorted by date
func (s *Statement) GetSortedTransactions() []Transaction {
	calculator := NewCalculator()
	return calculator.SortTransactions(s.Input.Transactions)
}

// GetTransactionsWithBalances returns all transactions with calculated balances
func (s *Statement) GetTransactionsWithBalances() []Transaction {
	calculator := NewCalculator()
	sorted := calculator.SortTransactions(s.Input.Transactions)
	return calculator.ApplyRunningBalances(sorted, s.Input.OpeningBalance)
}

// Validate validates a statement input without generating a statement
func Validate(input StatementInput) error {
	validator := NewValidator()
	return validator.ValidateStatementInput(input)
}

// CalculateTotals calculates totals without generating a full statement
func CalculateTotals(transactions []Transaction, openingBalance decimal.Decimal) (*CalculationResult, error) {
	calculator := NewCalculator()
	return calculator.CalculateBalances(transactions, openingBalance)
}

// Builder provides a fluent interface for building statements
type Builder struct {
	input StatementInput
	opts  []Option
}

// NewBuilder creates a new statement builder
func NewBuilder() *Builder {
	return &Builder{
		input: StatementInput{
			Transactions: []Transaction{},
		},
	}
}

// WithAccount sets the account information
func (b *Builder) WithAccount(number, holderName, currency string) *Builder {
	b.input.Account = Account{
		Number:     number,
		HolderName: holderName,
		Currency:   currency,
	}
	return b
}

// WithAccountAddress sets the account address
func (b *Builder) WithAccountAddress(address Address) *Builder {
	b.input.Account.Address = &address
	return b
}

// WithPeriod sets the statement period
func (b *Builder) WithPeriod(start, end time.Time) *Builder {
	b.input.PeriodStart = start
	b.input.PeriodEnd = end
	return b
}

// WithOpeningBalance sets the opening balance
func (b *Builder) WithOpeningBalance(balance decimal.Decimal) *Builder {
	b.input.OpeningBalance = balance
	return b
}

// AddTransaction adds a transaction
func (b *Builder) AddTransaction(txn Transaction) *Builder {
	b.input.Transactions = append(b.input.Transactions, txn)
	return b
}

// WithInstitution sets the institution
func (b *Builder) WithInstitution(institution Institution) *Builder {
	b.input.Institution = &institution
	return b
}

// WithOptions adds generator options
func (b *Builder) WithOptions(opts ...Option) *Builder {
	b.opts = append(b.opts, opts...)
	return b
}

// Build generates the statement
func (b *Builder) Build() (*Statement, error) {
	generator := New(b.opts...)
	return generator.Generate(b.input)
}