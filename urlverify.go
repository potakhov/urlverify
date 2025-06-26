// Package urlverify provides functionality to extract and validate URLs and domains from text.
//
// The package uses the Public Suffix List to validate domains and supports:
// - HTTP/HTTPS URLs
// - Plain domain names
// - IPv4 and IPv6 addresses
// - Dynamic DNS services (e.g., dyndns.org, no-ip.org)
//
// Example usage:
//
//	text := "Visit https://example.com and check foo.dyndns.org"
//	validDomains := urlverify.ExtractAll(text)
//
//	Returns: ["https://example.com", "foo.dyndns.org"]
//
//	result := urlverify.ValidateDomain("foo.dyndns.org")
//	if result.Valid {
//	    fmt.Printf("Valid domain: %s (TLD: %s)\n", "foo.dyndns.org", result.TLD)
//	}
package urlverify

import (
	"net"
	"net/url"
	"regexp"
	"strings"

	"golang.org/x/net/publicsuffix"
)

var urlRegex = regexp.MustCompile(`https?://[^\s]+|(?:\[[0-9a-fA-F:]+\]|\d{1,3}(?:\.\d{1,3}){3}|[a-zA-Z0-9][-a-zA-Z0-9]*(?:\.[a-zA-Z0-9][-a-zA-Z0-9]*)+)(?::\d+)?(?:/[^\s]*)?`)

type URLType int

const (
	URLTypeInvalid URLType = iota
	URLTypeIP
	URLTypeICANN
	URLTypeNonICANN
)

func (t URLType) String() string {
	switch t {
	case URLTypeInvalid:
		return "Invalid"
	case URLTypeIP:
		return "IP Address"
	case URLTypeICANN:
		return "ICANN Domain"
	case URLTypeNonICANN:
		return "Non-ICANN Domain"
	default:
		return "Unknown"
	}
}

// ValidationResult represents the result of domain validation
type ValidationResult struct {
	Valid  bool     // Whether the URL or domain is valid
	Reason string   // Explanation of the validation result
	Type   URLType  // Type of URL or domain
	TLD    string   // The effective TLD, if applicable or an IP address
	URL    *url.URL // Only set if the URL was successfully parsed
}

// ExtractAll extracts and validates all URLs and domains from the given text,
// returning them exactly as they appeared in the original text (without adding schema)
func ExtractAll(text string) []string {
	matches := urlRegex.FindAllString(text, -1)

	var validURLs []string
	for _, raw := range matches {
		raw = strings.TrimRight(raw, ".,)") // Strip trailing punctuation
		if result := ValidateDomain(raw); result.Valid {
			validURLs = append(validURLs, raw)
		}
	}

	return validURLs
}

// ValidateDomain validates a single URL or domain string and returns detailed validation result
func ValidateDomain(raw string) ValidationResult {
	// Try to parse as-is first
	url, err := url.Parse(raw)

	if err != nil || url.Host == "" {
		// Might be a naked domain like "example.com"
		testURL := "http://" + raw
		url, err = url.Parse(testURL)
		if err != nil {
			return ValidationResult{
				Valid:  false,
				Reason: "parse error: " + err.Error(),
				Type:   URLTypeInvalid,
			}
		}
	}

	// Check if it's an IP address
	if ip := net.ParseIP(url.Hostname()); ip != nil {
		return ValidationResult{
			Valid:  true,
			Reason: "valid IP address",
			Type:   URLTypeIP,
			TLD:    ip.String(),
			URL:    url,
		}
	}

	// Validate domain using publicsuffix
	return validateDomainName(url)
}

// validateDomainName validates a domain name using the public suffix list
func validateDomainName(url *url.URL) ValidationResult {
	hostname := url.Hostname()

	// Handle edge cases first
	if hostname == "" {
		return ValidationResult{
			Valid:  false,
			Reason: "empty hostname",
			Type:   URLTypeInvalid,
		}
	}

	// Check if it has any dots - if not, it's not a valid domain
	if !strings.Contains(hostname, ".") {
		return ValidationResult{
			Valid:  false,
			Reason: "no valid TLD found",
			Type:   URLTypeInvalid,
		}
	}

	eTLD, icann := publicsuffix.PublicSuffix(hostname)

	if eTLD == "" {
		return ValidationResult{
			Valid:  false,
			Reason: "no valid TLD found",
			Type:   URLTypeInvalid,
		}
	}

	if icann {
		return ValidationResult{
			Valid:  true,
			Reason: "valid ICANN domain",
			Type:   URLTypeICANN,
			TLD:    eTLD,
			URL:    url,
		}
	}

	// For non-ICANN eTLD, check if it's built on a valid ICANN TLD
	// e.g., "foo.dyndns.org" -> eTLD is "dyndns.org", check if ".org" is ICANN
	if strings.Contains(eTLD, ".") {
		parts := strings.Split(eTLD, ".")
		actualTLD := parts[len(parts)-1]
		// Test if this actual TLD is an ICANN TLD
		testDomain := "test." + actualTLD
		if _, testICANN := publicsuffix.PublicSuffix(testDomain); testICANN {
			return ValidationResult{
				Valid:  true,
				Reason: "valid domain built on ICANN TLD",
				Type:   URLTypeNonICANN,
				TLD:    eTLD,
				URL:    url,
			}
		}
	}

	return ValidationResult{
		Valid:  false,
		Reason: "invalid or non-ICANN TLD",
		Type:   URLTypeInvalid,
		TLD:    eTLD,
	}
}
