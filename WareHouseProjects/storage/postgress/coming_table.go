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

type coming_tableRepo struct {
	db *pgxpool.Pool
}

func NewComingTableRepo(db *pgxpool.Pool) *coming_tableRepo {
	return &coming_tableRepo{
		db: db,
	}
}

func (c *coming_tableRepo) CreateComingTable(req *models.CreateComingTable) (resp string, err error) {
	id := uuid.NewString()
	// comingId := helper.NewCustomIDGenerator().GenerateID()

	query := `
	INSERT INTO coming_table(
	  id,
	  coming_id,
	  branch_id,
	  date_time
	) VALUES($1,$2,$3,$4)	`

	_, err = c.db.Exec(context.Background(), query,
		id,
		req.Coming_id,
		req.Branch_id,
		req.DateTime,
	)

	if err != nil {
		return
	}

	return id, nil
}
func (c *coming_tableRepo) GetComingTable(req *models.ComingTableIdRequest) (resp *models.ComingTable, err error) {

	query := `
		SELECT
		    "id", 
		    "coming_id",
		    "branch_id",
		    "date_time",
		    "status",
		    "created_at",
			"updated_at" 
		FROM "coming_table"
		WHERE id = $1
	`
	var (
		createdAt time.Time
		updatedAt sql.NullTime
	)

	ComingTable := models.ComingTable{}
	err = c.db.QueryRow(context.Background(), query, req.Id).Scan(
		&ComingTable.ID,
		&ComingTable.ComingID,
		&ComingTable.BranchID,
		&ComingTable.DateTime,
		&ComingTable.Status,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf(" ComingTable not found")
	}
	ComingTable.CreatedAt = createdAt.Format(time.RFC3339)
	if updatedAt.Valid {
		ComingTable.UpdatedAt = updatedAt.Time.Format(time.RFC3339)
	}

	return &ComingTable, nil
}

func (c *coming_tableRepo) GetAllComingTable(req *models.GetAllComingTableRequest) (*models.GetAllComingTableResponse, error) {
	params := make(map[string]interface{})
	var resp = &models.GetAllComingTableResponse{}

	resp.ComingTables = make([]models.ComingTable, 0)

	filter := " WHERE true "
	query := `
			SELECT
				COUNT(*) OVER(),
				"id", 
				"coming_id",
				"branch_id",
				"date_time",
				"status",
				"created_at",
				"updated_at" 
			FROM "coming_table"
		`

	if req.ComingID != "" {
		filter += ` AND ("coming_id" ILIKE '%' || :coming_id ) `
		params["search"] = req.ComingID
	}
	if req.BranchID != "" {
		filter += ` AND ("branch_id" ILIKE '%' || :branch_id ) `
		params["search"] = req.BranchID
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
			coming_id sql.NullString
			branch_id sql.NullString
			date_time sql.NullTime
			status    sql.NullString
			createdAt sql.NullString
			updatedAt sql.NullString
		)
		err := rows.Scan(
			&resp.Count,
			&id,
			&coming_id,
			&branch_id,
			&date_time,
			&status,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, err
		}
		resp.ComingTables = append(resp.ComingTables, models.ComingTable{
			ID:        id.String,
			ComingID:  coming_id.String,
			BranchID:  branch_id.String,
			DateTime:  date_time.Time.Format(time.DateTime),
			Status:    models.TableType(status.String),
			CreatedAt: createdAt.String,
			UpdatedAt: updatedAt.String,
		})
	}
	return resp, nil
}

func (c *coming_tableRepo) UpdateComingTable(req *models.UpdateComingTable) (string, error) {

	query := `UPDATE coming_table 
	            SET  coming_id = $1, 
				     branch_id = $2, 
					 date_time=$3,
					 updated_at = NOW() 
					 WHERE id = $4 RETURNING id`

	result, err := c.db.Exec(context.Background(), query, req.ComingID, req.BranchID, req.DateTime, req.ID)
	if err != nil {
		return "Error Update Coming_Table", err
	}

	if result.RowsAffected() == 0 {
		return "", fmt.Errorf("Coming_Table not found")
	}

	return req.ID, nil
}

func (c *coming_tableRepo) DeleteComingTable(req *models.ComingTableIdRequest) (resp string, err error) {
	query := `DELETE FROM coming_table 
	            WHERE id = $1 RETURNING id`

	result, err := c.db.Exec(context.Background(), query, req.Id)
	if err != nil {
		return "Error from Delete Coming_Table", err
	}

	if result.RowsAffected() == 0 {
		return "", fmt.Errorf("Coming_Table not found")
	}

	return req.Id, nil
}
func (c *coming_tableRepo) UpdateStatus(req *models.ComingTableIdRequest) (string, error) {
	query := `Update coming_table Set
	            status=$1,
				updated_at=now()
				where id=$2`
	resp, err := c.db.Exec(context.Background(), query, "finished", req.Id)
	if err != nil {
		return "", err
	}

	if resp.RowsAffected() == 0 {
		return "", fmt.Errorf("coming_table with ID %s not found", req.Id)
	}

	return req.Id, nil
}

func (c *coming_tableRepo) GetStatus(req *models.ComingTableIdRequest) (string, error) {
	var status sql.NullString

	var branch_id sql.NullString
	parsedUUID, err := uuid.Parse(req.Id)
	if err != nil {
		return "", fmt.Errorf("invalid UUID format: %v", err)
	}

	query := `
		SELECT 
		   status,
		   branch_id
		FROM coming_table
		WHERE id = $1::uuid
	`

	err = c.db.QueryRow(context.Background(), query, parsedUUID).Scan(&status, &branch_id)
	if err != nil {
		return "", err
	}

	if status.Valid && status.String == "finished" {
		return "", fmt.Errorf("coming table already finished")
	}

	return branch_id.String, nil
}
