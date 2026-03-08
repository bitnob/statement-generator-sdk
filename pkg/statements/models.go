package statements

import (
	"encoding/json"
	"time"

	"github.com/shopspring/decimal"
)

// TransactionType represents the type of transaction
type TransactionType string

const (
	// Credit represents money coming into the account
	Credit TransactionType = "credit"
	// Debit represents money going out of the account
	Debit TransactionType = "debit"
)

// Transaction represents a single transaction in the statement
type Transaction struct {
	ID          string                 `json:"id"`
	Date        time.Time              `json:"date"`
	Description string                 `json:"description"`
	Amount      decimal.Decimal        `json:"amount"`
	Type        TransactionType        `json:"type"`
	Balance     *decimal.Decimal       `json:"balance,omitempty"`
	Reference   string                 `json:"reference,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// Address represents a physical address for proof-of-address functionality
type Address struct {
	Line1      string `json:"line1"`
	Line2      string `json:"line2,omitempty"`
	City       string `json:"city"`
	State      string `json:"state,omitempty"`
	PostalCode string `json:"postal_code,omitempty"`
	Country    string `json:"country"`
}

// Account represents the account holder information
type Account struct {
	Number     string   `json:"number"`      // Can be masked (e.g., "****1234")
	HolderName string   `json:"holder_name"` // Account holder's name
	Currency   string   `json:"currency"`    // ISO 4217 currency code
	Type       string   `json:"type,omitempty"`    // Account type (Savings, Current, etc.)
	Address    *Address `json:"address,omitempty"` // Customer address for proof-of-address
}

// Institution represents the financial institution information
type Institution struct {
	Name         string   `json:"name"`
	Logo         []byte   `json:"logo,omitempty"`          // Logo image data (PNG/JPEG)
	LogoSVG      string   `json:"logo_svg,omitempty"`      // Logo in SVG format (preferred)
	Address      *Address `json:"address,omitempty"`       // Institution address
	RegNumber    string   `json:"reg_number,omitempty"`    // Registration number
	TaxID        string   `json:"tax_id,omitempty"`        // Tax identification number
	ContactPhone string   `json:"contact_phone,omitempty"` // Customer service phone
	ContactEmail string   `json:"contact_email,omitempty"` // Customer service email
	Website      string   `json:"website,omitempty"`       // Institution website
	FooterText   string   `json:"footer_text,omitempty"`   // Custom footer text
}

// StatementInput represents the complete input for statement generation
type StatementInput struct {
	Account        Account         `json:"account"`
	Transactions   []Transaction   `json:"transactions"`
	PeriodStart    time.Time       `json:"period_start"`
	PeriodEnd      time.Time       `json:"period_end"`
	OpeningBalance decimal.Decimal `json:"opening_balance"`
	Institution    *Institution    `json:"institution,omitempty"`
}

// Statement represents a generated statement with calculated values
type Statement struct {
	Input            StatementInput  `json:"input"`
	ClosingBalance   decimal.Decimal `json:"closing_balance"`
	TotalCredits     decimal.Decimal `json:"total_credits"`
	TotalDebits      decimal.Decimal `json:"total_debits"`
	TransactionCount int             `json:"transaction_count"`
	GeneratedAt      time.Time       `json:"generated_at"`

	// Internal fields (not exported to JSON)
	calculatedBalances []decimal.Decimal `json:"-"`
	generator          *StatementGenerator `json:"-"`
}

// CalculationResult holds the results of balance calculations
type CalculationResult struct {
	ClosingBalance decimal.Decimal
	TotalCredits   decimal.Decimal
	TotalDebits    decimal.Decimal
	Balances       []decimal.Decimal
}

// ValidationError represents a validation error with field information
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e ValidationError) Error() string {
	return e.Field + ": " + e.Message
}

// CalculationError represents an error in balance calculations
type CalculationError struct {
	Expected decimal.Decimal `json:"expected"`
	Actual   decimal.Decimal `json:"actual"`
	Message  string          `json:"message"`
}

func (e CalculationError) Error() string {
	return e.Message
}

// Config holds configuration for the statement generator
type Config struct {
	Institution     *Institution
	Locale          string
	HTMLTemplate    string
	TimeZone        *time.Location
	PDFRenderer     string // "minimalist" (default), "enhanced", or "simple"
	EnableColors    bool   // Enable colored text (default: false for black/white)
	AlternatingRows bool   // Enable alternating row colors (default: true)
}

// Option is a function that configures the StatementGenerator
type Option func(*Config)

// WithInstitution sets the default institution for statements
func WithInstitution(institution *Institution) Option {
	return func(c *Config) {
		c.Institution = institution
	}
}

// WithLocale sets the locale for formatting
func WithLocale(locale string) Option {
	return func(c *Config) {
		c.Locale = locale
	}
}

// WithHTMLTemplate sets a custom HTML template
func WithHTMLTemplate(template string) Option {
	return func(c *Config) {
		c.HTMLTemplate = template
	}
}

// WithTimeZone sets the timezone for date formatting
func WithTimeZone(tz *time.Location) Option {
	return func(c *Config) {
		c.TimeZone = tz
	}
}

// WithPDFRenderer sets the PDF renderer style
// Options: "minimalist" (default), "enhanced", "simple"
func WithPDFRenderer(renderer string) Option {
	return func(c *Config) {
		c.PDFRenderer = renderer
	}
}

// WithMinimalistDesign uses the minimalist black and white design (default)
func WithMinimalistDesign() Option {
	return func(c *Config) {
		c.PDFRenderer = "minimalist"
		c.EnableColors = false
		c.AlternatingRows = true
	}
}

// WithEnhancedDesign uses the colorful enhanced design
func WithEnhancedDesign() Option {
	return func(c *Config) {
		c.PDFRenderer = "enhanced"
		c.EnableColors = true
		c.AlternatingRows = true
	}
}

// WithSimpleDesign uses the basic simple design
func WithSimpleDesign() Option {
	return func(c *Config) {
		c.PDFRenderer = "simple"
		c.EnableColors = false
		c.AlternatingRows = false
	}
}

// MarshalJSON implements custom JSON marshaling for Transaction
func (t Transaction) MarshalJSON() ([]byte, error) {
	type Alias Transaction
	return json.Marshal(&struct {
		Date string `json:"date"`
		*Alias
	}{
		Date:  t.Date.Format(time.RFC3339),
		Alias: (*Alias)(&t),
	})
}

// UnmarshalJSON implements custom JSON unmarshaling for Transaction
func (t *Transaction) UnmarshalJSON(data []byte) error {
	type Alias Transaction
	aux := &struct {
		Date string `json:"date"`
		*Alias
	}{
		Alias: (*Alias)(t),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if aux.Date != "" {
		parsedDate, err := time.Parse(time.RFC3339, aux.Date)
		if err != nil {
			// Try alternative formats
			parsedDate, err = time.Parse("2006-01-02", aux.Date)
			if err != nil {
				return err
			}
		}
		t.Date = parsedDate
	}

	return nil
}