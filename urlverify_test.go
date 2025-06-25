package urlverify

import (
	"testing"
)

func TestValidateDomain(t *testing.T) {
	tests := []struct {
		input          string
		expectedValid  bool
		expectedReason string
		expectedType   URLType
		description    string
	}{
		{
			input:          "example.com",
			expectedValid:  true,
			expectedReason: "valid ICANN domain",
			expectedType:   URLTypeICANN,
			description:    "basic .com domain",
		},
		{
			input:          "https://example.co.uk/path?query=123",
			expectedValid:  true,
			expectedReason: "valid ICANN domain",
			expectedType:   URLTypeICANN,
			description:    "full URL with multi-part TLD",
		},
		{
			input:          "192.168.1.1:8080",
			expectedValid:  true,
			expectedReason: "valid IP address",
			expectedType:   URLTypeIP,
			description:    "IPv4 address with port",
		},
		{
			input:          "[2001:db8::1]:443",
			expectedValid:  true,
			expectedReason: "valid IP address",
			expectedType:   URLTypeIP,
			description:    "IPv6 address with port",
		},
		{
			input:          "foo.dyndns.org",
			expectedValid:  true,
			expectedReason: "valid domain built on ICANN TLD",
			expectedType:   URLTypeNonICANN,
			description:    "dynamic DNS service on .org",
		},
		{
			input:          "test.local",
			expectedValid:  false,
			expectedReason: "invalid or non-ICANN TLD",
			expectedType:   URLTypeInvalid,
			description:    ".local domain (mDNS)",
		},
		{
			input:          "not_a_valid_domain.dse",
			expectedValid:  false,
			expectedReason: "invalid or non-ICANN TLD",
			expectedType:   URLTypeInvalid,
			description:    "invalid TLD",
		},
		{
			input:          "justtext",
			expectedValid:  false,
			expectedReason: "no valid TLD found",
			expectedType:   URLTypeInvalid,
			description:    "plain text without domain structure",
		},
		{
			input:          "weird:123",
			expectedValid:  false,
			expectedReason: "no valid TLD found",
			expectedType:   URLTypeInvalid,
			description:    "invalid format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			result := ValidateDomain(tt.input)

			if result.Valid != tt.expectedValid {
				t.Errorf("ValidateDomain(%q) valid = %v, want %v", tt.input, result.Valid, tt.expectedValid)
			}

			if result.Reason != tt.expectedReason {
				t.Errorf("ValidateDomain(%q) reason = %q, want %q", tt.input, result.Reason, tt.expectedReason)
			}

			if result.Type != tt.expectedType {
				t.Errorf("ValidateDomain(%q) type = %q, want %q", tt.input, result.Type, tt.expectedType)
			}
		})
	}
}

func TestExtractAll(t *testing.T) {
	text := `
Here are some links:
- http://example.com
- https://example.co.uk/path?query=123
- http://192.168.1.1:8080/dashboard
- http://[2001:db8::1]:443/page
- weird:123
- justtext
- foo.bar
- test.local
- http://not_a_valid_domain.dse
- foo.dyndns.org
- http://foo.dot.uno/index.html
- https://example.com/path/to/resource.html
- https://dns.17ex.com
- example.com/index.html
- server.com:8080/index.html
`

	expected := []string{
		"http://example.com",
		"https://example.co.uk/path?query=123",
		"http://192.168.1.1:8080/dashboard",
		"http://[2001:db8::1]:443/page",
		"foo.bar",
		"foo.dyndns.org",
		"http://foo.dot.uno/index.html",
		"https://example.com/path/to/resource.html",
		"https://dns.17ex.com",
		"example.com/index.html",
		"server.com:8080/index.html",
	}

	result := ExtractAll(text)

	if len(result) != len(expected) {
		t.Errorf("ExtractAll() returned %d results, want %d", len(result), len(expected))
		t.Logf("Got: %v", result)
		t.Logf("Expected: %v", expected)
		return
	}

	for i, want := range expected {
		if i >= len(result) || result[i] != want {
			t.Errorf("ExtractAll() result[%d] = %q, want %q", i, result[i], want)
		}
	}
}

func TestExtractAllExactText(t *testing.T) {
	// Test that ExtractAll returns domains exactly as they appear in text
	text := "Visit example.com and also check https://google.com/search?q=test"

	result := ExtractAll(text)
	expected := []string{"example.com", "https://google.com/search?q=test"}

	if len(result) != len(expected) {
		t.Errorf("ExtractAll() returned %d results, want %d", len(result), len(expected))
		return
	}

	for i, want := range expected {
		if result[i] != want {
			t.Errorf("ExtractAll() result[%d] = %q, want %q (should be exact text from input)", i, result[i], want)
		}
	}
}
