package main

import (
	"fmt"
	"net"
	"net/url"
	"regexp"
	"strings"

	"golang.org/x/net/publicsuffix"
)

func main() {
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
`

	// Matches URLs with http/https and also plain domains like foo.bar
	re := regexp.MustCompile(`https?://[^\s]+|[a-zA-Z0-9][-a-zA-Z0-9]*(\.[a-zA-Z0-9][-a-zA-Z0-9]*)+`)

	matches := re.FindAllString(text, -1)

	for _, raw := range matches {
		raw = strings.TrimRight(raw, ".,)") // Strip trailing punctuation
		u, err := url.Parse(raw)

		if err != nil || u.Host == "" {
			// Might be a naked domain like "example.com"
			raw = "http://" + raw
			u, err = url.Parse(raw)
			if err != nil {
				fmt.Printf("Rejected (parse error): %q\n", raw)
				continue
			}
		}

		host := u.Hostname()

		if ip := net.ParseIP(host); ip != nil {
			fmt.Printf("✅ Valid IP-based URL: %s (IP: %s)\n", raw, ip)
			continue
		}

		// Validate domain using publicsuffix
		if eTLD, icann := publicsuffix.PublicSuffix(host); eTLD != "" {
			if icann {
				fmt.Printf("✅ Valid domain-based URL: %s (TLD: %s)\n", raw, eTLD)
			} else {
				// For non-ICANN eTLD, check if it's built on a valid ICANN TLD
				// e.g., "foo.dyndns.org" -> eTLD is "dyndns.org", check if ".org" is ICANN
				if strings.Contains(eTLD, ".") {
					parts := strings.Split(eTLD, ".")
					actualTLD := parts[len(parts)-1]
					// Test if this actual TLD is an ICANN TLD
					testDomain := "test." + actualTLD
					if _, testICANN := publicsuffix.PublicSuffix(testDomain); testICANN {
						fmt.Printf("✅ Valid domain-based URL: %s (TLD: %s, built on ICANN TLD: %s)\n", raw, eTLD, actualTLD)
					} else {
						fmt.Printf("❌ Rejected (invalid domain/TLD): %s\n", raw)
					}
				} else {
					fmt.Printf("❌ Rejected (invalid domain/TLD): %s\n", raw)
				}
			}
		} else {
			fmt.Printf("❌ Rejected (invalid domain/TLD): %s\n", raw)
		}
	}
}
