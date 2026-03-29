package domain

type Product struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"category"`
	CategoryID  string `json:"category_id"`
	Price       int    `json:"price" `
	Brand       string `json:"brand"`
	Quantity    int    `json:"quantity"`
}
