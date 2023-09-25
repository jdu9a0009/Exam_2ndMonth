package postgres

import (
	"WareHouseProjects/config"
	"WareHouseProjects/storage"
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

type store struct {
	db                  *pgxpool.Pool
	branches            *branchRepo
	category            *categoryRepo
	product             *productRepo
	coming_table        *coming_tableRepo
	coming_tableProduct *coming_TableProductRepo
	remain              *remainRepo
}

func NewStorage(ctx context.Context, cfg config.Config) (storage.StorageI, error) {
	connect, err := pgxpool.ParseConfig(fmt.Sprintf(
		"host=%s user=%s dbname=%s password=%s port=%d sslmode=disable",
		cfg.PostgresHost,
		cfg.PostgresUser,
		cfg.PostgresDatabase,
		cfg.PostgresPassword,
		cfg.PostgresPort,
	))

	if err != nil {
		return nil, err
	}
	connect.MaxConns = cfg.PostgresMaxConnections

	pgxpool, err := pgxpool.ConnectConfig(context.Background(), connect)
	if err != nil {
		return nil, err
	}

	return &store{
		db: pgxpool,
	}, nil
}

func (b *store) Branch() storage.BranchesI {
	if b.branches == nil {
		b.branches = NewBranchRepo(b.db)
	}
	return b.branches
}

func (b *store) Category() storage.CategoriesI {
	if b.category == nil {
		b.category = NewCategoryRepo(b.db)
	}
	return b.category
}

func (b *store) Product() storage.ProdouctsI {
	if b.product == nil {
		b.product = NewProductRepo(b.db)
	}
	return b.product
}

func (b *store) Coming_Table() storage.Coming_TableI {
	if b.coming_table == nil {
		b.coming_table = NewComingTableRepo(b.db)
	}
	return b.coming_table
}

func (b *store) Coming_TableProduct() storage.Coming_TableProductI {
	if b.coming_tableProduct == nil {
		b.coming_tableProduct = NewComingTableProductRepo(b.db)
	}
	return b.coming_tableProduct
}

func (b *store) Remaining() storage.RemainingI {
	if b.remain == nil {
		b.remain = NewRemainRepo(b.db)
	}
	return b.remain
}

func (s *store) Close() {
	s.db.Close()
}
