package pparser

type ParsePageRequestDTO struct {
	URL string `json:"url" validate:"required,url"`
}

type ParsePageResponseDTO struct {
	URL      string `json:"url"`
	Markdown string `json:"markdown"`
}
