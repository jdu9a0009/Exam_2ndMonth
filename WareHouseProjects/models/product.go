package models

type CreateProduct struct {
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Barcode     string  `json:"barcode"`
	Category_id string  `json:"category_id"`
}

type Product struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Barcode     string  `json:"barcode"`
	Category_id string  `json:"category_id"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}
type UpdateProduct struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Barcode     string  `json:"barcode"`
	Category_id string  `json:"category_id"`
}

type RespBarcodeProduct struct {
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Category_id string  `json:"category_id"`
}

type ProductIdRequest struct {
	Id string `json:"id"`
}

type GetAllProductRequest struct {
	Page    int    `json:"page"`
	Limit   int    `json:"limit"`
	Name    string `json:"name"`
	Barcode string `json:"barcode"`
}

type GetAllProductResponse struct {
	Products []Product `json:"product"`
	Count    int       `json:"count"`
}
