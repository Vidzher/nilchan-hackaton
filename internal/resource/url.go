package resource

import (
	"net"
	"net/netip"
	"net/url"
	"strings"
)

const MaxURLLength = 2048

func NormalizeURL(rawURL string) (string, error) {
	rawURL = strings.TrimSpace(rawURL)
	if rawURL == "" || len(rawURL) > MaxURLLength {
		return "", ErrInvalidURL
	}

	parsed, err := url.Parse(rawURL)
	if err != nil {
		return "", ErrInvalidURL
	}

	parsed.Scheme = strings.ToLower(parsed.Scheme)
	if !parsed.IsAbs() || parsed.Host == "" || (parsed.Scheme != "http" && parsed.Scheme != "https") {
		return "", ErrInvalidURL
	}
	if parsed.User != nil {
		return "", ErrBlockedURL
	}

	hostname := strings.ToLower(strings.TrimSuffix(parsed.Hostname(), "."))
	if isInternalHost(hostname) {
		return "", ErrBlockedURL
	}

	port := parsed.Port()
	if (parsed.Scheme == "http" && port == "80") || (parsed.Scheme == "https" && port == "443") {
		port = ""
	}
	if port != "" {
		parsed.Host = net.JoinHostPort(hostname, port)
	} else if address, err := netip.ParseAddr(hostname); err == nil && address.Is6() {
		parsed.Host = "[" + hostname + "]"
	} else {
		parsed.Host = hostname
	}

	query := parsed.Query()
	for key := range query {
		keyLower := strings.ToLower(key)
		if strings.HasPrefix(keyLower, "utm_") || keyLower == "fbclid" || keyLower == "gclid" {
			query.Del(key)
		}
	}

	parsed.RawQuery = query.Encode()
	parsed.Fragment = ""
	return parsed.String(), nil
}

func isInternalHost(hostname string) bool {
	if hostname == "localhost" || strings.HasSuffix(hostname, ".localhost") || !strings.Contains(hostname, ".") {
		return true
	}

	address, err := netip.ParseAddr(hostname)
	if err != nil {
		return false
	}
	return !address.IsGlobalUnicast() || address.IsPrivate()
}
