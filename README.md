# Statement Generator SDK

[![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.20-blue)](https://golang.org/doc/devel/release.html)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

A lightweight, embeddable Go SDK for generating professional bank and financial statements in multiple formats (PDF, CSV, HTML). Perfect for fintech applications requiring proof-of-address documentation and transaction reporting.

## Documentation

- [API Reference](docs/API.md) - Complete API documentation
- [Specification](docs/SPECIFICATION.md) - Detailed SDK specification
- [Contributing Guide](docs/CONTRIBUTING.md) - How to contribute
- [Changelog](docs/CHANGELOG.md) - Version history
- [Test Coverage](docs/TEST_COVERAGE.md) - Test suite documentation

## Features

- **Multi-Format Output**: Generate statements in PDF, CSV, and HTML formats
- **Proof of Address**: Full address support for account holders and institutions
- **International Support**: 200+ currency codes (ISO 4217) and locale-aware formatting
- **High Precision**: Decimal arithmetic for accurate financial calculations
- **Zero External Dependencies**: Core functionality requires no external services
- **Thread-Safe**: Concurrent generation support for high-throughput applications
- **Comprehensive Validation**: Built-in validation for all inputs
- **Flexible API**: Builder pattern and quick methods for rapid development

## Installation

```bash
go get github.com/bitnob/statement-generator-sdk
```

## Quick Start

```go
package main

import (
    "log"
    "time"
    "github.com/bitnob/statement-generator-sdk/pkg/statements"
    "github.com/shopspring/decimal"
)

func main() {
    // Create a simple statement with QuickPDF
    transactions := []statements.Transaction{
        {
            ID:          "TXN001",
            Date:        time.Now(),
            Description: "Salary Deposit",
            Amount:      decimal.NewFromFloat(5000),
            Type:        statements.Credit,
        },
        {
            ID:          "TXN002",
            Date:        time.Now(),
            Description: "Rent Payment",
            Amount:      decimal.NewFromFloat(-1200),
            Type:        statements.Debit,
        },
    }

    pdfBytes, err := statements.QuickPDF(
        transactions,
        decimal.NewFromFloat(1000), // Opening balance
    )
    if err != nil {
        log.Fatal(err)
    }

    // Save to file
    os.WriteFile("statement.pdf", pdfBytes, 0644)
}
```

## Full Example with Address (Proof of Address)

```go
// Create account with address for proof-of-address
account := statements.Account{
    Number:     "****1234",
    HolderName: "John Doe",
    Currency:   "USD",
    Type:       "Checking",
    Address: &statements.Address{
        Line1:      "123 Main Street",
        Line2:      "Apt 4B",
        City:       "New York",
        State:      "NY",
        PostalCode: "10001",
        Country:    "USA",
    },
}

// Create institution with contact details and custom footer
institution := &statements.Institution{
    Name: "Example Bank",
    Address: &statements.Address{
        Line1:   "100 Financial Plaza",
        City:    "New York",
        State:   "NY",
        Country: "USA",
    },
    RegNumber:    "Member FDIC | REG-123456",
    ContactPhone: "1-800-EXAMPLE (392-6753)",
    ContactEmail: "support@examplebank.com",
    Website:      "www.examplebank.com",
    FooterText:   "Thank you for banking with us. Please contact us with your account number for any questions.",
}

// Generate statement
generator := statements.New(
    statements.WithLocale("en-US"),
    statements.WithInstitution(institution),
)

input := statements.StatementInput{
    Account:        account,
    Transactions:   transactions,
    PeriodStart:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
    PeriodEnd:      time.Date(2024, 1, 31, 23, 59, 59, 0, time.UTC),
    OpeningBalance: decimal.NewFromFloat(500),
}

statement, err := generator.Generate(input)
if err != nil {
    log.Fatal(err)
}

// Export to multiple formats
pdfBytes, _ := statement.ToPDF()
csvString := statement.ToCSV()
htmlString := statement.ToHTML()
```

## Enhanced Footer and Contact Information

The SDK supports customizable footer content and contact information for professional statements:

```go
// Create institution with full contact details
institution := &statements.Institution{
    Name: "Premier Bank",
    Address: &statements.Address{
        Line1:   "500 Banking Center",
        City:    "New York",
        State:   "NY",
        Country: "USA",
    },
    RegNumber:    "Member FDIC | Equal Housing Lender",
    ContactPhone: "1-800-555-BANK (2265)",
    ContactEmail: "support@premierbank.com",
    Website:      "www.premierbank.com",
    FooterText:   "For questions about your statement, please contact our Customer Service team with your account number ready. We're available 24/7 to assist you.",
    LogoSVG:      logoSVGString, // Optional: SVG logo
    // Logo:      logoBytes,     // Alternative: PNG/JPEG bytes
}
```

### Footer Features

- **Customer Service Section**: Displays phone, email, and website prominently
- **Custom Footer Text**: Set your own message for customers
- **Account Reference**: Automatically includes "Please reference your account number (****1234) when contacting us"
- **Generation Timestamp**: Shows exact date and time statement was generated
- **Page Numbering**: "Page 1 of 3" format for multi-page statements
- **Logo Support**: Include your institution's logo in SVG or image format

### Default Behavior

If no custom footer text is provided, the SDK uses:
> "This is an official bank statement. Please review all transactions and report any discrepancies immediately."

The footer automatically includes:
- Generation timestamp at the bottom of every page
- Page numbers for multi-page statements
- Copyright notice with institution name

## API Documentation

### Core Types

#### Transaction
Represents a single financial transaction.

```go
type Transaction struct {
    ID          string          // Unique identifier
    Date        time.Time       // Transaction date
    Description string          // Transaction description
    Amount      decimal.Decimal // Transaction amount (positive or negative)
    Type        TransactionType // Credit or Debit
    Reference   string          // Optional reference number
    Metadata    interface{}     // Optional metadata
}
```

#### Account
Represents a bank account with optional address for proof-of-address.

```go
type Account struct {
    Number     string   // Account number (can be masked)
    HolderName string   // Account holder's name
    Currency   string   // ISO 4217 currency code
    Type       string   // Account type (e.g., Checking, Savings)
    Address    *Address // Optional address for proof-of-address
}
```

#### Address
Physical address information for proof-of-address documentation.

```go
type Address struct {
    Line1      string // Primary address line
    Line2      string // Secondary address line (optional)
    City       string // City
    State      string // State/Province (optional)
    PostalCode string // Postal/ZIP code (optional)
    Country    string // Country
}
```

### Statement Generator

#### Creating a Generator

```go
// Default generator
generator := statements.New()

// With options
generator := statements.New(
    statements.WithLocale("en-US"),              // Set locale for formatting
    statements.WithInstitution(institution),      // Set default institution
    statements.WithDateFormat("Jan 02, 2006"),   // Custom date format
    statements.WithCurrencySymbol("$"),          // Override currency symbol
)
```

#### Available Options

- `WithLocale(locale string)` - Set locale for formatting (en-US, en-GB, fr-FR, etc.)
- `WithInstitution(institution *Institution)` - Set default institution
- `WithDateFormat(format string)` - Override date format
- `WithCurrencySymbol(symbol string)` - Override currency symbol
- `WithDecimalPlaces(places int)` - Set decimal places for amounts

### Statement Output

#### Generated Statement Object

```go
type Statement struct {
    Account          Account
    Transactions     []Transaction
    PeriodStart      time.Time
    PeriodEnd        time.Time
    OpeningBalance   decimal.Decimal
    ClosingBalance   decimal.Decimal
    TotalCredits     decimal.Decimal
    TotalDebits      decimal.Decimal
    TransactionCount int
    Institution      *Institution
    GeneratedAt      time.Time
    Locale           string
}
```

#### Export Methods

```go
// Generate PDF (returns byte array)
pdfBytes, err := statement.ToPDF()

// Generate CSV (returns string)
csvString := statement.ToCSV()

// Generate HTML (returns string)
htmlString := statement.ToHTML()

// Generate JSON (returns byte array)
jsonBytes, err := statement.ToJSON()
```

### Builder Pattern

For complex statement creation:

```go
builder := statements.NewBuilder().
    SetAccount(account).
    SetPeriodStart(startDate).
    SetPeriodEnd(endDate).
    SetOpeningBalance(openingBalance).
    AddTransaction(transaction1).
    AddTransaction(transaction2).
    SetInstitution(institution)

statement, err := builder.Build()
```

### Quick Methods

For simple use cases:

```go
// Quick PDF generation
pdfBytes, err := statements.QuickPDF(transactions, openingBalance)

// Quick CSV generation
csvString, err := statements.QuickCSV(transactions, openingBalance)
```

## Validation

The SDK includes comprehensive validation:

- **Currency Codes**: Validates against 200+ ISO 4217 currency codes
- **Date Ranges**: Ensures transactions fall within statement period
- **Balances**: Validates balance calculations
- **Required Fields**: Ensures all required fields are present
- **Address Validation**: Validates address completeness for proof-of-address

```go
validator := statements.NewValidator()
err := validator.ValidateStatementInput(input)
if err != nil {
    // Handle validation error
}
```

## Formatting

### Currency Formatting

Supports multiple locales and currencies:

```go
formatter := statements.NewCurrencyFormatter("en-US")
formatted := formatter.FormatAmount(amount, "USD") // $1,234.56

formatter = statements.NewCurrencyFormatter("fr-FR")
formatted = formatter.FormatAmount(amount, "EUR") // 1 234,56 €
```

### Date Formatting

Locale-aware date formatting:

```go
formatter := statements.NewDateFormatter("en-US")
formatted := formatter.Format(date) // Jan 15, 2024

formatter = statements.NewDateFormatter("fr-FR")
formatted = formatter.Format(date) // 15 janv. 2024
```

## Performance

The SDK is optimized for performance:

- **1,000 transactions**: < 50ms generation time
- **10,000 transactions**: < 500ms generation time
- **Concurrent safe**: Thread-safe for parallel generation
- **Memory efficient**: Streaming for large datasets

### Benchmarks

```bash
# Run benchmarks
go test -bench=. ./test/...

# Results (M1 MacBook Pro)
BenchmarkStatementGeneration_100-8      2419    495337 ns/op
BenchmarkStatementGeneration_1000-8      247   4836421 ns/op
BenchmarkPDFGeneration_1000-8            134   8901234 ns/op
BenchmarkCSVGeneration_1000-8           3428    350123 ns/op
```

## Testing

```bash
# Run all tests
go test ./test/... -v

# Run specific test suite
go test ./test/... -run TestStatementGenerator

# Run with coverage
go test ./test/... -cover

# Run benchmarks
go test -bench=. ./test/...
```

## Error Handling

The SDK provides detailed error messages:

```go
statement, err := generator.Generate(input)
if err != nil {
    switch e := err.(type) {
    case *statements.ValidationError:
        // Handle validation error
        log.Printf("Validation failed: %s", e.Field)
    case *statements.FormatError:
        // Handle formatting error
        log.Printf("Formatting failed: %s", e.Message)
    default:
        // Handle other errors
        log.Printf("Generation failed: %v", err)
    }
}
```

## Examples

See the [examples](examples/) directory for working examples:

- [Basic Example](examples/basic/) - Complete statement generation with all export formats
- [Simple Example](examples/simple/) - Minimal configuration example
- [Footer Demo](examples/footer_demo/) - Various footer and contact information styles
- [Template Demo](examples/template_demo/) - Different PDF template styles (minimalist, enhanced, simple)

## Requirements

- Go 1.20 or higher
- `github.com/jung-kurt/gofpdf/v2` - PDF generation
- `github.com/shopspring/decimal` - High-precision decimal arithmetic
- `github.com/stretchr/testify` - Testing framework (dev only)

## Contributing

Contributions are welcome! Please read our [Contributing Guide](docs/CONTRIBUTING.md) for details.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

For issues, questions, or contributions, please visit our [GitHub repository](https://github.com/bitnob/statement-generator-sdk).

## Acknowledgments

- Built with love by the Statement Generator SDK Contributors
- Special thanks to all contributors
- Powered by Go's excellent standard library