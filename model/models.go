package model

// Product Model
type Product struct{
	ID int `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
	Category string `json:"category"`
	Price float64 `json:"price"`
	Brand string `json:"brand"`
	Quantity int `json:"quantity"`
}

var ProductMap map[int]Product

// Using map as a temporary database
func InitializeInventory() {
    ProductMap = make(map[int]Product)
    ProductMap[1] = Product{1,"Pearl Necklace","A pearl necklace of standard size for any age group.","Accessories",49.99,"La Vie",10}
	ProductMap[2] = Product{2, "Gaming Mouse", "Ergonomic RGB mouse with 16000 DPI sensor.", "Electronics", 59.99, "Logitech", 25}
    ProductMap[3] = Product{3, "Leather Journal", "Hand-bound A5 notebook with recycled paper.", "Stationery", 24.50, "Moleskine", 50}
    ProductMap[4] = Product{4, "Espresso Machine", "Semi-automatic machine with milk frothing wand.", "Appliances", 299.00, "Breville", 5}
    ProductMap[5] = Product{5, "Yoga Mat", "Non-slip 6mm eco-friendly rubber mat.", "Fitness", 35.00, "Items", 40}
    ProductMap[6] = Product{6, "Wireless Earbuds", "Noise-canceling buds with 24-hour battery life.", "Electronics", 129.99, "Sony", 15}
}