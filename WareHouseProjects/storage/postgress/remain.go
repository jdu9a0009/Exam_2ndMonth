package postgres

import (
	"WareHouseProjects/models"
	"WareHouseProjects/pkg/helper"
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type remainRepo struct {
	db *pgxpool.Pool
}

func NewRemainRepo(db *pgxpool.Pool) *remainRepo {
	return &remainRepo{
		db: db,
	}
}

func (c *remainRepo) CreateRemain(req *models.CreateRemain) (string, error) {
	var (
		id = uuid.NewString()
	)

	query := `
		INSERT INTO "remain"(
			"id", 
			"branch_id",
			"category_id",
			"name",
			"price",
			"barcode",
			"count",
			"total_price",
			"created_at" )
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW())`

	_, err := c.db.Exec(context.Background(), query,
		id,
		req.Branch_id,
		req.Category_id,
		req.Name,
		req.Price,
		req.Barcode,
		req.Count,
		req.TotalPrice,
	)

	if err != nil {
		return "", err
	}

	return id, nil
}

func (c *remainRepo) GetRemain(req *models.RemainIdRequest) (*models.Remain, error) {
	query := `
		SELECT
		    "id",
		    "branch_id",
		    "category_id",
		    "name",
		    "price",
		    "barcode",
		    "count",
		    "total_price",
		    "created_at",
			   "updated_at"
		FROM "remain"
		WHERE id = $1
	`
	var (
		createdAt  time.Time
		updatedAt  sql.NullTime
		totalPrice float64
	)

	rem := models.Remain{}
	err := c.db.QueryRow(context.Background(), query, req.Id).Scan(
		&rem.ID,
		&rem.Branch_id,
		&rem.Category_id,
		&rem.Name,
		&rem.Price,
		&rem.Barcode,
		&rem.Count,
		&totalPrice,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("Remain not found")
	}
	rem.TotalPrice = totalPrice
	rem.CreatedAt = createdAt.Format(time.RFC3339)
	if updatedAt.Valid {
		rem.UpdatedAt = updatedAt.Time.Format(time.RFC3339)
	}

	return &rem, nil
}

func (c *remainRepo) GetAllRemain(req *models.GetAllRemainRequest) (*models.GetAllRemainResponse, error) {
	params := make(map[string]interface{})
	resp := &models.GetAllRemainResponse{}

	resp.Remainings = make([]models.Remain, 0)

	filter := " WHERE true "
	query := `
		SELECT
			COUNT(*) OVER(),
			"id",
			"branch_id",
			"category_id",
			"name",
			"price",
			"barcode",
			"count",
			"total_price",
			"created_at",
			"updated_at"
		FROM "remain"
	`
	if req.Category_id != "" {
		filter += ` AND ("category_id" ILIKE '%' || :search ) `
		params["search"] = req.Category_id
	}

	if req.Branch_id != "" {
		filter += ` AND ("branch_id" ILIKE '%' || :search ) `
		params["search"] = req.Branch_id
	}

	if req.Barcode != "" {
		filter += ` AND ("barcode" ILIKE '%' || :search ) `
		params["search"] = req.Barcode
	}

	offset := (req.Page - 1) * req.Limit
	params["limit"] = req.Limit
	params["offset"] = offset

	query = query + filter + " ORDER BY created_at DESC OFFSET :offset LIMIT :limit "
	rquery, pArr := helper.ReplaceQueryParams(query, params)

	rows, err := c.db.Query(context.Background(), rquery, pArr...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var (
		totalPrice float64
		createdAt  time.Time
		updatedAt  sql.NullTime
	)

	count := 0
	for rows.Next() {
		var rem models.Remain
		count++
		err := rows.Scan(
			&count,
			&rem.ID,
			&rem.Branch_id,
			&rem.Category_id,
			&rem.Name,
			&rem.Price,
			&rem.Barcode,
			&rem.Count,
			&totalPrice,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, err
		}
		rem.TotalPrice = totalPrice
		rem.CreatedAt = createdAt.Format(time.RFC3339)
		if updatedAt.Valid {
			rem.UpdatedAt = updatedAt.Time.Format(time.RFC3339)
		}
		resp.Remainings = append(resp.Remainings, rem)
	}

	resp.Count = count
	return resp, nil
}
func (c *remainRepo) UpdateRemain(req *models.UpdateRemain) (string, error) {
	totalPrice := req.Count * req.Price

	query := `UPDATE remain 
	            SET  branch_id = $1, 
				     category_id = $2,
					 name=$3,
					 price=$4,
					 barcode=$5,
					 count=$6,
					 total_price=$7, 
					 updated_at = NOW() 
					 WHERE id = $8 RETURNING id`

	result, err := c.db.Exec(context.Background(), query, req.Branch_id, req.Category_id, req.Name, req.Price, req.Barcode, req.Count, totalPrice, req.ID)
	if err != nil {
		return "Error Update Remain", err
	}

	if result.RowsAffected() == 0 {
		return "", fmt.Errorf("Remain not found")
	}

	return req.ID, nil
}

func (c *remainRepo) DeleteRemain(req *models.RemainIdRequest) (resp string, err error) {
	query := `DELETE FROM remain 
	            WHERE id = $1 RETURNING id`

	result, err := c.db.Exec(context.Background(), query, req.Id)
	if err != nil {
		return "Error from Delete Remain", err
	}

	if result.RowsAffected() == 0 {
		return "", fmt.Errorf("Remain not found")
	}

	return req.Id, nil
}

func (c *remainRepo) CheckRemain(req *models.CheckRemain) (string, error) {
	var id sql.NullString
	var params map[string]interface{}

	query := `
		SELECT
			"id"
		FROM "remaining"
		WHERE "branch_id" = :branch_id AND "barcode" = :barcode
	`

	params = map[string]interface{}{
		"branch_id": req.Branch_id,
		"barcode":   req.Barcode,
	}
	queryN, args := helper.ReplaceQueryParams(query, params)

	err := c.db.QueryRow(context.Background(), queryN, args...).Scan(
		&id,
	)
	if err != nil {
		return id.String, err
	}

	return id.String, nil
}

func (c *remainRepo) UpdateIdAviable(req *models.UpdateRemain) (string, error) {
	query := `UPDATE  remain SET
	                 "branch_id" = $1,
	                 "category_id" = $2,
	                 "name" = $3,
	                 "price" = $4,
	                 "barcode" =$5,
	                 "count" = "count" + $6,
                	 "total_price" = "total_price" + $7,
	                 "updated_at" = NOW()
                    WHERE id = $8    `

	resp, err := c.db.Exec(context.Background(), query,
		req.Branch_id,
		req.Category_id,
		req.Name,
		req.Price,
		req.Barcode,
		req.Count,
		req.TotalPrice,
		req.ID,
	)
	if err != nil {
		return "", err
	}

	if resp.RowsAffected() == 0 {
		return "", fmt.Errorf("remaining with ID %s not found", req.ID)
	}

	return req.ID, nil
}

func (c *coming_TableProductRepo) GetComingTableById(req *models.ComingTableProductIdRequest) (*models.ComingTableProduct, error) {
	query := `
	SELECT
		"id",
		"category_id",
		"name",
		"price",
		"barcode",
		sum("count"),
		sum("total_price"),
	FROM "remain"
	WHERE coming_table_id = $1
	GROUPD BY id,barcode
`

	rem := models.ComingTableProduct{}
	err := c.db.QueryRow(context.Background(), query, req.Id).Scan(
		&rem.ID,
		&rem.Category_id,
		&rem.Name,
		&rem.Price,
		&rem.Barcode,
		&rem.Count,
		&rem.TotalPrice,
	)
	if err != nil {
		return nil, fmt.Errorf(" Error in GetComingTableById func")
	}

	return &rem, nil
}
