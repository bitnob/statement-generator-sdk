# Statement Generator SDK - API Reference

## Table of Contents
- [Core Package](#core-package)
- [Types](#types)
- [Functions](#functions)
- [Methods](#methods)
- [Options](#options)
- [Validators](#validators)
- [Formatters](#formatters)
- [Errors](#errors)

## Core Package

```go
import "github.com/bitnob/statement-generator-sdk/pkg/statements"
```

## Types

### Transaction

Represents a single financial transaction.

```go
type Transaction struct {
    ID          string          `json:"id"`
    Date        time.Time       `json:"date"`
    Description string          `json:"description"`
    Amount      decimal.Decimal `json:"amount"`
    Type        TransactionType `json:"type"`
    Reference   string          `json:"reference,omitempty"`
    Metadata    interface{}     `json:"metadata,omitempty"`
}
```

#### Fields
- `ID` (string, required): Unique identifier for the transaction
- `Date` (time.Time, required): Transaction date and time
- `Description` (string, required): Human-readable description
- `Amount` (decimal.Decimal, required): Transaction amount (positive for credits, negative for debits)
- `Type` (TransactionType, required): Either `Credit` or `Debit`
- `Reference` (string, optional): External reference number
- `Metadata` (interface{}, optional): Additional data for custom use

### TransactionType

```go
type TransactionType string

const (
    Credit TransactionType = "credit"
    Debit  TransactionType = "debit"
)
```

### Account

Represents a financial account with optional address for proof-of-address.

```go
type Account struct {
    Number     string   `json:"number"`
    HolderName string   `json:"holder_name"`
    Currency   string   `json:"currency"`
    Type       string   `json:"type,omitempty"`
    Address    *Address `json:"address,omitempty"`
}
```

#### Fields
- `Number` (string, required): Account number (can be masked, e.g., "****1234")
- `HolderName` (string, required): Full name of account holder
- `Currency` (string, required): ISO 4217 currency code (e.g., "USD", "EUR", "NGN")
- `Type` (string, optional): Account type (e.g., "Checking", "Savings", "Investment")
- `Address` (*Address, optional): Physical address for proof-of-address

### Address

Physical address information for proof-of-address documentation.

```go
type Address struct {
    Line1      string `json:"line1"`
    Line2      string `json:"line2,omitempty"`
    City       string `json:"city"`
    State      string `json:"state,omitempty"`
    PostalCode string `json:"postal_code,omitempty"`
    Country    string `json:"country"`
}
```

#### Fields
- `Line1` (string, required): Primary address line
- `Line2` (string, optional): Secondary address line (apartment, suite, etc.)
- `City` (string, required): City name
- `State` (string, optional): State or province
- `PostalCode` (string, optional): Postal or ZIP code
- `Country` (string, required): Country name

### Institution

Represents the financial institution issuing the statement.

```go
type Institution struct {
    Name         string   `json:"name"`
    Logo         []byte   `json:"logo,omitempty"`
    LogoSVG      string   `json:"logo_svg,omitempty"`
    Address      *Address `json:"address,omitempty"`
    RegNumber    string   `json:"reg_number,omitempty"`
    TaxID        string   `json:"tax_id,omitempty"`
    ContactPhone string   `json:"contact_phone,omitempty"`
    ContactEmail string   `json:"contact_email,omitempty"`
    Website      string   `json:"website,omitempty"`
    FooterText   string   `json:"footer_text,omitempty"`
}
```

#### Fields
- `Name` (string, required): Institution name
- `Logo` ([]byte, optional): Logo image data in PNG or JPEG format
- `LogoSVG` (string, optional): Logo in SVG format (preferred for scalability)
- `Address` (*Address, optional): Institution's physical address
- `RegNumber` (string, optional): Registration or license number (e.g., "Member FDIC")
- `TaxID` (string, optional): Tax identification number
- `ContactPhone` (string, optional): Customer service phone number
- `ContactEmail` (string, optional): Customer support email address
- `Website` (string, optional): Institution website URL
- `FooterText` (string, optional): Custom footer message for statements

### StatementInput

Input data for generating a statement.

```go
type StatementInput struct {
    Account        Account         `json:"account"`
    Transactions   []Transaction   `json:"transactions"`
    PeriodStart    time.Time       `json:"period_start"`
    PeriodEnd      time.Time       `json:"period_end"`
    OpeningBalance decimal.Decimal `json:"opening_balance"`
    Institution    *Institution    `json:"institution,omitempty"`
}
```

#### Fields
- `Account` (Account, required): Account information
- `Transactions` ([]Transaction, required): List of transactions
- `PeriodStart` (time.Time, required): Statement period start
- `PeriodEnd` (time.Time, required): Statement period end
- `OpeningBalance` (decimal.Decimal, required): Balance at period start
- `Institution` (*Institution, optional): Issuing institution

### Statement

Generated statement with calculated balances and summaries.

```go
type Statement struct {
    Account          Account         `json:"account"`
    Transactions     []Transaction   `json:"transactions"`
    PeriodStart      time.Time       `json:"period_start"`
    PeriodEnd        time.Time       `json:"period_end"`
    OpeningBalance   decimal.Decimal `json:"opening_balance"`
    ClosingBalance   decimal.Decimal `json:"closing_balance"`
    TotalCredits     decimal.Decimal `json:"total_credits"`
    TotalDebits      decimal.Decimal `json:"total_debits"`
    TransactionCount int            `json:"transaction_count"`
    Institution      *Institution    `json:"institution,omitempty"`
    GeneratedAt      time.Time       `json:"generated_at"`
    Locale           string         `json:"locale"`
}
```

## Functions

### New

Creates a new StatementGenerator with optional configuration.

```go
func New(options ...GeneratorOption) *StatementGenerator
```

#### Parameters
- `options` (variadic GeneratorOption): Configuration options

#### Returns
- `*StatementGenerator`: Configured generator instance

#### Example
```go
generator := statements.New(
    statements.WithLocale("en-US"),
    statements.WithInstitution(institution),
)
```

### QuickPDF

Generates a PDF statement with minimal configuration.

```go
func QuickPDF(transactions []Transaction, openingBalance decimal.Decimal) ([]byte, error)
```

#### Parameters
- `transactions` ([]Transaction): List of transactions
- `openingBalance` (decimal.Decimal): Opening balance

#### Returns
- `[]byte`: PDF document bytes
- `error`: Error if generation fails

### QuickCSV

Generates a CSV statement with minimal configuration.

```go
func QuickCSV(transactions []Transaction, openingBalance decimal.Decimal) (string, error)
```

#### Parameters
- `transactions` ([]Transaction): List of transactions
- `openingBalance` (decimal.Decimal): Opening balance

#### Returns
- `string`: CSV content
- `error`: Error if generation fails

### NewBuilder

Creates a new statement builder for fluent API usage.

```go
func NewBuilder() *StatementBuilder
```

#### Returns
- `*StatementBuilder`: New builder instance

### NewValidator

Creates a new validator instance.

```go
func NewValidator() *Validator
```

#### Returns
- `*Validator`: New validator instance

### NewCalculator

Creates a new balance calculator.

```go
func NewCalculator() *Calculator
```

#### Returns
- `*Calculator`: New calculator instance

### NewCurrencyFormatter

Creates a locale-aware currency formatter.

```go
func NewCurrencyFormatter(locale string) *CurrencyFormatter
```

#### Parameters
- `locale` (string): Locale code (e.g., "en-US", "fr-FR")

#### Returns
- `*CurrencyFormatter`: Configured formatter

### NewDateFormatter

Creates a locale-aware date formatter.

```go
func NewDateFormatter(locale string) *DateFormatter
```

#### Parameters
- `locale` (string): Locale code

#### Returns
- `*DateFormatter`: Configured formatter

## Methods

### StatementGenerator Methods

#### Generate

Generates a statement from input data.

```go
func (g *StatementGenerator) Generate(input StatementInput) (*Statement, error)
```

#### Parameters
- `input` (StatementInput): Statement input data

#### Returns
- `*Statement`: Generated statement
- `error`: Validation or generation error

### Statement Methods

#### ToPDF

Exports statement as PDF.

```go
func (s *Statement) ToPDF() ([]byte, error)
```

#### Returns
- `[]byte`: PDF document bytes
- `error`: Export error

#### ToCSV

Exports statement as CSV.

```go
func (s *Statement) ToCSV() string
```

#### Returns
- `string`: CSV content

#### ToHTML

Exports statement as HTML.

```go
func (s *Statement) ToHTML() string
```

#### Returns
- `string`: HTML content

#### ToJSON

Exports statement as JSON.

```go
func (s *Statement) ToJSON() ([]byte, error)
```

#### Returns
- `[]byte`: JSON bytes
- `error`: Marshaling error

### StatementBuilder Methods

#### SetAccount

```go
func (b *StatementBuilder) SetAccount(account Account) *StatementBuilder
```

#### SetPeriodStart

```go
func (b *StatementBuilder) SetPeriodStart(start time.Time) *StatementBuilder
```

#### SetPeriodEnd

```go
func (b *StatementBuilder) SetPeriodEnd(end time.Time) *StatementBuilder
```

#### SetOpeningBalance

```go
func (b *StatementBuilder) SetOpeningBalance(balance decimal.Decimal) *StatementBuilder
```

#### AddTransaction

```go
func (b *StatementBuilder) AddTransaction(transaction Transaction) *StatementBuilder
```

#### SetInstitution

```go
func (b *StatementBuilder) SetInstitution(institution *Institution) *StatementBuilder
```

#### Build

```go
func (b *StatementBuilder) Build() (*Statement, error)
```

### Validator Methods

#### ValidateStatementInput

Validates statement input data.

```go
func (v *Validator) ValidateStatementInput(input StatementInput) error
```

#### ValidateTransaction

Validates a single transaction.

```go
func (v *Validator) ValidateTransaction(tx Transaction) error
```

#### ValidateAccount

Validates account information.

```go
func (v *Validator) ValidateAccount(account Account) error
```

#### ValidateCurrency

Validates ISO 4217 currency code.

```go
func (v *Validator) ValidateCurrency(currency string) error
```

### Calculator Methods

#### CalculateBalances

Calculates running balances and totals.

```go
func (c *Calculator) CalculateBalances(
    transactions []Transaction,
    openingBalance decimal.Decimal,
) (*BalanceResult, error)
```

#### Parameters
- `transactions` ([]Transaction): Sorted transactions
- `openingBalance` (decimal.Decimal): Opening balance

#### Returns
- `*BalanceResult`: Calculated balances and totals
- `error`: Calculation error

### CurrencyFormatter Methods

#### FormatAmount

Formats amount with currency symbol and separators.

```go
func (f *CurrencyFormatter) FormatAmount(
    amount decimal.Decimal,
    currency string,
) string
```

#### FormatBalance

Formats balance with sign indicator.

```go
func (f *CurrencyFormatter) FormatBalance(
    balance decimal.Decimal,
    currency string,
) string
```

### DateFormatter Methods

#### Format

Formats date according to locale.

```go
func (f *DateFormatter) Format(date time.Time) string
```

#### FormatRange

Formats date range.

```go
func (f *DateFormatter) FormatRange(start, end time.Time) string
```

## Options

### Generator Options

Options for configuring StatementGenerator.

#### WithLocale

Sets the locale for formatting.

```go
func WithLocale(locale string) GeneratorOption
```

#### WithInstitution

Sets default institution.

```go
func WithInstitution(institution *Institution) GeneratorOption
```

#### WithDateFormat

Overrides default date format.

```go
func WithDateFormat(format string) GeneratorOption
```

#### WithCurrencySymbol

Overrides currency symbol.

```go
func WithCurrencySymbol(symbol string) GeneratorOption
```

#### WithDecimalPlaces

Sets decimal places for amounts.

```go
func WithDecimalPlaces(places int) GeneratorOption
```

## Errors

### ValidationError

Returned when validation fails.

```go
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string
```

### FormatError

Returned when formatting fails.

```go
type FormatError struct {
    Type    string
    Message string
}

func (e *FormatError) Error() string
```

### GenerationError

Returned when statement generation fails.

```go
type GenerationError struct {
    Step    string
    Message string
    Cause   error
}

func (e *GenerationError) Error() string
```

## Constants

### Supported Currencies

The SDK supports all ISO 4217 currency codes. Common ones include:

```go
const (
    USD = "USD" // US Dollar
    EUR = "EUR" // Euro
    GBP = "GBP" // British Pound
    NGN = "NGN" // Nigerian Naira
    JPY = "JPY" // Japanese Yen
    CNY = "CNY" // Chinese Yuan
    // ... 200+ more
)
```

### Supported Locales

```go
const (
    LocaleEnUS = "en-US" // English (United States)
    LocaleEnGB = "en-GB" // English (United Kingdom)
    LocaleFrFR = "fr-FR" // French (France)
    LocaleDeDE = "de-DE" // German (Germany)
    LocaleEsES = "es-ES" // Spanish (Spain)
    LocalePtBR = "pt-BR" // Portuguese (Brazil)
    LocaleZhCN = "zh-CN" // Chinese (Simplified)
    LocaleJaJP = "ja-JP" // Japanese (Japan)
)
```

## Thread Safety

All public methods in the SDK are thread-safe and can be called concurrently. The StatementGenerator can be shared across goroutines.

```go
generator := statements.New()

// Safe to use concurrently
go func() {
    statement, _ := generator.Generate(input1)
}()

go func() {
    statement, _ := generator.Generate(input2)
}()
```

## Performance Considerations

- For large datasets (>10,000 transactions), consider streaming or pagination
- PDF generation is the most resource-intensive operation
- CSV generation is the fastest export format
- Use builder pattern for complex statements to avoid intermediate allocations

## Usage Examples

### Example 1: Basic Statement with Contact Information

```go
institution := &statements.Institution{
    Name:         "Community Bank",
    ContactPhone: "1-800-555-1234",
    ContactEmail: "help@communitybank.com",
}

generator := statements.New(
    statements.WithInstitution(institution),
)

// Statement will include contact info in footer
```

### Example 2: Full Institution Configuration

```go
institution := &statements.Institution{
    Name: "International Bank Corp",
    Address: &statements.Address{
        Line1:   "100 Wall Street",
        City:    "New York",
        State:   "NY",
        Country: "USA",
    },
    RegNumber:    "FDIC #12345 | SWIFT: IBCUUS33",
    ContactPhone: "+1-212-555-0100",
    ContactEmail: "support@intlbank.com",
    Website:      "www.intlbank.com",
    FooterText:   "Thank you for choosing International Bank. For immediate assistance with your account, please call our 24/7 hotline or visit our website. Always have your account number ready when contacting us.",
    LogoSVG:      svgLogoData, // Your SVG logo as string
}

generator := statements.New(
    statements.WithInstitution(institution),
    statements.WithLocale("en-US"),
)
```

### Example 3: Fintech/Digital Bank Configuration

```go
institution := &statements.Institution{
    Name:         "NeoBank",
    ContactEmail: "support@neobank.app",
    Website:      "app.neobank.io",
    FooterText:   "Questions? Open the NeoBank app and tap 'Support' for instant help, or email us with your account number.",
}

// Minimal contact info for digital-first bank
```

### Example 4: Multi-language Support Footer

```go
institution := &statements.Institution{
    Name:         "Global Bank",
    ContactPhone: "1-800-GLOBAL1",
    FooterText:   "For assistance in English, press 1. Para español, oprima 2. Pour le français, appuyez sur 3.",
}
```

## Statement Output Features

### Footer Information Display

The generated statements will automatically include:

1. **Contact Section** (if provided):
   ```
   Customer Service
   Phone: 1-800-555-1234
   Email: support@bank.com
   Website: www.bank.com
   ```

2. **Custom Footer Text** with account reference:
   ```
   [Your custom message]. Please reference your account
   number (****1234) when contacting us.
   ```

3. **Page Footer** (on every page):
   ```
   Page 1 of 3        Generated: Mar 8, 2026 at 2:30 PM EST
   © Your Bank Name. All rights reserved.
   ```

### Logo Placement

Logos appear at the top of the statement:
- SVG logos are preferred for scalability
- PNG/JPEG logos are supported via byte array
- Maximum recommended height: 15mm (PDF)

## Version History

- v1.0.0 - Initial release with PDF, CSV, HTML support
- v1.1.0 - Added address support for proof-of-address
- v1.2.0 - Added multi-currency support with 200+ ISO codes
- v1.3.0 - Added customizable footer, contact info, and logo support