package api

import (
	_ "WareHouseProjects/api/docs"
	"WareHouseProjects/api/handler"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewServer(h *handler.Handler) *gin.Engine {
	r := gin.Default()
	//Branches
	r.POST("/branch", h.CreateBranch)
	r.GET("/branch/:id", h.GetBranch)
	r.GET("/branch", h.GetAllBranch)
	r.PUT("/branch/:id", h.UpdateBranch)
	r.DELETE("/branch/:id", h.DeleteBranch)

	//Categories
	r.POST("/category", h.CreateCategory)
	r.GET("/category/:id", h.GetCategory)
	r.GET("/category", h.GetAllCategory)
	r.PUT("/category/:id", h.UpdateCategory)
	r.DELETE("/category/:id", h.DeleteCategory)

	//Product
	r.POST("/product", h.CreateProduct)
	r.GET("/product/:id", h.GetProduct)
	r.GET("/product", h.GetAllProduct)
	r.PUT("/product/:id", h.UpdateProduct)
	r.DELETE("/product/:id", h.DeleteProduct)

	//ComingTable
	r.POST("/coming_table", h.CreateComingTable)
	r.GET("/coming_table/:id", h.GetComingTable)
	r.GET("/coming_table", h.GetAllComingTable)
	r.PUT("/coming_table/:id", h.UpdateComingTable)
	r.DELETE("/coming_table/:id", h.DeleteComingTable)

	//ComingTableProduct
	r.POST("/coming_table_product", h.CreateComingTableProduct)
	r.GET("/coming_table_product/:id", h.GetComingTableProduct)
	r.GET("/coming_table_product", h.GetAllComingTableProduct)
	r.PUT("/coming_table_product/:id", h.UpdateComingTableProduct)
	r.DELETE("/coming_table_product/:id", h.DeleteComingTableProduct)

	//Remain
	r.POST("/do_income/:coming_table_id", h.CreateRemain)
	r.GET("/remain/:id", h.GetRemain)
	r.GET("/remain", h.GetAllRemain)
	r.PUT("/remain/:id", h.UpdateRemain)
	r.DELETE("/remain/:id", h.DeleteRemain)

	url := ginSwagger.URL("swagger/doc.json")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	return r
}
