package handler

import (
	"WareHouseProjects/api/handler/response"
	"WareHouseProjects/models"
	"WareHouseProjects/pkg/logger"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateCategory godoc
// @Router       /category [POST]
// @Summary      CREATES Category
// @Description  CREATES Category BASED ON GIVEN DATA
// @Tags         category
// @Accept       json
// @Produce      json
// @Param        data  body      models.CreateCategory  true  "Category data"
// @Success      200  {string}  string
// @Failure      400  {object}  response.ErrorResp
// @Failure      404  {object}  response.ErrorResp
// @Failure      500  {object}  response.ErrorResp
func (h *Handler) CreateCategory(c *gin.Context) {
	var category models.CreateCategory
	err := c.ShouldBind(&category)
	if err != nil {
		h.log.Error("error while binding:", logger.Error(err))
		c.JSON(http.StatusBadRequest, "invalid body")
		return
	}

	resp, err := h.storage.Category().CreateCategory(&category)
	if err != nil {
		h.log.Error("error Category Create:", logger.Error(err))
		c.JSON(http.StatusInternalServerError, "internal server error")
		return
	}
	c.JSON(http.StatusCreated, response.CreateResponse{Message: "Succesfully created", Id: resp})
}

// Get category godoc
// @Router       /category/{id} [GET]
// @Summary      GET BY ID
// @Description  get category by ID
// @Tags         category
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Category ID" format(uuid)
// @Success      200  {object}  models.Category
// @Failure      400  {object}  response.ErrorResp
// @Failure      404  {object}  response.ErrorResp
// @Failure      500  {object}  response.ErrorResp
func (h *Handler) GetCategory(c *gin.Context) {
	id := c.Param("id")

	resp, err := h.storage.Category().GetCategory(&models.CategoryIdRequest{Id: id})
	if err != nil {
		h.log.Error("error Category Get:", logger.Error(err))
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// GetAllCategories godoc
// @Router       /category [GET]
// @Summary      GET  ALL CATEGORIES
// @Description  get all categories based on limit, page and search by name
// @Tags         category
// @Accept       json
// @Produce      json
// @Param   limit         query     int        false  "limit"          minimum(1)     default(10)
// @Param   page         query     int        false  "page"          minimum(1)     default(1)
// @Param   search         query     string        false  "search"
// @Success      200  {object}  models.GetAllCategoryRequest
// @Failure      400  {object}  response.ErrorResp
// @Failure      404  {object}  response.ErrorResp
// @Failure      500  {object}  response.ErrorResp
func (h *Handler) GetAllCategory(c *gin.Context) {
	h.log.Info("request GetAllCategory")
	page, err := strconv.Atoi(c.DefaultQuery("page", "fmt.sprintf(`%d`,cfg.DefaultPage)"))
	if err != nil {
		h.log.Error("error getting page:", logger.Error(err))
		c.JSON(http.StatusBadRequest, "invalid page param")
		return
	}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		h.log.Error("error getting limit:", logger.Error(err))
		c.JSON(http.StatusBadRequest, "invalid limit param")
		return
	}

	resp, err := h.storage.Category().GetAllCategory(&models.GetAllCategoryRequest{
		Page:  page,
		Limit: limit,
		Name:  c.Query("search"),
	})
	if err != nil {
		h.log.Error("error Category GetAllCategory:", logger.Error(err))
		c.JSON(http.StatusInternalServerError, "internal server error")
		return
	}
	h.log.Warn("response to GetAllCategory")
	c.JSON(http.StatusOK, resp)
}

// UpdateCategory godoc
// @Router       /category/{id} [PUT]
// @Summary      UPDATE Category BY ID
// @Description  UPDATES Category BASED ON GIVEN DATA AND ID
// @Tags         category
// @Accept       json
// @Produce      json
// @Param        id    path     string  true  "id of category" format(uuid)
// @Param        data  body      models.UpdateCategory true  "category data"
// @Success      200  {string}  string
// @Failure      400  {object}  response.ErrorResp
// @Failure      404  {object}  response.ErrorResp
// @Failure      500  {object}  response.ErrorResp
func (h *Handler) UpdateCategory(ctx *gin.Context) {
	var category models.UpdateCategory

	err := ctx.ShouldBind(&category)
	if err != nil {
		h.log.Error("error while binding:", logger.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	category.Id = ctx.Param("id")
	resp, err := h.storage.Category().UpdateCategory(&category)
	if err != nil {
		h.log.Error("error category update:", logger.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "success", "resp": resp})
}

// DeleteCategory godoc
// @Router       /category/{id} [DELETE]
// @Summary      DELETE CATEGORY BY ID
// @Description  DELETES CATEGORY BASED ON ID
// @Tags         category
// @Accept       json
// @Produce      json
// @Param        id    path     string  true  "id of category" format(uuid)
// @Success      200  {string}  string
// @Failure      400  {object}  response.ErrorResp
// @Failure      404  {object}  response.ErrorResp
// @Failure      500  {object}  response.ErrorResp
func (h *Handler) DeleteCategory(c *gin.Context) {
	id := c.Param("id")

	resp, err := h.storage.Category().DeleteCategory(&models.CategoryIdRequest{Id: id})
	if err != nil {
		h.log.Error("error deleting category:", logger.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete category"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category successfully deleted", "id": resp})
}
