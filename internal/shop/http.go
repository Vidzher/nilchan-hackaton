package shop

type PurchaseRequest struct {
	ItemID string `json:"itemId" validate:"required"`
}

type PurchaseResponse struct {
	Item    CatalogItemResponse `json:"item"`
	EPoints int64               `json:"ePoints"`
}

type CatalogItemResponse struct {
	ID           string               `json:"id"`
	Type         string               `json:"type"`
	DisplayName  string               `json:"displayName"`
	Price        int64                `json:"price"`
	Presentation PresentationResponse `json:"presentation"`
}

type PresentationResponse struct {
	Emoji    string `json:"emoji,omitempty"`
	AssetKey string `json:"assetKey,omitempty"`
	CSSClass string `json:"cssClass,omitempty"`
}
