package domain

type ProductCategory struct {
	ID          string `json:"id,omitempty"`
	Title       string `json:"title" `
	Description string `json:"description"`
}
