package handler

import (
	"WareHouseProjects/models"
	"WareHouseProjects/pkg/logger"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateComingTableProduct godoc
// @Router       /coming_table_product  [POST]
// @Summary      CREATE ComingTableProduct
// @Description add ComingTableProduct data to db based on given info in body
// @Tags         coming_table_product
// @Accept       json
// @Produce      json
// @Param        data  body      models.CreateComingTableProductSwagger true  "ComingTableProduct data"
// @Success      200  {string}  string
// @Failure      400  {object}  response.ErrorResp
// @Failure      404  {object}  response.ErrorResp
// @Failure      500  {object}  response.ErrorResp
func (h *Handler) CreateComingTableProduct(c *gin.Context) {
	var coming_tableProduct models.CreateComingTableProduct
	err := c.ShouldBind(&coming_tableProduct)
	if err != nil {
		h.log.Error("error while binding coming table_product:", logger.Error(err))
		c.JSON(http.StatusBadRequest, "invalid body")
		return
	}

	//check status
	comingTableId := c.Param("coming_table_id")
	coming_table_id := models.ComingTableIdRequest{Id: comingTableId}
	_, err = h.storage.Coming_Table().GetStatus(&coming_table_id)
	if err != nil {
		h.log.Error("error getting coming table status:", logger.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//get  product details
	CheckBarcodeComingTable := models.CheckBarcodeComingTable{Barcode: coming_tableProduct.Barcode, Coming_Table_id: coming_tableProduct.Coming_Table_id}
	respondProduct, err := h.storage.Product().GetProductByBarcode(&CheckBarcodeComingTable)
	if err != nil {
		h.log.Error("error Getting Product info:", logger.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	coming_tableProduct.Name = respondProduct.Name
	coming_tableProduct.Price = respondProduct.Price
	coming_tableProduct.Category_id = respondProduct.Category_id
	coming_tableProduct.TotalPrice = respondProduct.Price * coming_tableProduct.Count

	barcodeQuery := c.Query("barcode")
	barcode := models.CheckBarcodeComingTable{Barcode: barcodeQuery, Coming_Table_id: comingTableId}
	id, err := h.storage.Coming_TableProduct().CheckAviableProduct(&barcode)

	if err != nil {
		h.log.Error("error   id or barcode not found in comingTable", logger.Error(err))
		// if this product didnt exist Add it
		resp, err := h.storage.Coming_TableProduct().CreateComingTableProduct(&coming_tableProduct)
		if err != nil {
			h.log.Error("error Coming_Table_Product create:", logger.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return

		}

		c.JSON(http.StatusCreated, gin.H{"code": http.StatusCreated, "message": "successfully added ", "resp": resp})
		return
	}
	updatingData := models.UpdateComingTableProduct{
		ID:              id,
		Category_id:     coming_tableProduct.Category_id,
		Name:            coming_tableProduct.Name,
		Price:           coming_tableProduct.Price,
		Barcode:         coming_tableProduct.Barcode,
		Count:           coming_tableProduct.Count,
		TotalPrice:      coming_tableProduct.TotalPrice,
		Coming_Table_id: coming_tableProduct.Coming_Table_id,
	}

	resp, err := h.storage.Coming_TableProduct().UpdateIdAviable(&updatingData)
	if err != nil {
		h.log.Error("error Updating  coming_table_product:", logger.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return

	}
	c.JSON(http.StatusOK, gin.H{"message": "updated existing coming_product_table", "resp": resp})
}

/*
func (h *Handler) CreateComingTableProduct(c *gin.Context) {
	var coming_tableProduct models.CreateComingTableProduct
	err := c.ShouldBind(&coming_tableProduct)
	if err != nil {
		h.log.Error("error while binding coming table_product:", logger.Error(err))
		c.JSON(http.StatusBadRequest, "invalid body")
		return
	}

	resp, err := h.storage.Coming_TableProduct().CreateComingTableProduct(&coming_tableProduct)
	if err != nil {
		h.log.Error("error Coming_Table_Product create:", logger.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"code": http.StatusCreated, "message": "success", "resp": resp})
}
*/

// GetComingTableProduct godoc
// @Router       /coming_table_product/{id} [GET]
// @Summary      GET BY ID
// @Description  gets ComingTableProduct by ID
// @Tags         coming_table_product
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "ComingTableProduct ID" format(uuid)
// @Success      200  {object}  models.ComingTableProduct
// @Failure      400  {object}  response.ErrorResp
// @Failure      404  {object}  response.ErrorResp
// @Failure      500  {object}  response.ErrorResp
func (h *Handler) GetComingTableProduct(c *gin.Context) {
	id := c.Param("id")

	resp, err := h.storage.Coming_TableProduct().GetComingTableProduct(&models.ComingTableProductIdRequest{Id: id})
	if err != nil {
		h.log.Error("error get ComingTableProduct:", logger.Error(err))
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetAllComingTableProduct godoc
// @Router       /coming_table_product [GET]
// @Summary      LIST Coming_Table_Product
// @Description  gets all Coming_TableProduct based on limit, page and search by name
// @Tags         coming_table_product
// @Accept       json
// @Produce      json
// @Param  		 limit         query     int        false  "limit"          minimum(1)     default(10)
// @Param  		 page          query     int        false  "page"           minimum(1)     default(1)
// @Param   	 category_id   query     string     false  "category_id"
// @Param   	 barcode       query     string     false  "barcode"
// @Success      200  {object}  models.GetAllComingTableProductRequest
// @Failure      400  {object}  response.ErrorResp
// @Failure      404  {object}  response.ErrorResp
// @Failure      500  {object}  response.ErrorResp
func (h *Handler) GetAllComingTableProduct(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		h.log.Error("error get page:", logger.Error(err))
		c.JSON(http.StatusBadRequest, "invalid page param")
		return
	}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		h.log.Error("error get limit:", logger.Error(err))
		c.JSON(http.StatusBadRequest, "invalid page param")
		return
	}

	resp, err := h.storage.Coming_TableProduct().GetAllComingTableProduct(&models.GetAllComingTableProductRequest{
		Page:        page,
		Limit:       limit,
		Category_id: c.Query("search"),
		Barcode:     c.Query("search"),
	})
	if err != nil {
		h.log.Error("error ComingTableProduct GetAllComingTableProduct:", logger.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// UpdateComingTableProduct godoc
// @Router       /coming_table_product/{id} [PUT]
// @Summary      UPDATE COMINGTableProduct
// @Description  UPDATES COMINGTableProduct BASED ON GIVEN DATA AND ID
// @Tags         coming_table_product
// @Accept       json
// @Produce      json
// @Param        id    path     string  true  "id of ComingTableProduct" format(uuid)
// @Param        data  body      models.UpdateComingTableProduct  true  "ComingTableProduct data"
// @Success      200  {string}  string
// @Failure      400  {object}  response.ErrorResp
// @Failure      404  {object}  response.ErrorResp
// @Failure      500  {object}  response.ErrorResp
func (h *Handler) UpdateComingTableProduct(c *gin.Context) {
	var ComingTableProduct models.UpdateComingTableProduct

	err := c.ShouldBind(&ComingTableProduct)
	if err != nil {
		h.log.Error("error while binding:", logger.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	ComingTableProduct.ID = c.Param("id")
	resp, err := h.storage.Coming_TableProduct().UpdateComingTableProduct(&ComingTableProduct)
	if err != nil {
		h.log.Error("error ComingTableProduct update:", logger.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "success", "resp": resp})
}

// DeleteComingTableProduct godoc
// @Router       /coming_table_product/{id} [DELETE]
// @Summary      DELETE ComingTableProduct BY ID
// @Description  deletes ComingTableProduct by id
// @Tags         coming_table_product
// @Accept       json
// @Produce      json
// @Param        id    path     string  true  "id of ComingTableProduct" format(uuid)
// @Success      200  {string}  string
// @Failure      400  {object}  response.ErrorResp
// @Failure      404  {object}  response.ErrorResp
// @Failure      500  {object}  response.ErrorResp
func (h *Handler) DeleteComingTableProduct(c *gin.Context) {
	id := c.Param("id")

	resp, err := h.storage.Coming_TableProduct().DeleteComingTableProduct(&models.ComingTableProductIdRequest{Id: id})
	if err != nil {
		h.log.Error("error deleting ComingTableProduct:", logger.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete ComingTableProduct"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ComingTableProduct successfully deleted", "id": resp})
}
