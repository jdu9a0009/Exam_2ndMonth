package storage

import models "WareHouseProjects/models"

type StorageI interface {
	Branch() BranchesI
	Category() CategoriesI
	Product() ProdouctsI
	Coming_Table() Coming_TableI
	Coming_TableProduct() Coming_TableProductI
	Remaining() RemainingI

	Close()
}

type BranchesI interface {
	CreateBranch(*models.CreateBranch) (string, error)
	GetBranch(*models.BranchIdRequest) (*models.Branch, error)
	GetAllBranch(*models.GetAllBranchRequest) (*models.GetAllBranchResponse, error)
	UpdateBranch(*models.UpdateBranch) (string, error)
	DeleteBranch(*models.BranchIdRequest) (string, error)
}

type CategoriesI interface {
	CreateCategory(*models.CreateCategory) (string, error)
	GetCategory(*models.CategoryIdRequest) (*models.Category, error)
	GetAllCategory(*models.GetAllCategoryRequest) (*models.GetAllCategoryResponse, error)
	UpdateCategory(*models.UpdateCategory) (string, error)
	DeleteCategory(*models.CategoryIdRequest) (string, error)
}

type ProdouctsI interface {
	CreateProduct(*models.CreateProduct) (string, error)
	GetProduct(*models.ProductIdRequest) (*models.Product, error)
	GetAllProduct(*models.GetAllProductRequest) (*models.GetAllProductResponse, error)
	UpdateProduct(*models.UpdateProduct) (string, error)
	DeleteProduct(*models.ProductIdRequest) (string, error)

	GetProductByBarcode(*models.CheckBarcodeComingTable) (*models.RespBarcodeProduct, error)
}

type Coming_TableI interface {
	CreateComingTable(*models.CreateComingTable) (string, error)
	GetComingTable(*models.ComingTableIdRequest) (*models.ComingTable, error)
	GetAllComingTable(*models.GetAllComingTableRequest) (*models.GetAllComingTableResponse, error)
	UpdateComingTable(*models.UpdateComingTable) (string, error)
	DeleteComingTable(*models.ComingTableIdRequest) (string, error)

	GetStatus(*models.ComingTableIdRequest) (string, error)
	UpdateStatus(req *models.ComingTableIdRequest) (string, error)
}

type Coming_TableProductI interface {
	CreateComingTableProduct(*models.CreateComingTableProduct) (string, error)
	GetComingTableProduct(*models.ComingTableProductIdRequest) (*models.ComingTableProduct, error)
	GetAllComingTableProduct(*models.GetAllComingTableProductRequest) (*models.GetAllComingTableProductResponse, error)
	UpdateComingTableProduct(*models.UpdateComingTableProduct) (string, error)
	DeleteComingTableProduct(*models.ComingTableProductIdRequest) (string, error)

	CheckAviableProduct(*models.CheckBarcodeComingTable) (string, error)
	UpdateIdAviable(*models.UpdateComingTableProduct) (string, error)
	GetComingTableById(*models.ComingTableProductIdRequest) (*models.ComingTableProduct, error)
}

type RemainingI interface {
	CreateRemain(*models.CreateRemain) (string, error)
	GetRemain(*models.RemainIdRequest) (*models.Remain, error)
	GetAllRemain(*models.GetAllRemainRequest) (resp *models.GetAllRemainResponse, err error)
	UpdateRemain(*models.UpdateRemain) (string, error)
	DeleteRemain(*models.RemainIdRequest) (string, error)

	UpdateIdAviable(req *models.UpdateRemain) (string, error)
	CheckRemain(req *models.CheckRemain) (string, error)
}
