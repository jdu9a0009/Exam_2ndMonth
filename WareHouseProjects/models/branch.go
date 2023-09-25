package models

type CreateBranch struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Phone   string `json:"phone"`
}

type Branch struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Address   string `json:"address"`
	Phone     string `json:"phone"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
type UpdateBranch struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
	Phone   string `json:"phone"`
}

type BranchIdRequest struct {
	Id string `json:"id"`
}

type GetAllBranchRequest struct {
	Page  int    `json:"page"`
	Limit int    `json:"limit"`
	Name  string `json:"name"`
}

type GetAllBranchResponse struct {
	Branches []Branch `json:"branches"`
	Count    int      `json:"count"`
}
