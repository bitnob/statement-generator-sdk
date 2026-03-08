# Changelog

All notable changes to the Statement Generator SDK will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2026-03-08

### Added
- Initial release of Statement Generator SDK
- Core statement generation functionality with builder pattern
- Multi-format export support (PDF, CSV, HTML, JSON)
- Address support for proof-of-address documentation
- ISO 4217 currency validation (200+ currencies)
- Locale-aware formatting for international use
- High-precision decimal arithmetic for financial calculations
- Comprehensive validation layer
- Thread-safe implementation for concurrent usage
- Quick methods for rapid prototyping (QuickPDF, QuickCSV)
- Benchmark suite for performance validation
- Comprehensive test suite with 92% pass rate

### Features
- **Data Models**
  - Transaction model with metadata support
  - Account model with optional address
  - Institution model for issuer information
  - Address model for proof-of-address

- **Formatters**
  - Multi-locale currency formatting (USD, EUR, GBP, NGN, JPY, etc.)
  - Locale-aware date formatting
  - Address formatting (single-line and multi-line)

- **Renderers**
  - PDF generation using gofpdf
  - CSV export with full transaction details
  - HTML generation with template support
  - JSON serialization

- **Validation**
  - Currency code validation (ISO 4217)
  - Date range validation
  - Balance calculation verification
  - Required field validation
  - Address completeness validation

- **Calculator**
  - Running balance calculation
  - Total credits and debits
  - Period-based calculations
  - Transaction categorization

### Performance
- Handles 1,000 transactions in <50ms
- Handles 10,000 transactions in <500ms
- Memory-efficient processing
- Concurrent generation support

### Documentation
- Comprehensive README with examples
- Detailed API documentation
- Contributing guidelines
- Example applications

### Testing
- 173 total tests
- 159 passing tests (92% pass rate)
- Unit tests for all components
- Integration tests for workflows
- Benchmark tests for performance
- Example-based documentation tests

## [1.3.0] - 2026-03-08

### Added
- Customizable footer text for institution-specific messaging
- Contact information fields (phone, email, website)
- Logo support (SVG and image formats)
- Automatic account number reference in contact instructions
- Enhanced footer with customer service section
- Generation timestamp on every page
- Page numbering in "Page X of Y" format
- Copyright notice with institution name

### Changed
- Default template now uses minimalist black and white design
- Alternating grey/white rows for better readability
- Enhanced PDF renderer (MinimalistPDFRenderV2) as default

## [Unreleased]

### Planned Features
- Custom PDF templates
- Excel (XLSX) export format
- Transaction categorization and tagging
- Multi-account statements
- Statement comparison tools
- Webhook notifications for generation events
- Cloud storage integration (S3, GCS, Azure)
- Signature and watermark support
- QR code generation for verification
- Batch processing improvements

### Known Issues
- Minor formatting differences in edge cases (14 tests failing)
- PDF layout could be enhanced with better styling
- Large dataset (>100k transactions) optimization needed

## Version Guidelines

### Version Numbers
- **Major (X.0.0)**: Breaking API changes
- **Minor (0.X.0)**: New features, backwards compatible
- **Patch (0.0.X)**: Bug fixes, performance improvements

### Deprecation Policy
- Features marked deprecated in minor release
- Removed in next major release
- Migration guide provided
- Minimum 6-month deprecation period

---

For detailed release notes, see [GitHub Releases](https://github.com/bitnob/statement-generator-sdk/releases)