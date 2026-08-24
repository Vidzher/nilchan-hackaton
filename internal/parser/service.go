package parser

import "context"

type Page struct {
	Content  string
	Title    string
	Keywords string
}

type pageParser interface {
	ParsePage(ctx context.Context, pageURL string) (Page, error)
}

type Service struct{ parser pageParser }

func NewService(parser pageParser) *Service { return &Service{parser: parser} }

func (ps *Service) ParsePage(ctx context.Context, pageURL string) (Page, error) {
	return ps.parser.ParsePage(ctx, pageURL)
}
