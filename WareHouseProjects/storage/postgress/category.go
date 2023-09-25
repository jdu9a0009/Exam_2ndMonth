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

type categoryRepo struct {
	db *pgxpool.Pool
}

func NewCategoryRepo(db *pgxpool.Pool) *categoryRepo {
	return &categoryRepo{
		db: db,
	}
}

func (r *categoryRepo) CreateCategory(req *models.CreateCategory) (string, error) {
	var (
		id = uuid.NewString()
	)

	query := `
		  INSERT INTO "category"(
			"id",
			"name",
			"created_at")
		  VALUES ($1, $2, NOW())`

	if req.Parent_id != "" {
		query = `
		  INSERT INTO "category"(
			"id",
			"name",
			"parent_id",
			"created_at")
		  VALUES ($1, $2, $3, NOW())`
	}

	if req.Parent_id != "" {
		_, err := r.db.Exec(context.Background(), query,
			id,
			req.Name,
			req.Parent_id,
		)

		if err != nil {
			return "", err
		}
	} else {
		_, err := r.db.Exec(context.Background(), query,
			id,
			req.Name,
		)

		if err != nil {
			return "", err
		}
	}

	return id, nil
}

func (c *categoryRepo) GetCategory(req *models.CategoryIdRequest) (resp *models.Category, err error) {

	query := `
		SELECT
		   "id", 
		    "name",
		    "parent_id",
		    "created_at", 
			"updated_at" 
		FROM "category"
		WHERE id = $1
	`
	var (
		createdAt time.Time
		updatedAt sql.NullTime
	)

	category := models.Category{}
	err = c.db.QueryRow(context.Background(), query, req.Id).Scan(
		&category.ID,
		&category.Name,
		&category.Parent_id,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("category not found")
	}
	category.CreatedAt = createdAt.Format(time.RFC3339)
	if updatedAt.Valid {
		category.UpdatedAt = updatedAt.Time.Format(time.RFC3339)
	}

	return &category, nil
}

func (c *categoryRepo) GetAllCategory(req *models.GetAllCategoryRequest) (*models.GetAllCategoryResponse, error) {
	params := make(map[string]interface{})
	var resp = &models.GetAllCategoryResponse{}

	resp.Categories = make([]models.Category, 0)

	filter := " WHERE true "
	query := `
			SELECT
				COUNT(*) OVER(),
				"id", 
				"name",
				"parent_id",
				"created_at",
				"updated_at" 
			FROM "category"
		`
	if req.Name != "" {
		filter += ` AND "name" ILIKE '%' || :search || '%' `
		params["search"] = req.Name
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
			id        sql.NullString
			name      sql.NullString
			parent_id sql.NullString
			createdAt sql.NullString
			updatedAt sql.NullString
		)
		err := rows.Scan(
			&resp.Count,
			&id,
			&name,
			&parent_id,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, err
		}
		resp.Categories = append(resp.Categories, models.Category{
			ID:        id.String,
			Name:      name.String,
			Parent_id: parent_id.String,
			CreatedAt: createdAt.String,
			UpdatedAt: updatedAt.String,
		})
	}
	return resp, nil
}

func (c *categoryRepo) UpdateCategory(req *models.UpdateCategory) (string, error) {

	query := `UPDATE category 
	            SET  name = $1, 
				     parent_id = $2, 
					 updated_at = NOW() 
					 WHERE id = $3 RETURNING id`

	result, err := c.db.Exec(context.Background(), query, req.Name, req.Parent_id, req.Id)
	if err != nil {
		return "Error Update Category", err
	}

	if result.RowsAffected() == 0 {
		return "", fmt.Errorf("category not found")
	}

	return req.Id, nil
}

func (c *categoryRepo) DeleteCategory(req *models.CategoryIdRequest) (resp string, err error) {
	query := `DELETE FROM category 
	            WHERE id = $1 RETURNING id`

	result, err := c.db.Exec(context.Background(), query, req.Id)
	if err != nil {
		return "Error from Delete Category", err
	}

	if result.RowsAffected() == 0 {
		return "", fmt.Errorf("Category not found")
	}

	return req.Id, nil
}
