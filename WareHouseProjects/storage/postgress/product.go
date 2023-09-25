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

type productRepo struct {
	db *pgxpool.Pool
}

func NewProductRepo(db *pgxpool.Pool) *productRepo {
	return &productRepo{
		db: db,
	}
}

func (r *productRepo) CreateProduct(req *models.CreateProduct) (string, error) {
	var (
		id = uuid.NewString()
	)

	query := `
				INSERT INTO "product"(
					"id",
					"name",
					"price",
					"barcode",
					"category_id",
					"created_at")
				VALUES ($1, $2, $3, $4, $5, NOW())`

	_, err := r.db.Exec(context.Background(), query,
		id,
		req.Name,
		req.Price,
		req.Barcode,
		req.Category_id,
	)

	if err != nil {
		return "", err
	}

	return id, nil

}

func (c *productRepo) GetProduct(req *models.ProductIdRequest) (resp *models.Product, err error) {

	query := `
		SELECT
		   "id", 
		    "name",
		    "price",
			"barcode",
			"category_id",
		    "created_at", 
			"updated_at" 
		FROM "product"
		WHERE id = $1
	`
	var (
		createdAt time.Time
		updatedAt sql.NullTime
	)

	Product := models.Product{}
	err = c.db.QueryRow(context.Background(), query, req.Id).Scan(
		&Product.ID,
		&Product.Name,
		&Product.Price,
		&Product.Barcode,
		&Product.Category_id,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("Product not found")
	}
	Product.CreatedAt = createdAt.Format(time.RFC3339)
	if updatedAt.Valid {
		Product.UpdatedAt = updatedAt.Time.Format(time.RFC3339)
	}

	return &Product, nil
}

func (c *productRepo) GetProductByBarcode(req *models.CheckBarcodeComingTable) (resp *models.RespBarcodeProduct, err error) {

	query := `
		SELECT
		    "name",
		    "price",
			"category_id"
		FROM "product"
		WHERE barcode = $1 and coming_table_id=$2
	`

	Product := models.RespBarcodeProduct{}
	err = c.db.QueryRow(context.Background(), query, req.Barcode, req.Coming_Table_id).Scan(
		&Product.Name,
		&Product.Price,
		&Product.Category_id,
	)
	if err != nil {
		return nil, fmt.Errorf("Product not found")
	}

	return &Product, nil
}

func (c *productRepo) GetAllProduct(req *models.GetAllProductRequest) (*models.GetAllProductResponse, error) {
	params := make(map[string]interface{})
	var resp = &models.GetAllProductResponse{}

	resp.Products = make([]models.Product, 0)

	filter := " WHERE true "
	query := `
		SELECT
			COUNT(*) OVER(),
			"id",
			"name",
			"price",		
			"barcode",
			"category_id",
			"created_at",
			"updated_at" 
		FROM "product"
	`
	if req.Name != "" {
		filter += ` AND ("name" ILIKE '%' || :name ) `
		params["search"] = req.Name
	}
	if req.Barcode != "" {
		filter += ` AND ("barcode" ILIKE '%' || :barcode ) `
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
			id          sql.NullString
			name        sql.NullString
			price       sql.NullFloat64
			barcode     sql.NullString
			category_id sql.NullString
			createdAt   sql.NullString
			updatedAt   sql.NullString
		)
		err := rows.Scan(
			&resp.Count,
			&id,
			&name,
			&price,
			&barcode,
			&category_id,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, err
		}
		resp.Products = append(resp.Products, models.Product{
			ID:          id.String,
			Name:        name.String,
			Price:       price.Float64,
			Barcode:     barcode.String,
			Category_id: category_id.String,
			CreatedAt:   createdAt.String,
			UpdatedAt:   updatedAt.String,
		})
	}
	return resp, nil
}

func (c *productRepo) UpdateProduct(req *models.UpdateProduct) (string, error) {

	query := `
		UPDATE
			"product"
		SET
			"name" = $1,
			"price" = $2,
			"barcode" = $3,
			"category_id" = $4,
			"updated_at" = NOW()
			WHERE id= $5 RETURNING id	`

	result, err := c.db.Exec(context.Background(), query, req.Name, req.Price, req.Barcode, req.Category_id, req.ID)
	if err != nil {
		return "Error Update Product", err
	}

	if result.RowsAffected() == 0 {
		return "", fmt.Errorf("Product not found")
	}

	return req.ID, nil
}

func (c *productRepo) DeleteProduct(req *models.ProductIdRequest) (resp string, err error) {
	query := `DELETE FROM product 
	            WHERE id = $1 RETURNING id`

	result, err := c.db.Exec(context.Background(), query, req.Id)
	if err != nil {
		return "Error from Delete Product", err
	}

	if result.RowsAffected() == 0 {
		return "", fmt.Errorf("Product not found")
	}

	return req.Id, nil
}
