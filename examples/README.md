# Statement Generator SDK Examples

This directory contains example applications demonstrating various features of the Statement Generator SDK.

## Available Examples

### 1. Simple Example (`simple_example.go`)
Basic statement generation with minimal configuration:
- Creates a checking account with address
- Adds sample transactions
- Generates PDF with default minimalist template
- Shows institution footer customization

```bash
go run simple_example.go
```

### 2. Template Styles (`test_templates.go`)
Demonstrates different PDF template styles:
- Minimalist (default) - Black and white with alternating rows
- Enhanced - Colored headers and amounts
- Simple - Basic text layout
- Custom - With locale settings

```bash
go run test_templates.go
```

### 3. Footer and Contact Information (`test_footer_logo.go`)
Shows various footer configurations:
- Standard bank with basic contact
- Premium bank with full details
- Fintech with digital-first approach
- International bank with compliance text

```bash
go run test_footer_logo.go
```

### 4. Basic Example (`basic/main.go`)
Complete example with all export formats:
- Generates PDF, CSV, and HTML
- Shows balance calculations
- Demonstrates proof-of-address feature

```bash
cd basic && go run main.go
```

## Features Demonstrated

- **Proof of Address**: Account holder addresses included in statements
- **Multi-format Export**: PDF, CSV, HTML generation
- **Custom Footers**: Institution-specific messaging
- **Contact Information**: Phone, email, website display
- **Logo Support**: SVG and image logo integration
- **Locale Support**: International formatting
- **Template Styles**: Multiple design options
- **Balance Calculations**: Automatic running balance

## Running the Examples

1. Install dependencies:
```bash
go mod download
```

2. Run any example:
```bash
go run [example_file].go
```

3. Check generated files:
- PDF files: `*.pdf`
- CSV files: `*.csv`
- HTML files: `*.html`

## Customization

Each example can be modified to test different features:
- Change transaction data
- Modify institution details
- Test different locales
- Experiment with footer text
- Try various template styles

## Documentation

For full API documentation, see:
- [API Reference](../API.md)
- [Main README](../README.md)
- [Contributing Guide](../CONTRIBUTING.md)