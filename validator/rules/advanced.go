package rules

import (
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"
)

// IP validates IP addresses (v4 or v6)
type IP struct {
	// AllowV4 allows IPv4 addresses
	AllowV4 bool
	// AllowV6 allows IPv6 addresses
	AllowV6 bool
	// AllowEmpty allows empty values
	AllowEmpty bool
}

func (i IP) Validate(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("value must be a string")
	}

	if str == "" {
		if i.AllowEmpty {
			return nil
		}
		return fmt.Errorf("value is required")
	}

	ip := net.ParseIP(str)
	if ip == nil {
		return fmt.Errorf("invalid IP address format")
	}

	ipv4 := ip.To4() != nil
	if ipv4 && !i.AllowV4 {
		return fmt.Errorf("IPv4 addresses are not allowed")
	}
	if !ipv4 && !i.AllowV6 {
		return fmt.Errorf("IPv6 addresses are not allowed")
	}

	return nil
}

// Domain validates domain names
type Domain struct {
	// AllowSubdomains allows subdomains (e.g., sub.example.com)
	AllowSubdomains bool
	// AllowEmpty allows empty values
	AllowEmpty bool
}

func (d Domain) Validate(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("value must be a string")
	}

	if str == "" {
		if d.AllowEmpty {
			return nil
		}
		return fmt.Errorf("value is required")
	}

	// Domain name validation rules:
	// 1. Length between 1 and 253 characters
	// 2. Each label between 1 and 63 characters
	// 3. Only alphanumeric characters and hyphens
	// 4. Labels cannot start or end with hyphens
	// 5. TLD cannot be all numeric

	if len(str) > 253 {
		return fmt.Errorf("domain name too long")
	}

	labels := strings.Split(str, ".")
	if len(labels) < 2 {
		return fmt.Errorf("invalid domain name format")
	}

	if !d.AllowSubdomains && len(labels) > 2 {
		return fmt.Errorf("subdomains are not allowed")
	}

	for i, label := range labels {
		if len(label) == 0 {
			return fmt.Errorf("invalid domain name format")
		}
		if len(label) > 63 {
			return fmt.Errorf("domain name too long")
		}
		if !regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9-]*[a-zA-Z0-9])?$`).MatchString(label) {
			return fmt.Errorf("invalid domain name format")
		}
		if i == len(labels)-1 && regexp.MustCompile(`^[0-9]+$`).MatchString(label) {
			return fmt.Errorf("TLD cannot be all numeric")
		}
	}

	return nil
}

// Password validates password strength
type Password struct {
	// MinLength is the minimum length required
	MinLength int
	// MaxLength is the maximum length allowed
	MaxLength int
	// RequireUpper requires at least one uppercase letter
	RequireUpper bool
	// RequireLower requires at least one lowercase letter
	RequireLower bool
	// RequireDigit requires at least one digit
	RequireDigit bool
	// RequireSpecial requires at least one special character
	RequireSpecial bool
}

func (p Password) Validate(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("value must be a string")
	}

	if len(str) < p.MinLength {
		return fmt.Errorf("password must be at least %d characters", p.MinLength)
	}
	if p.MaxLength > 0 && len(str) > p.MaxLength {
		return fmt.Errorf("password must not exceed %d characters", p.MaxLength)
	}

	if p.RequireUpper && !regexp.MustCompile(`[A-Z]`).MatchString(str) {
		return fmt.Errorf("password must contain at least one uppercase letter")
	}
	if p.RequireLower && !regexp.MustCompile(`[a-z]`).MatchString(str) {
		return fmt.Errorf("password must contain at least one lowercase letter")
	}
	if p.RequireDigit && !regexp.MustCompile(`[0-9]`).MatchString(str) {
		return fmt.Errorf("password must contain at least one digit")
	}
	if p.RequireSpecial && !regexp.MustCompile(`[^a-zA-Z0-9]`).MatchString(str) {
		return fmt.Errorf("password must contain at least one special character")
	}

	return nil
}

// CreditCard validates credit card numbers using the Luhn algorithm
type CreditCard struct {
	// AllowEmpty allows empty values
	AllowEmpty bool
}

func (c CreditCard) Validate(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("value must be a string")
	}

	if str == "" {
		if c.AllowEmpty {
			return nil
		}
		return fmt.Errorf("value is required")
	}

	// Remove spaces and hyphens
	str = regexp.MustCompile(`[\s-]`).ReplaceAllString(str, "")

	if !regexp.MustCompile(`^[0-9]{13,19}$`).MatchString(str) {
		return fmt.Errorf("invalid credit card number format")
	}

	// Luhn algorithm
	var sum int
	nDigits := len(str)
	parity := nDigits % 2

	for i := 0; i < nDigits; i++ {
		digit := int(str[i] - '0')
		if i%2 == parity {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
	}

	if sum%10 != 0 {
		return fmt.Errorf("invalid credit card number format")
	}

	return nil
}

// CIDR validates IPv4 CIDR notation
type CIDR struct {
	// AllowEmpty allows empty values
	AllowEmpty bool
}

func (c CIDR) Validate(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("value must be a string")
	}

	if str == "" {
		if c.AllowEmpty {
			return nil
		}
		return fmt.Errorf("value is required")
	}

	_, _, err := net.ParseCIDR(str)
	if err != nil {
		return fmt.Errorf("invalid CIDR format")
	}

	return nil
}

// MAC validates MAC addresses
type MAC struct {
	// AllowEmpty allows empty values
	AllowEmpty bool
}

func (m MAC) Validate(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("value must be a string")
	}

	if str == "" {
		if m.AllowEmpty {
			return nil
		}
		return fmt.Errorf("value is required")
	}

	// Remove colons and hyphens
	str = regexp.MustCompile(`[:-]`).ReplaceAllString(str, "")

	if !regexp.MustCompile(`^[0-9A-Fa-f]{12}$`).MatchString(str) {
		return fmt.Errorf("invalid MAC address format")
	}

	return nil
}

// LatLong validates latitude and longitude coordinates
type LatLong struct {
	// AllowEmpty allows empty values
	AllowEmpty bool
}

func (l LatLong) Validate(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("value must be a string")
	}

	if str == "" {
		if l.AllowEmpty {
			return nil
		}
		return fmt.Errorf("value is required")
	}

	// Format: "latitude,longitude"
	parts := strings.Split(str, ",")
	if len(parts) != 2 {
		return fmt.Errorf("invalid lat/long format, expected 'latitude,longitude'")
	}

	lat, err := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
	if err != nil || lat < -90 || lat > 90 {
		return fmt.Errorf("invalid latitude value")
	}

	long, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
	if err != nil || long < -180 || long > 180 {
		return fmt.Errorf("invalid longitude value")
	}

	return nil
}

// Color validates color codes in various formats
type Color struct {
	// AllowHEX allows HEX color codes (#RRGGBB or #RGB)
	AllowHEX bool
	// AllowRGB allows RGB color codes (rgb(r,g,b))
	AllowRGB bool
	// AllowHSL allows HSL color codes (hsl(h,s%,l%))
	AllowHSL bool
	// AllowEmpty allows empty values
	AllowEmpty bool
}

func (c Color) Validate(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("value must be a string")
	}

	if str == "" {
		if c.AllowEmpty {
			return nil
		}
		return fmt.Errorf("value is required")
	}

	str = strings.TrimSpace(strings.ToLower(str))

	if c.AllowHEX {
		if regexp.MustCompile(`^#([0-9a-f]{3}|[0-9a-f]{6})$`).MatchString(str) {
			return nil
		}
	}

	if c.AllowRGB {
		if match := regexp.MustCompile(`^rgb\((\d{1,3}),\s*(\d{1,3}),\s*(\d{1,3})\)$`).FindStringSubmatch(str); match != nil {
			r, _ := strconv.Atoi(match[1])
			g, _ := strconv.Atoi(match[2])
			b, _ := strconv.Atoi(match[3])
			if r >= 0 && r <= 255 && g >= 0 && g <= 255 && b >= 0 && b <= 255 {
				return nil
			}
		}
	}

	if c.AllowHSL {
		if match := regexp.MustCompile(`^hsl\((\d{1,3}),\s*(\d{1,3})%,\s*(\d{1,3})%\)$`).FindStringSubmatch(str); match != nil {
			h, _ := strconv.Atoi(match[1])
			s, _ := strconv.Atoi(match[2])
			l, _ := strconv.Atoi(match[3])
			if h >= 0 && h <= 360 && s >= 0 && s <= 100 && l >= 0 && l <= 100 {
				return nil
			}
		}
	}

	return fmt.Errorf("invalid color format")
}

// EmailDNS validates email addresses and optionally checks DNS records
type EmailDNS struct {
	// CheckDNS enables MX record validation
	CheckDNS bool
	// AllowEmpty allows empty values
	AllowEmpty bool
}

func (e EmailDNS) Validate(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("value must be a string")
	}

	if str == "" {
		if e.AllowEmpty {
			return nil
		}
		return fmt.Errorf("value is required")
	}

	// Basic email format validation
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(str) {
		return fmt.Errorf("invalid email format")
	}

	if e.CheckDNS {
		parts := strings.Split(str, "@")
		_, err := net.LookupMX(parts[1])
		if err != nil {
			return fmt.Errorf("domain does not have valid MX records")
		}
	}

	return nil
}

// Hostname validates hostnames according to RFC 1123
type Hostname struct {
	// AllowWildcard allows wildcard in hostname (e.g., *.example.com)
	AllowWildcard bool
	// AllowEmpty allows empty values
	AllowEmpty bool
}

func (h Hostname) Validate(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("value must be a string")
	}

	if str == "" {
		if h.AllowEmpty {
			return nil
		}
		return fmt.Errorf("value is required")
	}

	if h.AllowWildcard && strings.HasPrefix(str, "*.") {
		str = "host" + str[1:]
	}

	// RFC 1123 hostname validation
	if len(str) > 255 {
		return fmt.Errorf("hostname too long")
	}

	hostnameRegex := regexp.MustCompile(`^[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$`)
	if !hostnameRegex.MatchString(str) {
		return fmt.Errorf("invalid hostname format")
	}

	return nil
}

// Port validates port numbers
type Port struct {
	// Min is the minimum allowed port number (default: 1)
	Min int
	// Max is the maximum allowed port number (default: 65535)
	Max int
	// AllowPrivileged allows ports below 1024
	AllowPrivileged bool
	// AllowEmpty allows empty values
	AllowEmpty bool
}

func (p Port) Validate(value interface{}) error {
	if p.Min == 0 {
		p.Min = 1
	}
	if p.Max == 0 {
		p.Max = 65535
	}
	if !p.AllowPrivileged && p.Min < 1024 {
		p.Min = 1024
	}

	// Handle string input
	if str, ok := value.(string); ok {
		if str == "" {
			if p.AllowEmpty {
				return nil
			}
			return fmt.Errorf("value is required")
		}
		port, err := strconv.Atoi(str)
		if err != nil {
			return fmt.Errorf("invalid port number")
		}
		value = port
	}

	// Handle numeric input
	port, ok := value.(int)
	if !ok {
		return fmt.Errorf("value must be a string or integer")
	}

	if port < p.Min || port > p.Max {
		return fmt.Errorf("port must be between %d and %d", p.Min, p.Max)
	}

	return nil
}

// SemVer validates semantic version strings
type SemVer struct {
	// AllowPrefix allows 'v' prefix (e.g., v1.0.0)
	AllowPrefix bool
	// RequirePrefix requires 'v' prefix when AllowPrefix is true
	RequirePrefix bool
	// AllowPrerelease allows prerelease versions (e.g., 1.0.0-alpha)
	AllowPrerelease bool
	// AllowBuild allows build metadata (e.g., 1.0.0+001)
	AllowBuild bool
	// AllowEmpty allows empty values
	AllowEmpty bool
}

func (s SemVer) Validate(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("value must be a string")
	}

	if str == "" {
		if s.AllowEmpty {
			return nil
		}
		return fmt.Errorf("value is required")
	}

	// Handle v prefix
	if strings.HasPrefix(str, "v") {
		if !s.AllowPrefix {
			return fmt.Errorf("v prefix not allowed")
		}
		str = str[1:]
	} else if s.RequirePrefix && s.AllowPrefix {
		return fmt.Errorf("v prefix is required")
	}

	// Split version into parts
	parts := strings.SplitN(str, "+", 2)
	versionParts := strings.SplitN(parts[0], "-", 2)

	// Validate core version (X.Y.Z)
	core := strings.Split(versionParts[0], ".")
	if len(core) != 3 {
		return fmt.Errorf("version must be in format X.Y.Z")
	}

	for _, num := range core {
		if !regexp.MustCompile(`^\d+$`).MatchString(num) {
			return fmt.Errorf("version components must be numeric")
		}
	}

	// Validate prerelease
	if len(versionParts) > 1 {
		if !s.AllowPrerelease {
			return fmt.Errorf("prerelease versions not allowed")
		}
		if !regexp.MustCompile(`^[0-9A-Za-z-]+(\.[0-9A-Za-z-]+)*$`).MatchString(versionParts[1]) {
			return fmt.Errorf("invalid prerelease format")
		}
	}

	// Validate build metadata
	if len(parts) > 1 {
		if !s.AllowBuild {
			return fmt.Errorf("build metadata not allowed")
		}
		if !regexp.MustCompile(`^[0-9A-Za-z-]+(\.[0-9A-Za-z-]+)*$`).MatchString(parts[1]) {
			return fmt.Errorf("invalid build metadata format")
		}
	}

	return nil
}
