package models

type CreateRemain struct {
	Branch_id   string  `json:"branch_id"`
	Category_id string  `json:"category_id"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Barcode     string  `json:"barcode"`
	Count       float64 `json:"count"`
	TotalPrice  float64 `json:"total_price"`
}

type CheckRemain struct {
	Barcode   string `json:"barcode"`
	Branch_id string `json:"branch_id"`
}

type Remain struct {
	ID          string  `json:"id"`
	Branch_id   string  `json:"branch_id"`
	Category_id string  `json:"category_id"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Barcode     string  `json:"barcode"`
	Count       float64 `json:"count"`
	TotalPrice  float64 `json:"total_price"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

type RemainIdRequest struct {
	Id string `json:"id"`
}

type UpdateRemain struct {
	ID          string  `json:"id"`
	Branch_id   string  `json:"branch_id"`
	Category_id string  `json:"category_id"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Barcode     string  `json:"barcode"`
	Count       float64 `json:"count"`
	TotalPrice  float64 `json:"total_price"`
}

type GetAllRemainRequest struct {
	Page        int    `json:"page"`
	Limit       int    `json:"limit"`
	Branch_id   string `json:"branch_id"`
	Category_id string `json:"category_id"`
	Barcode     string `json:"barcode"`
}

type GetAllRemainResponse struct {
	Remainings []Remain `json:"remaining"`
	Count      int      `json:"count"`
}
