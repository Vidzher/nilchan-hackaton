package pparser

import "context"

type Parser interface {
	ParsePage(ctx context.Context, pageURL string) (string, error)
}

type ParserService struct{ parser Parser }

func NewParserService(parser Parser) *ParserService { return &ParserService{parser: parser} }

func (ps *ParserService) ParsePage(ctx context.Context, pageURL string) (string, error) {
	return ps.parser.ParsePage(ctx, pageURL)
}
