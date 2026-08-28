package resource

import (
	"strings"

	"nilchan-hackaton/internal/parser"
)

const (
	minContentWords = 300
	maxContentChars = 50_000
)

func validatePage(page parser.Page) (parser.Page, error) {
	if len(strings.Fields(page.Content)) < minContentWords {
		return parser.Page{}, ErrContentTooShort
	}
	if len([]rune(page.Content)) > maxContentChars {
		return parser.Page{}, ErrContentTooLong
	}
	return page, nil
}
