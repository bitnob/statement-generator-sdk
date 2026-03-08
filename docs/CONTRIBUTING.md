# Contributing to Statement Generator SDK

Thank you for considering contributing to the Statement Generator SDK! We welcome contributions from the community and are grateful for any help you can provide.

## Table of Contents
- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [How to Contribute](#how-to-contribute)
- [Pull Request Process](#pull-request-process)
- [Coding Standards](#coding-standards)
- [Testing Guidelines](#testing-guidelines)
- [Documentation](#documentation)
- [Community](#community)

## Code of Conduct

By participating in this project, you agree to abide by our Code of Conduct:

- Be respectful and inclusive
- Welcome newcomers and help them get started
- Focus on constructive criticism
- Accept feedback gracefully
- Prioritize the community's best interests

## Getting Started

1. **Fork the Repository**: Click the "Fork" button on GitHub
2. **Clone Your Fork**:
   ```bash
   git clone https://github.com/YOUR_USERNAME/sdk.git
   cd sdk
   ```
3. **Add Upstream Remote**:
   ```bash
   git remote add upstream https://github.com/statement-generator/sdk.git
   ```

## Development Setup

### Prerequisites

- Go 1.20 or higher
- Git
- Make (optional but recommended)

### Installation

1. **Install Dependencies**:
   ```bash
   go mod download
   ```

2. **Run Tests**:
   ```bash
   go test ./test/... -v
   ```

3. **Run Benchmarks**:
   ```bash
   go test -bench=. ./test/...
   ```

### Project Structure

```
statements/
├── pkg/
│   └── statements/       # Main SDK package
│       ├── models.go      # Core data models
│       ├── validator.go   # Validation logic
│       ├── calculator.go  # Balance calculations
│       ├── formatters.go  # Currency/date formatting
│       ├── generator.go   # Main generator logic
│       └── *.go          # Other components
├── test/                  # Test files
├── examples/             # Example applications
├── fixtures/             # Test fixtures
└── docs/                 # Documentation
```

## How to Contribute

### Reporting Bugs

Before creating bug reports, please check existing issues. When creating a bug report, include:

1. **Clear Title**: Descriptive problem summary
2. **Description**: What happened vs. what you expected
3. **Steps to Reproduce**: Minimal code example
4. **Environment**: Go version, OS, etc.
5. **Stack Trace**: If applicable

Example:
```markdown
### Bug: CSV export missing transaction references

**Expected**: CSV should include reference column
**Actual**: Reference column is empty

**Code to reproduce**:
```go
// Your code here
```

**Environment**: Go 1.21, macOS 14.0
```

### Suggesting Features

Feature requests are welcome! Please provide:

1. **Use Case**: Why is this feature needed?
2. **Proposed Solution**: How should it work?
3. **Alternatives**: Other approaches considered
4. **Examples**: Code examples of the proposed API

### Submitting Code

1. **Create a Branch**:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make Changes**: Write clean, tested code

3. **Test Your Changes**:
   ```bash
   go test ./test/... -v
   ```

4. **Commit**:
   ```bash
   git commit -m "feat: add new feature description"
   ```

## Pull Request Process

### Before Submitting

- [ ] Code compiles without errors
- [ ] All tests pass
- [ ] New code has tests
- [ ] Documentation is updated
- [ ] Commit messages follow conventions

### PR Guidelines

1. **Title**: Use conventional commit format
   - `feat:` New feature
   - `fix:` Bug fix
   - `docs:` Documentation only
   - `test:` Test only
   - `refactor:` Code refactoring
   - `perf:` Performance improvement

2. **Description**: Include:
   - What changed and why
   - Related issue numbers
   - Breaking changes (if any)
   - Screenshots (if UI changes)

3. **Size**: Keep PRs focused and small

### Review Process

1. Automated checks must pass
2. At least one maintainer review required
3. Address review feedback
4. Maintainer merges when approved

## Coding Standards

### Go Style

Follow standard Go conventions:

```go
// Good
func CalculateBalance(transactions []Transaction) decimal.Decimal {
    balance := decimal.Zero
    for _, tx := range transactions {
        balance = balance.Add(tx.Amount)
    }
    return balance
}

// Bad
func calc_balance(t []Transaction) decimal.Decimal {
    var b decimal.Decimal
    for i := 0; i < len(t); i++ {
        b = b.Add(t[i].Amount)
    }
    return b
}
```

### Best Practices

1. **Error Handling**:
   ```go
   if err != nil {
       return fmt.Errorf("failed to generate statement: %w", err)
   }
   ```

2. **Comments**: Export functions need comments
   ```go
   // Generate creates a new statement from the provided input
   func (g *StatementGenerator) Generate(input StatementInput) (*Statement, error) {
   ```

3. **Tests**: Table-driven tests preferred
   ```go
   tests := []struct {
       name    string
       input   StatementInput
       want    *Statement
       wantErr bool
   }{
       // test cases
   }
   ```

### Formatting

Run before committing:
```bash
go fmt ./...
go vet ./...
```

## Testing Guidelines

### Test Coverage

Aim for >80% coverage for new code:

```bash
go test ./test/... -cover
```

### Test Types

1. **Unit Tests**: Test individual functions
2. **Integration Tests**: Test component interactions
3. **Benchmark Tests**: Performance testing
4. **Example Tests**: Documentation examples

### Writing Tests

```go
func TestStatementGenerator_Generate(t *testing.T) {
    // Arrange
    generator := New()
    input := StatementInput{
        // test data
    }

    // Act
    statement, err := generator.Generate(input)

    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, statement)
    assert.Equal(t, expected, statement.ClosingBalance)
}
```

### Test Data

Use fixtures for complex test data:

```go
func loadFixture(t *testing.T, name string) StatementInput {
    data, err := os.ReadFile(fmt.Sprintf("fixtures/%s.json", name))
    require.NoError(t, err)

    var input StatementInput
    err = json.Unmarshal(data, &input)
    require.NoError(t, err)

    return input
}
```

## Documentation

### Code Documentation

- All exported types, functions, and methods need comments
- Use complete sentences
- Include examples for complex APIs

```go
// Account represents a bank account with optional address for proof-of-address.
// The account number can be masked for security (e.g., "****1234").
//
// Example:
//
//	account := Account{
//	    Number:     "****1234",
//	    HolderName: "John Doe",
//	    Currency:   "USD",
//	    Address: &Address{
//	        Line1: "123 Main St",
//	        City:  "New York",
//	    },
//	}
type Account struct {
    // fields...
}
```

### README Updates

Update README.md when adding:
- New features
- API changes
- Examples
- Dependencies

### API Documentation

Update API.md for:
- New types
- New methods
- Breaking changes
- Deprecations

## Community

### Getting Help

- **Issues**: Use GitHub issues for bugs and features
- **Discussions**: Use GitHub Discussions for questions
- **Email**: Open a GitHub issue for support

### Recognition

Contributors are recognized in:
- README.md contributors section
- Release notes
- GitHub contributors page

## Release Process

1. **Version Bump**: Follow semantic versioning
2. **Changelog**: Update CHANGELOG.md
3. **Documentation**: Update version references
4. **Tag**: Create git tag
5. **Release**: Create GitHub release

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

---

Thank you for contributing to Statement Generator SDK! 🎉