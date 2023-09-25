package postgres

import (
	"WareHouseProjects/models"
	"WareHouseProjects/pkg/helper"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type coming_TableProductRepo struct {
	db *pgxpool.Pool
}

func NewComingTableProductRepo(db *pgxpool.Pool) *coming_TableProductRepo {
	return &coming_TableProductRepo{
		db: db,
	}
}

func (r *coming_TableProductRepo) CreateComingTableProduct(req *models.CreateComingTableProduct) (string, error) {
	var (
		id    = uuid.NewString()
		query string
	)

	query = `
		INSERT INTO "coming_table_product"(
			"id", 
			"category_id",
			"name",
			"price",
			"barcode",
			"count",
			"total_price",
			"coming_table_id",
			"created_at" )
		VALUES ($1, $2, $3, $4, $5,$6,$7,$8 NOW())`

	_, err := r.db.Exec(context.Background(), query,
		id,
		req.Category_id,
		req.Name,
		req.Price,
		req.Barcode,
		req.Count,
		req.TotalPrice,
		req.Coming_Table_id,
	)

	if err != nil {
		return "", err
	}

	return id, nil
}

func (c *coming_TableProductRepo) GetComingTableProduct(req *models.ComingTableProductIdRequest) (resp *models.ComingTableProduct, err error) {

	query := `
		SELECT
		    "id", 
			"category_id",
			"name",
			"price",
			"barcode",
			"count",
			"total_price",
			"coming_table_id",
		    "created_at",
			"updated_at" 
		FROM "coming_table_product"
		WHERE id = $1
	`
	var (
		createdAt time.Time
		updatedAt sql.NullTime
	)

	ComingTableProduct := models.ComingTableProduct{}
	err = c.db.QueryRow(context.Background(), query, req.Id).Scan(
		&ComingTableProduct.ID,
		&ComingTableProduct.Category_id,
		&ComingTableProduct.Name,
		&ComingTableProduct.Price,
		&ComingTableProduct.Barcode,
		&ComingTableProduct.Count,
		&ComingTableProduct.TotalPrice,
		&ComingTableProduct.Coming_Table_id,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf(" ComingTableProduct not found")
	}
	ComingTableProduct.CreatedAt = createdAt.Format(time.RFC3339)
	if updatedAt.Valid {
		ComingTableProduct.UpdatedAt = updatedAt.Time.Format(time.RFC3339)
	}

	return &ComingTableProduct, nil
}

func (c *coming_TableProductRepo) GetAllComingTableProduct(req *models.GetAllComingTableProductRequest) (*models.GetAllComingTableProductResponse, error) {
	params := make(map[string]interface{})
	var resp = &models.GetAllComingTableProductResponse{}

	resp.ComingTableProducts = make([]models.ComingTableProduct, 0)

	filter := " WHERE true "
	query := `
			SELECT
				COUNT(*) OVER(),
				"id", 
				"category_id",
				"name",
				"price",
				"barcode",
				"count",
				"total_price",
				"coming_table_id",
				"created_at",
				"updated_at" 
			FROM "coming_table_product"
		`
	if req.Category_id != "" {
		filter += ` AND ("category_id" ILIKE '%' || :category_id) `
		params["search"] = req.Category_id
	}
	if req.Barcode != "" {
		filter += ` AND ("barcode" ILIKE '%' || :barcode) `
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

	for rows.Next() {
		var (
			id              sql.NullString
			category_id     sql.NullString
			name            sql.NullString
			price           sql.NullFloat64
			barcode         sql.NullString
			count           sql.NullFloat64
			total_price     sql.NullFloat64
			coming_table_id sql.NullString
			createdAt       sql.NullString
			updatedAt       sql.NullString
		)
		err := rows.Scan(
			&resp.Count,
			&id,
			&category_id,
			&name,
			&price,
			&barcode,
			&count,
			&total_price,
			&coming_table_id,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, err
		}
		resp.ComingTableProducts = append(resp.ComingTableProducts, models.ComingTableProduct{
			ID:              id.String,
			Category_id:     category_id.String,
			Name:            name.String,
			Price:           price.Float64,
			Barcode:         barcode.String,
			Count:           count.Float64,
			TotalPrice:      total_price.Float64,
			Coming_Table_id: coming_table_id.String,
			CreatedAt:       createdAt.String,
			UpdatedAt:       updatedAt.String,
		})
	}
	return resp, nil
}

func (c *coming_TableProductRepo) UpdateComingTableProduct(req *models.UpdateComingTableProduct) (string, error) {
	total_price := req.Count * req.Price

	query := `UPDATE coming_table_product 
	            SET  category_id = $1, 
				     name = $2, 
					 price=$3,
					 barcode=$4,
					 count=$5,
					 total_price=$6,
					 comint_table_product=$7,
					 updated_at = NOW() 
					 WHERE id = $7 RETURNING id`

	result, err := c.db.Exec(context.Background(), query, req.Category_id, req.Name, req.Price, req.Barcode, req.Count, total_price, req.Coming_Table_id, req.ID)
	if err != nil {
		return "Error Update Coming_TableProduct", err
	}

	if result.RowsAffected() == 0 {
		return "", fmt.Errorf("Coming_TableProduct not found")
	}

	return req.ID, nil
}

func (c *coming_TableProductRepo) DeleteComingTableProduct(req *models.ComingTableProductIdRequest) (resp string, err error) {
	query := `DELETE FROM coming_TableProduct 
	            WHERE id = $1 RETURNING id`

	result, err := c.db.Exec(context.Background(), query, req.Id)
	if err != nil {
		return "Error from Delete coming_table_product", err
	}

	if result.RowsAffected() == 0 {
		return "", fmt.Errorf("Coming_TableProduct not found")
	}

	return req.Id, nil
}

func (c *coming_TableProductRepo) CheckAviableProduct(req *models.CheckBarcodeComingTable) (string, error) {
	var id sql.NullString

	query := `Select
	             id
			from coming_table_product
			where barcode=$1 and coming_table_id=$2 `

	err := c.db.QueryRow(context.Background(), query, req.Barcode, req.Coming_Table_id).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.New("not found")
		}
		return "", err
	}

	return id.String, nil
}

func (c *coming_TableProductRepo) UpdateIdAviable(req *models.UpdateComingTableProduct) (string, error) {
	query := `Update coming_table_product Set
	           category_id=$1,
			   barcode=$2,
			   name=$3,
			   price=$4,
			   count=count+$5,
			   total_price=total_price+$6,
			   coming_table_id=$7,
			   updated_at=now()
			   where id = $8  `

	result, err := c.db.Exec(context.Background(), query,
		req.Category_id,
		req.Barcode,
		req.Name,
		req.Price,
		req.Count,
		req.TotalPrice,
		req.ID,
	)
	if err != nil {
		return "", err
	}
	if result.RowsAffected() == 0 {
		return "", fmt.Errorf("Not found this id in coming table product:  %s", req.Coming_Table_id)
	}
	return req.Coming_Table_id, nil

}
