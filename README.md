# URLVerify

A Go library for extracting and validating URLs and domains from text.

## Features

- Extracts URLs (with http/https) and plain domains from text
- Validates domains using the Public Suffix List
- Supports IPv4 and IPv6 addresses
- Handles dynamic DNS services (e.g., dyndns.org, no-ip.org) 
- Returns domains exactly as they appear in the original text
- Provides detailed validation results for testing and debugging

## Usage

### Basic Usage

```go
package main

import (
    "fmt"
    "github.com/potakhov/urlverify"
)

func main() {
    text := `
    Check out these sites:
    - https://example.com
    - foo.dyndns.org  
    - 192.168.1.1:8080
    - invalid.fake
    `
    
    validDomains := urlverify.ExtractAll(text)
    fmt.Printf("Found %d valid URLs/domains:\n", len(validDomains))
    for _, domain := range validDomains {
        fmt.Printf("- %s\n", domain)
    }
}
```

### Detailed Validation

```go
result := urlverify.ValidateDomain("foo.dyndns.org")
if result.Valid {
    fmt.Printf("✅ Valid %s: %s (TLD: %s)\n", result.Type, "foo.dyndns.org", result.TLD)
} else {
    fmt.Printf("❌ Invalid: %s (%s)\n", "foo.dyndns.org", result.Reason)
}
```

## API

### `ExtractAll(text string) []string`

Extracts all valid URLs and domains from the given text, returning them exactly as they appeared in the original text.

### `ValidateDomain(domain string) ValidationResult`

Validates a single URL or domain string and returns detailed validation information.

### `ValidationResult`

```go
type ValidationResult struct {
    Valid  bool    // Whether the domain is valid
    Reason string  // Explanation of the validation result
    Type   URLType // type of the URL
    TLD    string  // The effective TLD or IP address
}
```

## Testing

Run tests with:

```bash
go test ./pkg/urlverify
```