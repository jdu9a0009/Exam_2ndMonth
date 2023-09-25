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

type branchRepo struct {
	db *pgxpool.Pool
}

func NewBranchRepo(db *pgxpool.Pool) *branchRepo {
	return &branchRepo{
		db: db,
	}
}

func (b *branchRepo) CreateBranch(req *models.CreateBranch) (string, error) {

	var (
		id    = uuid.NewString()
		query string
	)

	query = `
		INSERT INTO "branches"(
			"id", 
			"name",
			"address",
			"phone",
			"created_at" )
		VALUES ($1, $2, $3, $4, NOW())`

	_, err := b.db.Exec(context.Background(), query,
		id,
		req.Name,
		req.Address,
		req.Phone,
	)

	if err != nil {
		return "", err
	}

	return id, nil
}

func (b *branchRepo) GetBranch(req *models.BranchIdRequest) (resp *models.Branch, err error) {

	query := `
		SELECT
			"id", 
			"name",
			"address",
			"phone",
			"created_at",
			"updated_at" 
		FROM "branches"
		WHERE id = $1
	`
	var (
		createdAt time.Time
		updatedAt sql.NullTime
	)

	branch := models.Branch{}
	err = b.db.QueryRow(context.Background(), query, req.Id).Scan(
		&branch.ID,
		&branch.Name,
		&branch.Address,
		&branch.Phone,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("branch not found")
	}
	branch.CreatedAt = createdAt.Format(time.RFC3339)
	if updatedAt.Valid {
		branch.UpdatedAt = updatedAt.Time.Format(time.RFC3339)
	}

	return &branch, nil
}

func (b *branchRepo) GetAllBranch(req *models.GetAllBranchRequest) (*models.GetAllBranchResponse, error) {
	params := make(map[string]interface{})
	var resp = &models.GetAllBranchResponse{}

	resp.Branches = make([]models.Branch, 0)

	filter := " WHERE true "
	query := `
			SELECT
				COUNT(*) OVER(),
				"id", 
				"name",
				"address",
				"phone",
				"created_at",
				"updated_at" 
			FROM "branches"
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

	rows, err := b.db.Query(context.Background(), rquery, pArr...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			id        sql.NullString
			name      sql.NullString
			address   sql.NullString
			phone     sql.NullString
			createdAt sql.NullString
			updatedAt sql.NullString
		)
		err := rows.Scan(
			&resp.Count,
			&id,
			&name,
			&address,
			&phone,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, err
		}
		resp.Branches = append(resp.Branches, models.Branch{
			ID:        id.String,
			Name:      name.String,
			Address:   address.String,
			Phone:     phone.String,
			CreatedAt: createdAt.String,
			UpdatedAt: updatedAt.String,
		})
	}
	return resp, nil
}

func (b *branchRepo) UpdateBranch(req *models.UpdateBranch) (string, error) {

	query := `UPDATE branches 
	            SET  name = $1, 
				     address = $2, 
					 phone = $3, 
					 updated_at = NOW() 
					 WHERE id = $4 RETURNING id`

	result, err := b.db.Exec(context.Background(), query, req.Name, req.Address, req.Phone, req.Id)
	if err != nil {
		return "Error Update Branch", err
	}

	if result.RowsAffected() == 0 {
		return "", fmt.Errorf("branch not found")
	}

	return req.Id, nil
}

func (b *branchRepo) DeleteBranch(req *models.BranchIdRequest) (resp string, err error) {
	query := `DELETE FROM branches 
	            WHERE id = $1 RETURNING id`

	result, err := b.db.Exec(context.Background(), query, req.Id)
	if err != nil {
		return "Error from Delete Branch", err
	}

	if result.RowsAffected() == 0 {
		return "", fmt.Errorf("branch not found")
	}

	return req.Id, nil
}
