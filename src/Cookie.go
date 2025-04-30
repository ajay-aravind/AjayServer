/*
Cookie
    CreateCookie
    Copy
    GetDomain()
    GetExpiry()
    GetHttpOnly()
    GetKey()
    GetMaxAge()
    ParseCookieString()
    ParseCookieBytes()
    IsPartitioned() // if it's partitioned cookie, it's set per top level domain. If abc.com loads a thirdparty script from cdn.com and cdn.com set's a cookie, usedid=123, and
    If i got to antoerh website xyz.com and loads same script cdn.com and cdn.com set's a cookie, usedid=456, then these cookies will be seperate. Even if the cookie and domain which
    set them are same i.e usedid, cdn.com, but top level domain in which the cdn.com is loaded is different. First cookie set in abc.com while second cookie set in xyz.com, so these two
    cookies even though they are same, they act like 2 different cookies. now cdn.com can't track users across sites.
    // It's more of a browser feature, don't really matter in backend i guess apart from the parsing part.
    Reset()
    SameSite() // Strict, lax, None
    //strict, send only with requests orginating for this domain i.e only when user trying to access this domain
    // lax, send cookies even if the user redirected from other site. user clicked this domain link from other domain tab. But
    not allowed to send, when a other domain web page sending ajax, post requests to this domain
    // none, Cookies will always be sent with cross-site requests. Need secure to send the requests.
    IsSecure() // Only send cookie via https, don't send via http requests
    SetDomain()
    SetExpiry()
    setHttonly()
    SetKey()
    SetMaxAge()
    SetParitioned()
    setPath() // The Path attribute in cookies determines which URL paths on a website can access the cookie.
    SetSameSite()
    SetSecure()
    SetValue()
    ToString()
    GetValue()
    WriteTo()
*/

package main

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type Cookie struct {
	key    string
	value  string
	path   string
	domain string
	expiry time.Time

	maxAge   int
	httpOnly bool
	isSecure bool
	sameSite string
}

func NewCookie(
	key string,
	value string,
	path string,
	domain string,
	expiry time.Time,
	maxAge int,
	httpOnly bool,
	isSecure bool,
	sameSite string,
	requestLocation string,
) (*Cookie, error) {

	cookie := &Cookie{}

	var err error
	err = cookie.SetKey(key)
	if err != nil {
		return nil, err
	}
	err = cookie.SetValue(value)
	if err != nil {
		return nil, err
	}

	if !expiry.IsZero() {
		err = cookie.SetExpiry(expiry)
		if err != nil {
			return nil, err
		}
	}

	err = cookie.SetMaxAge(maxAge)
	if err != nil {
		return nil, err
	}
	err = cookie.SetSameSite(sameSite)
	if err != nil {
		return nil, err
	}

	// empty domain is allowed
	err = cookie.SetDomain(domain)

	if err != nil {
		return nil, err
	}

	err = cookie.SetPath(path, requestLocation)

	if err != nil {
		return nil, err
	}

	cookie.SetHttpOnly(httpOnly)
	cookie.SetSecure(isSecure)
	return cookie, err
}

// ParseCookieBytes parses a byte slice of a cookie into a Cookie struct.
func ParseCookieBytes(cookieBytes []byte, requestLocation string) (*Cookie, error) {
	// Convert the byte slice to string
	cookieStr := string(cookieBytes)

	// Initialize default cookie fields
	var key, value, path, domain, sameSite string
	var expiry time.Time
	var maxAge int
	var httpOnly, isSecure bool

	// Split the cookie string into key-value pairs (separated by ';')
	parts := strings.Split(cookieStr, ";")
	for _, part := range parts {
		// Trim spaces around each part
		part = strings.TrimSpace(part)

		// Split the part by '=' to get key and value
		keyValue := strings.SplitN(part, "=", 2)

		// Ensure the keyValue is valid (contains both key and value)
		if len(keyValue) != 2 {
			continue
		}

		// Trim spaces around the key and value
		k, v := strings.TrimSpace(keyValue[0]), strings.TrimSpace(keyValue[1])

		// Map the key-value pair to the appropriate cookie fields
		switch k {
		case "key":
			key = v
		case "value":
			value = v
		case "path":
			path = v
		case "domain":
			domain = v
		case "maxAge":
			// Convert maxAge to integer
			var err error
			maxAge, err = strconv.Atoi(v)
			if err != nil {
				return nil, fmt.Errorf("invalid maxAge value: %v", err)
			}
		case "httpOnly":
			// Convert to boolean
			httpOnly = v == "true"
		case "secure":
			// Convert to boolean
			isSecure = v == "true"
		case "sameSite":
			sameSite = v
		case "expiry":
			// Parse the expiry date string into time
			expiry2, err := time.Parse(time.RFC3339, v)
			expiry = expiry2
			if err != nil {
				return nil, fmt.Errorf("invalid expiry date format: %v", err)
			}
		default:
			logrus.Trace("unknown cookie field, ignoreing:", k)
		}
	}

	return NewCookie(
		key,
		value,
		path,
		domain,
		expiry,
		maxAge,
		httpOnly,
		isSecure,
		sameSite,
		requestLocation)
}

func (c *Cookie) ToCookieString() string {
	var b strings.Builder

	// Required key-value pair
	b.WriteString(fmt.Sprintf("%s=%s", c.key, c.value))

	// Optional attributes
	if c.path != "" {
		b.WriteString(fmt.Sprintf("; Path=%s", c.path))
	}
	if c.domain != "" {
		b.WriteString(fmt.Sprintf("; Domain=%s", c.domain))
	}
	if !c.expiry.IsZero() {
		b.WriteString(fmt.Sprintf("; Expires=%s", c.expiry.UTC().Format(time.RFC1123)))
	}
	if c.maxAge > 0 {
		b.WriteString(fmt.Sprintf("; Max-Age=%d", c.maxAge))
	}
	if c.httpOnly {
		b.WriteString("; HttpOnly")
	}
	if c.isSecure {
		b.WriteString("; Secure")
	}
	if c.sameSite != "" {
		b.WriteString(fmt.Sprintf("; SameSite=%s", c.sameSite))
	}

	return b.String()
}

func (c *Cookie) GetKey() string {
	return c.key
}

func (c *Cookie) GetValue() string {
	return c.value
}

func (c *Cookie) GetPath() string {
	return c.path
}

func (c *Cookie) GetDomain() string {
	return c.domain
}

func (c *Cookie) GetExpiry() time.Time {
	return c.expiry
}

func (c *Cookie) GetMaxAge() int {
	return c.maxAge
}

func (c *Cookie) IsHttpOnly() bool {
	return c.httpOnly
}

func (c *Cookie) IsSecure() bool {
	return c.isSecure
}

func (c *Cookie) GetSameSite() string {
	return c.sameSite
}

// Setters with basic validation
func (c *Cookie) SetKey(key string) error {
	if key == "" {
		return errors.New("key cannot be empty")
	}
	c.key = key
	return nil
}

func (c *Cookie) SetValue(value string) error {
	if value == "" {
		return errors.New("value cannot be empty")
	}
	c.value = value
	return nil
}

func (c *Cookie) SetPath(path string, requestLocation string) error {
	//todo: is this needed?
	if path == "" {
		path = requestLocation
	}

	if !isValidCookiePath(path) {
		return errors.New("Invalid cookie path")
	}

	c.path = path
	return nil
}

func isValidDomain(domain string) bool {
	// Regex to validate domain names (as discussed previously)
	var domainRegex = `^(localhost|127\.0\.0\.1|\[::1\]|(\.?[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,})$`
	re := regexp.MustCompile(domainRegex)
	return re.MatchString(domain)
}

func isValidCookiePath(path string) bool {
	pathRegex := `^/(?:(?:(?:[a-zA-Z0-9\-._~!$&'()*+,;=:@]|%[0-9a-fA-F]{2})/?)*)?$`
	re := regexp.MustCompile(pathRegex)
	return re.MatchString(path)
}

// todo: Cookie domain must match or be a parent of the request's host (e.g., you cannot set a cookie for example.com from evil.com)
func (c *Cookie) SetDomain(domain string) error {
	if len(domain) <= 0 {
		return nil
	}
	if !isValidDomain(domain) {
		return errors.New("invalid domain name")
	}

	c.domain = domain
	return nil
}

func (c *Cookie) SetExpiry(expiry time.Time) error {
	if expiry.Before(time.Now()) {
		return errors.New("expiry time cannot be in the past")
	}
	c.expiry = expiry
	return nil
}

func (c *Cookie) SetMaxAge(age int) error {
	if age < 0 {
		return errors.New("maxAge cannot be negative")
	}
	c.maxAge = age
	return nil
}

func (c *Cookie) SetHttpOnly(httpOnly bool) {
	c.httpOnly = httpOnly
}

func (c *Cookie) SetSecure(isSecure bool) {
	c.isSecure = isSecure
}

func (c *Cookie) SetSameSite(mode string) error {

	if len(mode) == 0 {
		mode = "Lax"
	}

	switch mode {
	case "Strict", "Lax", "None", "":
		c.sameSite = mode
		return nil
	default:
		return errors.New("invalid SameSite mode:" + mode + "(valid: Strict, Lax, None)")
	}
}
