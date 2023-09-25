package models

type CreateCategory struct {
	Name      string `json:"name"`
	Parent_id string `json:"parent_id"`
}

type Category struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Parent_id string `json:"parent_id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type CategoryIdRequest struct {
	Id string `json:"id"`
}

type UpdateCategory struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Parent_id string `json:"parent_id"`
}
type GetAllCategoryRequest struct {
	Page  int    `json:"page"`
	Limit int    `json:"limit"`
	Name  string `json:"name"`
}

type GetAllCategoryResponse struct {
	Categories []Category `json:"category"`
	Count      int        `json:"count"`
}
