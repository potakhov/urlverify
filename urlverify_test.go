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
- https://مثال.السعودية
- https://книга.рф
- https://스타벅스코리아.com
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
		"https://مثال.السعودية",
		"https://книга.рф",
		"https://스타벅스코리아.com",
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

// TestCaseInsensitiveSchemes tests that URLs with uppercase or mixed-case schemes are properly detected
func TestCaseInsensitiveSchemes(t *testing.T) {
	text := `
Test case-insensitive URL schemes:
- HTTP://example.com
- HTTPS://example.co.uk/path
- HttpS://github.com
- hTtP://google.com
- HtTpS://stackoverflow.com/questions
`

	expected := []string{
		"HTTP://example.com",
		"HTTPS://example.co.uk/path",
		"HttpS://github.com",
		"hTtP://google.com",
		"HtTpS://stackoverflow.com/questions",
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

// Benchmark data
var (
	// Text without any URLs or domains
	textWithoutURLs = `
		This is a sample text that contains no URLs or domains whatsoever.
		It has multiple sentences and paragraphs to simulate real world text processing.
		The content talks about various topics like technology, science, literature.
		There are numbers like 123, 456, and special characters like @#$%^&*().
		This text is designed to test performance when no URL extraction is needed.
		Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor.
		Incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam quis.
		Nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.
		Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore.
		Eu fugiat nulla pariatur excepteur sint occaecat cupidatat non proident sunt.
	`

	// Text with a few URLs
	textWithFewURLs = `
		Check out these websites: https://example.com and http://google.com.
		You can also visit github.com for code repositories.
		The IP address 192.168.1.1 is commonly used for local networks.
		Don't forget about https://stackoverflow.com/questions for programming help.
		Some additional text to make this more realistic for benchmarking purposes.
	`

	// Text with many URLs
	textWithManyURLs = `
		Here are numerous websites to check out:
		https://example.com, http://google.com, https://github.com/user/repo,
		stackoverflow.com, https://reddit.com/r/programming, news.ycombinator.com,
		medium.com, dev.to, https://techcrunch.com, vercel.com,
		netlify.com, https://aws.amazon.com, cloud.google.com,
		azure.microsoft.com, https://digitalocean.com, heroku.com,
		https://stripe.com, paypal.com, https://shopify.com,
		192.168.1.1, 10.0.0.1, https://localhost:3000,
		test.dyndns.org, foo.no-ip.org, https://example.co.uk,
		subdomain.example.org, https://api.service.com/v1/endpoint,
		cdn.jsdelivr.net, unpkg.com, https://fonts.googleapis.com,
		bootstrap.com, jquery.com, https://reactjs.org,
		vuejs.org, angular.io, https://svelte.dev
	`

	// Mixed content with valid and invalid domains
	textMixed = `
		Valid domains: example.com, https://google.com, github.com
		Invalid stuff: just_text, weird:123, incomplete.
		IP addresses: 192.168.1.1, [2001:db8::1]
		Dynamic DNS: test.dyndns.org, home.no-ip.org
		More text content to simulate real world scenarios.
		Some numbers: 12345, dates: 2023-12-25, emails: user@domain.com
	`

	// Individual domain test cases
	domainExamples = []string{
		"example.com",
		"https://google.com/search?q=test",
		"192.168.1.1",
		"[2001:db8::1]:443",
		"test.dyndns.org",
		"invalid.local",
		"just_text",
		"subdomain.example.co.uk",
		"api.service.com:8080/v1/endpoint",
		"cdn.example.org",
	}
)

// BenchmarkExtractAll_NoURLs measures performance when no URLs are present
func BenchmarkExtractAll_NoURLs(b *testing.B) {
	for i := 0; i < b.N; i++ {
		result := ExtractAll(textWithoutURLs)
		_ = result // Prevent optimization
	}
}

// BenchmarkExtractAll_FewURLs measures performance with a few URLs
func BenchmarkExtractAll_FewURLs(b *testing.B) {
	for i := 0; i < b.N; i++ {
		result := ExtractAll(textWithFewURLs)
		_ = result
	}
}

// BenchmarkExtractAll_ManyURLs measures performance with many URLs
func BenchmarkExtractAll_ManyURLs(b *testing.B) {
	for i := 0; i < b.N; i++ {
		result := ExtractAll(textWithManyURLs)
		_ = result
	}
}

// BenchmarkExtractAll_Mixed measures performance with mixed valid/invalid content
func BenchmarkExtractAll_Mixed(b *testing.B) {
	for i := 0; i < b.N; i++ {
		result := ExtractAll(textMixed)
		_ = result
	}
}

// BenchmarkValidateDomain_ICANN tests validation of ICANN domains
func BenchmarkValidateDomain_ICANN(b *testing.B) {
	for i := 0; i < b.N; i++ {
		result := ValidateDomain("example.com")
		_ = result
	}
}

// BenchmarkValidateDomain_URL tests validation of full URLs
func BenchmarkValidateDomain_URL(b *testing.B) {
	for i := 0; i < b.N; i++ {
		result := ValidateDomain("https://example.com/path?query=value")
		_ = result
	}
}

// BenchmarkValidateDomain_IP tests validation of IP addresses
func BenchmarkValidateDomain_IP(b *testing.B) {
	for i := 0; i < b.N; i++ {
		result := ValidateDomain("192.168.1.1")
		_ = result
	}
}

// BenchmarkValidateDomain_DynamicDNS tests validation of dynamic DNS domains
func BenchmarkValidateDomain_DynamicDNS(b *testing.B) {
	for i := 0; i < b.N; i++ {
		result := ValidateDomain("test.dyndns.org")
		_ = result
	}
}

// BenchmarkValidateDomain_Invalid tests validation of invalid domains
func BenchmarkValidateDomain_Invalid(b *testing.B) {
	for i := 0; i < b.N; i++ {
		result := ValidateDomain("invalid.local")
		_ = result
	}
}

// BenchmarkValidateDomain_Mixed tests validation across various domain types
func BenchmarkValidateDomain_Mixed(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, domain := range domainExamples {
			result := ValidateDomain(domain)
			_ = result
		}
	}
}

// BenchmarkRegexOnly measures just the regex matching performance
func BenchmarkRegexOnly_NoURLs(b *testing.B) {
	for i := 0; i < b.N; i++ {
		matches := urlRegex.FindAllString(textWithoutURLs, -1)
		_ = matches
	}
}

func BenchmarkRegexOnly_ManyURLs(b *testing.B) {
	for i := 0; i < b.N; i++ {
		matches := urlRegex.FindAllString(textWithManyURLs, -1)
		_ = matches
	}
}

// BenchmarkLargeText tests performance with very large input
func BenchmarkExtractAll_LargeText(b *testing.B) {
	// Create a large text by repeating the mixed content
	largeText := ""
	for i := 0; i < 100; i++ {
		largeText += textMixed + "\n"
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := ExtractAll(largeText)
		_ = result
	}
}

// Comparative benchmarks to measure performance difference
// between text with URLs vs text without URLs

// BenchmarkComparison_WithoutURLs_Short tests short text without URLs
func BenchmarkComparison_WithoutURLs_Short(b *testing.B) {
	shortText := "This is a short text without any URLs or domains."
	for i := 0; i < b.N; i++ {
		result := ExtractAll(shortText)
		_ = result
	}
}

// BenchmarkComparison_WithURLs_Short tests short text with URLs
func BenchmarkComparison_WithURLs_Short(b *testing.B) {
	shortText := "Check out https://example.com and google.com for more info."
	for i := 0; i < b.N; i++ {
		result := ExtractAll(shortText)
		_ = result
	}
}

// BenchmarkComparison_WithoutURLs_Medium tests medium text without URLs
func BenchmarkComparison_WithoutURLs_Medium(b *testing.B) {
	mediumText := `
		This is a medium-sized text that contains no URLs or domains.
		It discusses various topics like technology, science, and literature.
		The text includes numbers, punctuation, and multiple sentences.
		It simulates typical content that might be processed in real applications.
		This helps us understand the baseline performance when no URL processing is needed.
	`
	for i := 0; i < b.N; i++ {
		result := ExtractAll(mediumText)
		_ = result
	}
}

// BenchmarkComparison_WithURLs_Medium tests medium text with URLs
func BenchmarkComparison_WithURLs_Medium(b *testing.B) {
	mediumText := `
		Check out these websites: https://example.com, google.com, and github.com.
		You can find more information at stackoverflow.com and reddit.com.
		For development tools, visit https://developer.mozilla.org and docs.microsoft.com.
		The server at 192.168.1.1 hosts our internal applications.
		Don't forget to check test.dyndns.org for the latest updates.
	`
	for i := 0; i < b.N; i++ {
		result := ExtractAll(mediumText)
		_ = result
	}
}
