package handler

import (
	"WareHouseProjects/models"
	"WareHouseProjects/pkg/logger"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateProduct godoc
// @Router       /product [POST]
// @Summary      CREATE PRODUCT
// @Description adds product data to db based on given info in body
// @Tags         product
// @Accept       json
// @Produce      json
// @Param        data  body      models.CreateProduct  true  "product data"
// @Success      200  {string}  string
// @Failure      400  {object}  response.ErrorResp
// @Failure      404  {object}  response.ErrorResp
// @Failure      500  {object}  response.ErrorResp
func (h *Handler) CreateProduct(c *gin.Context) {
	var product models.CreateProduct
	err := c.ShouldBind(&product)
	if err != nil {
		h.log.Error("error while binding product:", logger.Error(err))
		c.JSON(http.StatusBadRequest, "invalid body")
		return
	}

	resp, err := h.storage.Product().CreateProduct(&product)
	if err != nil {
		h.log.Error("error product create:", logger.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"code": http.StatusCreated, "message": "success", "resp": resp})
}

// GetProduct godoc
// @Router       /product/{id} [GET]
// @Summary      GET BY ID
// @Description  gets product by ID
// @Tags         product
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Product ID" format(uuid)
// @Success      200  {object}  models.Product
// @Failure      400  {object}  response.ErrorResp
// @Failure      404  {object}  response.ErrorResp
// @Failure      500  {object}  response.ErrorResp
func (h *Handler) GetProduct(c *gin.Context) {
	id := c.Param("id")

	resp, err := h.storage.Product().GetProduct(&models.ProductIdRequest{Id: id})
	if err != nil {
		h.log.Error("error get product:", logger.Error(err))
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetALlProducts godoc
// @Router       /product [GET]
// @Summary      LIST PRODUCT
// @Description  gets all product based on limit, page and search by name
// @Tags         product
// @Accept       json
// @Produce      json
// @Param  		 limit         query     int        false  "limit"          minimum(1)     default(10)
// @Param  		 page          query     int        false  "page"           minimum(1)     default(1)
// @Param   	 name          query     string     false  "name"
// @Param   	 barcode       query     string     false  "barcode"
// @Success      200  {object}  models.GetAllProductRequest
// @Failure      400  {object}  response.ErrorResp
// @Failure      404  {object}  response.ErrorResp
// @Failure      500  {object}  response.ErrorResp
func (h *Handler) GetAllProduct(c *gin.Context) {
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

	resp, err := h.storage.Product().GetAllProduct(&models.GetAllProductRequest{
		Page:    page,
		Limit:   limit,
		Barcode: c.Query("search"),
		Name:    c.Query("search"),
	})
	if err != nil {
		h.log.Error("error Product GetAllProduct:", logger.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// UpdateProduct godoc
// @Router       /product/{id} [PUT]
// @Summary      UPDATE PRODUCT
// @Description  UPDATES PRODUCT BASED ON GIVEN DATA AND ID
// @Tags         product
// @Accept       json
// @Produce      json
// @Param        id    path     string  true  "id of product" format(uuid)
// @Param        data  body      models.UpdateProduct  true  "product data"
// @Success      200  {string}  string
// @Failure      400  {object}  response.ErrorResp
// @Failure      404  {object}  response.ErrorResp
// @Failure      500  {object}  response.ErrorResp
func (h *Handler) UpdateProduct(c *gin.Context) {
	var product models.UpdateProduct

	err := c.ShouldBind(&product)
	if err != nil {
		h.log.Error("error while binding:", logger.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	product.ID = c.Param("id")
	resp, err := h.storage.Product().UpdateProduct(&product)
	if err != nil {
		h.log.Error("error product update:", logger.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "success", "resp": resp})
}

// DeleteProduct godoc
// @Router       /product/{id} [DELETE]
// @Summary      DELETE PRODUCT BY ID
// @Description  deletes product by id
// @Tags         product
// @Accept       json
// @Produce      json
// @Param        id    path     string  true  "id of product" format(uuid)
// @Success      200  {string}  string
// @Failure      400  {object}  response.ErrorResp
// @Failure      404  {object}  response.ErrorResp
// @Failure      500  {object}  response.ErrorResp
func (h *Handler) DeleteProduct(c *gin.Context) {
	id := c.Param("id")

	resp, err := h.storage.Product().DeleteProduct(&models.ProductIdRequest{Id: id})
	if err != nil {
		h.log.Error("error deleting Product:", logger.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete Product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product successfully deleted", "id": resp})
}
