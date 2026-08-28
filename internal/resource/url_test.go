package resource

import (
	"errors"
	"strings"
	"testing"
)

func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		want string
	}{
		{
			name: "normalizes host port fragment and tracking parameters",
			raw:  " HTTPS://Example.COM.:443/article?utm_source=test&b=2&a=1&FBCLID=id#section ",
			want: "https://example.com/article?a=1&b=2",
		},
		{
			name: "keeps non-default port",
			raw:  "http://Example.com:8080/article",
			want: "http://example.com:8080/article",
		},
		{
			name: "removes default HTTP port and gclid",
			raw:  "http://example.com:80/?gclid=id&q=go",
			want: "http://example.com/?q=go",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := NormalizeURL(test.raw)
			if err != nil {
				t.Fatalf("NormalizeURL() error = %v", err)
			}
			if got != test.want {
				t.Errorf("NormalizeURL() = %q, want %q", got, test.want)
			}
		})
	}
}

func TestNormalizeURLRejectsInvalidAndBlockedURLs(t *testing.T) {
	tests := []struct {
		name    string
		raw     string
		wantErr error
	}{
		{name: "relative URL", raw: "/article", wantErr: ErrInvalidURL},
		{name: "unsupported scheme", raw: "ftp://example.com/article", wantErr: ErrInvalidURL},
		{name: "missing host", raw: "https:///article", wantErr: ErrInvalidURL},
		{name: "credentials", raw: "https://user:password@example.com", wantErr: ErrBlockedURL},
		{name: "localhost", raw: "http://localhost/article", wantErr: ErrBlockedURL},
		{name: "single label hostname", raw: "http://intranet/article", wantErr: ErrBlockedURL},
		{name: "private IPv4", raw: "http://192.168.1.1/article", wantErr: ErrBlockedURL},
		{name: "loopback IPv4", raw: "http://127.0.0.1/article", wantErr: ErrBlockedURL},
		{name: "private IPv6", raw: "http://[fd00::1]/article", wantErr: ErrBlockedURL},
		{name: "too long", raw: "https://example.com/" + strings.Repeat("a", MaxURLLength), wantErr: ErrInvalidURL},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := NormalizeURL(test.raw)
			if !errors.Is(err, test.wantErr) {
				t.Errorf("NormalizeURL() error = %v, want %v", err, test.wantErr)
			}
		})
	}
}
