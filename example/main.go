package main

import (
	"fmt"

	"github.com/potakhov/urlverify"
)

func main() {
	// Example text with various URLs and domains
	text := `
	Visit our website at https://example.com for more info.
	You can also check out subdomain.example.org and our 
	mirror at backup.dyndns.org. Our server runs on 
	192.168.1.100:8080 and we also support IPv6 at 
	[2001:db8::1]:443.
	
	Test case-insensitive schemes: HTTP://github.com and HTTPS://stackoverflow.com
	
	Ignore these invalid ones: test.local, fake.invalid
	`

	fmt.Println("=== URLVerify Library Example ===")
	fmt.Println()

	// Extract all valid URLs/domains
	fmt.Println("1. Extract all valid URLs/domains:")
	validDomains := urlverify.ExtractAll(text)
	for i, domain := range validDomains {
		fmt.Printf("   %d. %s\n", i+1, domain)
	}

	fmt.Println("\n2. Detailed validation for each found item:")
	for _, domain := range validDomains {
		result := urlverify.ValidateDomain(domain)
		fmt.Printf("   ✅ %s\n", domain)
		fmt.Printf("      Type: %s, TLD: %s, Reason: %s\n", result.Type, result.TLD, result.Reason)
	}

	fmt.Println("\n3. Test individual domain validation:")
	testCases := []string{
		"google.com",
		"test.dyndns.org",
		"192.168.1.1",
		"invalid.fake",
		"justtext",
	}

	for _, testCase := range testCases {
		result := urlverify.ValidateDomain(testCase)
		status := "❌"
		if result.Valid {
			status = "✅"
		}
		fmt.Printf("   %s %s -> %s\n", status, testCase, result.Reason)
	}
}
