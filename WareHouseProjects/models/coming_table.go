package models

type TableType string

const (
	Finishied TableType = "finishied"
	InProcess TableType = "in_process"
)

type CreateComingTable struct {
	Coming_id string `json:"coming_id"`
	Branch_id string `json:"branch_id"`
	DateTime  string `json:"date_time"`
}

type ComingTable struct {
	ID        string    `json:"id"`
	ComingID  string    `json:"coming_id"`
	BranchID  string    `json:"branch_id"`
	DateTime  string    `json:"date_time"`
	Status    TableType `json:"status"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}
type UpdateComingTable struct {
	ID       string `json:"id"`
	ComingID string `json:"coming_id"`
	BranchID string `json:"branch_id"`
	DateTime string `json:"date_time"`
}

type ComingTableIdRequest struct {
	Id string `json:"id"`
}

type GetAllComingTableRequest struct {
	Page     int    `json:"page"`
	Limit    int    `json:"limit"`
	ComingID string `json:"coming_id"`
	BranchID string `json:"branch_id"`
}

type GetAllComingTableResponse struct {
	ComingTables []ComingTable `json:"coming_table"`
	Count        int           `json:"count"`
}
