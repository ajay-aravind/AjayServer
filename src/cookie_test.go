package main

import (
	"strings"
	"testing"
	"time"
)

func TestNewCookie_Valid(t *testing.T) {
	expiry := time.Now().Add(24 * time.Hour)
	c, err := NewCookie("token", "xyz123", "/", "example.com", expiry, 3600, true, true, "Lax", "/")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if c.GetKey() != "token" || c.GetValue() != "xyz123" {
		t.Errorf("Unexpected key/value: %s/%s", c.GetKey(), c.GetValue())
	}
}

func TestNewCookie_Invalid(t *testing.T) {
	_, err := NewCookie("", "value", "/", "example.com", time.Now(), 3600, true, true, "Lax", "/")
	if err == nil {
		t.Error("Expected error for empty key")
	}
}

func TestSetExpiryValidation(t *testing.T) {
	c := &Cookie{}
	err := c.SetExpiry(time.Now().Add(-1 * time.Hour))
	if err == nil {
		t.Error("Expected error for past expiry time")
	}
}

func TestSetSameSiteValidation(t *testing.T) {
	c := &Cookie{}
	err := c.SetSameSite("Invalid")
	if err == nil {
		t.Error("Expected error for invalid SameSite value")
	}
	err = c.SetSameSite("Lax")
	if err != nil {
		t.Errorf("Expected no error for valid SameSite, got %v", err)
	}
}

func TestToCookieString(t *testing.T) {
	expiry := time.Date(2025, 12, 31, 23, 59, 59, 0, time.UTC)
	c := &Cookie{
		key:      "user",
		value:    "abc123",
		path:     "/",
		domain:   "example.com",
		expiry:   expiry,
		maxAge:   3600,
		httpOnly: true,
		isSecure: true,
		sameSite: "Strict",
	}
	str := c.ToCookieString()
	expectedSubstrings := []string{
		"user=abc123",
		"Path=/",
		"Domain=example.com",
		"Expires=Wed, 31 Dec 2025 23:59:59 UTC",
		"Max-Age=3600",
		"HttpOnly",
		"Secure",
		"SameSite=Strict",
	}
	for _, s := range expectedSubstrings {
		if !strings.Contains(str, s) {
			t.Errorf("Cookie string missing expected part: %s", s)
		}
	}
}

func TestParseCookieBytes(t *testing.T) {
	raw := []byte("key=user; value=abc123; path=/;      domain=example.com; maxAge=3600; httpOnly=true; secure=true; sameSite=Lax; expiry=2025-12-31T23:59:59Z")
	c, err := ParseCookieBytes(raw, "/")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if c.GetKey() != "user" || c.GetValue() != "abc123" {
		t.Error("Parsed values do not match expected")
	}
	if !c.IsSecure() || !c.IsHttpOnly() {
		t.Error("Secure/HttpOnly flags not parsed correctly")
	}
}

func TestParseCookieBytes_EmptyString(t *testing.T) {
	_, err := ParseCookieBytes([]byte(""), "/")
	if err == nil {
		t.Error("Expected error for empty cookie string")
	}
}

func TestParseCookieBytes_MissingKey(t *testing.T) {
	_, err := ParseCookieBytes([]byte("value=abc; path=/"), "/")
	if err == nil {
		t.Error("Expected error for missing key")
	}
}

func TestParseCookieBytes_MissingValue(t *testing.T) {
	_, err := ParseCookieBytes([]byte("key=session; path=/"), "/")
	if err == nil {
		t.Error("Expected error for missing value")
	}
}

func TestParseCookieBytes_InvalidMaxAge(t *testing.T) {
	_, err := ParseCookieBytes([]byte("key=session; value=abc; maxAge=abc"), "/")
	if err == nil {
		t.Error("Expected error for non-numeric maxAge")
	}
}

func TestParseCookieBytes_InvalidExpiryFormat(t *testing.T) {
	_, err := ParseCookieBytes([]byte("key=session; value=abc; expiry=not-a-date"), "/")
	if err == nil {
		t.Error("Expected error for invalid expiry format")
	}
}

func TestParseCookieBytes_InvalidSameSite(t *testing.T) {
	_, err := ParseCookieBytes([]byte("key=session; value=abc; sameSite=Loosey"), "/")
	if err == nil {
		t.Error("Expected error for invalid SameSite value")
	}
}

func TestParseCookieBytes_ExtraSemicolonsAndSpaces(t *testing.T) {
	_, err := ParseCookieBytes([]byte(";; ; key = session ;; value = 123 ;; maxAge = 60 ;;"), "/")
	if err != nil {
		t.Errorf("Did not expect error for sloppy format: %v", err)
	}
}

func TestParseCookieBytes_WithGarbageKey(t *testing.T) {
	_, err := ParseCookieBytes([]byte("key=session; value=abc; g@rbage!?=value"), "/")
	if err != nil {
		t.Errorf("Unexpected error for unrecognized fields: %v", err)
	}
}

func TestParseCookieBytes_RandomCharactersInKey(t *testing.T) {
	_, err := ParseCookieBytes([]byte("key=!@#$%^&*(); value=abc"), "/")
	if err != nil {
		t.Errorf("Unexpected error for strange but syntactically valid input: %v", err)
	}
}

func TestParseCookieBytesSetDefaultPath(t *testing.T) {
	c, err := ParseCookieBytes([]byte("key=!@#$%^&*(); value=abc"), "/users")
	if err != nil {
		t.Errorf("Unexpected error for strange but syntactically valid input: %v", err)
	}

	if c.path != "/users" {
		t.Errorf("cookie path is not set to request location")
	}
}

func TestParseCookieBytesValidatePath(t *testing.T) {
	validPaths := []string{
		"/",
		"/path",
		"/path/",
		"/path/to/resource",
		"/path/with-hyphen",
		"/path/with.dot",
		"/path/with_underscore",
		"/path/with~tilde",
		"/path/with!exclamation",
		"/path/with$dollar",
		"/path/with&amp;ampersand",
		"/path/with'apostrophe",
		"/path/with(parentheses)",
		"/path/with*asterisk",
		"/path/with+plus",
		"/path/with,comma",
		"/path/with;semicolon",
		"/path/with=equals",
		"/path/with:colon",
		"/path/with@at",
		"/path/with%20space", // Percent-encoded space
		"/p%61th",            // Percent-encoded 'a'
	}

	for _, path := range validPaths {
		_, err := ParseCookieBytes([]byte("key=!@#$%^&*(); value=abc;path="+path), "/users")
		if err != nil {
			t.Errorf("cookie path should be parse without any issues")
		}
	}

	cookie := Cookie{}

	invalidPaths := []string{
		"relative/path",
		"/path with space",
		"/path%",
		"/path%G0",
		"//path",
		"/path?query=value",
		"/path#fragment",
		" path/leading-space",
	}

	for _, path := range invalidPaths {
		err := cookie.SetPath(path, "/")
		if err == nil {
			t.Errorf("cookie path should throw validation error")
		}
	}
}

func TestParseCookieBytesValidateDomain(t *testing.T) {
	var validDomains = []string{
		"example.com",
		".example.com",
		"sub.example.com",
		"example.co.uk",
		"localhost",
		"127.0.0.1",
		"[::1]",
		"xn--bcher-kva.example",
		"my-domain123.org",
		"EXAMPLE.COM",
		"", // Empty string, because by default cookie domain is request location
	}

	cookie := Cookie{}
	for _, domain := range validDomains {
		err := cookie.SetDomain(domain)
		if err != nil {
			t.Errorf("cookie domain should parse without any issues:" + domain)
		}
	}

	var invalidDomains = []string{
		"example",            // TLD-only
		".com",               // Public suffix only
		"-invalid.com",       // Starts with hyphen
		"example..com",       // Double dot
		"ex ample.com",       // Space in domain
		"example.com.",       // Trailing dot
		"http://example.com", // Scheme included
		"/example.com",       // Slash included
		"example,com",        // Comma separator
		"@example.com",       // Invalid char "@"
		"example!.com",       // Special character "!"
		"1234",               // Number, not a valid domain
		"256.300.999.1",      // Invalid IP address
	}

	for _, domain := range invalidDomains {
		err := cookie.SetDomain(domain)
		if err == nil {
			t.Errorf("cookie domain parsing should fail: " + domain)
		}
	}
}
