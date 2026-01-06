---
sidebar_position: 4
---

# Testing

Guide to testing tf2report.

## Running Tests

### All Tests

```bash
make test
```

### Specific Package

```bash
go test ./pkg/terraform
go test ./pkg/report
go test ./pkg/config
```

### Verbose Output

```bash
go test -v ./...
```

## Test Coverage

### Generate Coverage

```bash
make test-coverage
```

Opens `coverage.html` in browser.

### View Coverage in Terminal

```bash
go test -cover ./...
```

### Detailed Coverage

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Writing Tests

### Unit Tests

```go
// pkg/terraform/parser_test.go
package terraform

import "testing"

func TestParseFile(t *testing.T) {
    parser := NewParser()
    plan, err := parser.ParseFile("testdata/sample.json")

    if err != nil {
        t.Fatalf("expected no error, got %v", err)
    }

    if plan == nil {
        t.Fatal("expected plan, got nil")
    }
}
```

### Table-Driven Tests

```go
func TestFilterByType(t *testing.T) {
    tests := []struct {
        name     string
        changes  []ResourceChange
        types    []string
        expected int
    }{
        {
            name: "filter single type",
            changes: []ResourceChange{
                {Type: "aws_instance"},
                {Type: "aws_s3_bucket"},
            },
            types: []string{"aws_instance"},
            expected: 1,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := FilterByType(tt.changes, tt.types)
            if len(result) != tt.expected {
                t.Errorf("expected %d, got %d", tt.expected, len(result))
            }
        })
    }
}
```

## Test Data

Test data located in package-specific `testdata/` directories:

```
pkg/terraform/testdata/
├── sample-plan.json
├── empty-plan.json
└── complex-plan.json
```

## Integration Tests

Test with example files:

```bash
./bin/tf2report --plan examples/sample-plan.json
./bin/tf2report --plan examples/sample-plan.json --format text
./bin/tf2report --plan examples/sample-plan.json --format json
```

## Continuous Integration

Tests run automatically on:
- Pull requests
- Pushes to main branch
- Release builds

## Next Steps

- **[API Reference](./api-reference.md)** - Package documentation
- **[Contributing](./contributing.md)** - Contribution guidelines
