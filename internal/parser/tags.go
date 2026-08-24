package parser

import (
	"net/url"
	"strings"
	"unicode/utf8"
)

const maxTags = 3

func extractTags(pageURL, keywords string) []string {
	tags := make([]string, 0, maxTags)
	seen := make(map[string]struct{}, maxTags)

	for _, keyword := range strings.FieldsFunc(keywords, func(r rune) bool {
		return r == ',' || r == ';'
	}) {
		keyword = strings.ToLower(strings.Trim(keyword, " \t\r\n\"'#"))
		if keyword == "" || utf8.RuneCountInString(keyword) > 32 {
			continue
		}
		if _, exists := seen[keyword]; exists {
			continue
		}
		seen[keyword] = struct{}{}
		tags = append(tags, keyword)
		if len(tags) == maxTags {
			return tags
		}
	}
	if len(tags) > 0 {
		return tags
	}

	parsedURL, err := url.Parse(pageURL)
	if err != nil {
		return nil
	}
	hostname := strings.ToLower(parsedURL.Hostname())
	if hostname == "" || utf8.RuneCountInString(hostname) > 32 {
		return nil
	}
	return []string{hostname}
}
