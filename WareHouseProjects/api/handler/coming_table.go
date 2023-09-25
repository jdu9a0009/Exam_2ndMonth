package handler

import (
	"WareHouseProjects/models"
	"WareHouseProjects/pkg/logger"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateComingTable godoc
// @Router       /coming_table  [POST]
// @Summary      CREATE ComingTable
// @Description add ComingTable data to db based on given info in body
// @Tags         coming_table
// @Accept       json
// @Produce      json
// @Param        data  body      models.CreateComingTable true  "ComingTable data"
// @Success      200  {string}  string
// @Failure      400  {object}  response.ErrorResp
// @Failure      404  {object}  response.ErrorResp
// @Failure      500  {object}  response.ErrorResp
func (h *Handler) CreateComingTable(c *gin.Context) {
	var coming_table models.CreateComingTable
	err := c.ShouldBind(&coming_table)
	if err != nil {
		h.log.Error("error while binding coming table:", logger.Error(err))
		c.JSON(http.StatusBadRequest, "invalid body")
		return
	}

	resp, err := h.storage.Coming_Table().CreateComingTable(&coming_table)
	if err != nil {
		h.log.Error("error Coming_Table create:", logger.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"code": http.StatusCreated, "message": "success", "resp": resp})
}

// GetComingTable godoc
// @Router       /coming_table/{id} [GET]
// @Summary      GET BY ID
// @Description  gets ComingTable by ID
// @Tags         coming_table
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "ComingTable ID" format(uuid)
// @Success      200  {object}  models.ComingTable
// @Failure      400  {object}  response.ErrorResp
// @Failure      404  {object}  response.ErrorResp
// @Failure      500  {object}  response.ErrorResp
func (h *Handler) GetComingTable(c *gin.Context) {
	id := c.Param("id")

	resp, err := h.storage.Coming_Table().GetComingTable(&models.ComingTableIdRequest{Id: id})
	if err != nil {
		h.log.Error("error get ComingTable:", logger.Error(err))
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetAllComingTable godoc
// @Router       /coming_table [GET]
// @Summary      LIST Coming_Table
// @Description  gets all Coming_Table based on limit, page and search by name
// @Tags         coming_table
// @Accept       json
// @Produce      json
// @Param  		 limit         query     int        false  "limit"          minimum(1)     default(10)
// @Param  		 page          query     int        false  "page"           minimum(1)     default(1)
// @Param   	 coming_id        query     string     false  "coming_id"
// @Param   	 branch_id        query     string     false  "branch_id"
// @Success      200  {object}  models.GetAllComingTableRequest
// @Failure      400  {object}  response.ErrorResp
// @Failure      404  {object}  response.ErrorResp
// @Failure      500  {object}  response.ErrorResp
func (h *Handler) GetAllComingTable(c *gin.Context) {
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

	resp, err := h.storage.Coming_Table().GetAllComingTable(&models.GetAllComingTableRequest{
		Page:     page,
		Limit:    limit,
		ComingID: c.Query("search"),
		BranchID: c.Query("search"),
	})
	if err != nil {
		h.log.Error("error ComingTable GetAllComingTable:", logger.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// UpdateComingTable godoc
// @Router       /coming_table/{id} [PUT]
// @Summary      UPDATE COMINGTABLE
// @Description  UPDATES COMINGTABLE BASED ON GIVEN DATA AND ID
// @Tags         coming_table
// @Accept       json
// @Produce      json
// @Param        id    path     string  true  "id of ComingTable" format(uuid)
// @Param        data  body      models.UpdateComingTable  true  "ComingTable data"
// @Success      200  {string}  string
// @Failure      400  {object}  response.ErrorResp
// @Failure      404  {object}  response.ErrorResp
// @Failure      500  {object}  response.ErrorResp
func (h *Handler) UpdateComingTable(c *gin.Context) {
	var ComingTable models.UpdateComingTable

	err := c.ShouldBind(&ComingTable)
	if err != nil {
		h.log.Error("error while binding:", logger.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	ComingTable.ID = c.Param("id")
	resp, err := h.storage.Coming_Table().UpdateComingTable(&ComingTable)
	if err != nil {
		h.log.Error("error ComingTable update:", logger.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "success", "resp": resp})
}

// DeleteComingTable godoc
// @Router       /coming_table/{id} [DELETE]
// @Summary      DELETE ComingTable BY ID
// @Description  deletes ComingTable by id
// @Tags         coming_table
// @Accept       json
// @Produce      json
// @Param        id    path     string  true  "id of ComingTable" format(uuid)
// @Success      200  {string}  string
// @Failure      400  {object}  response.ErrorResp
// @Failure      404  {object}  response.ErrorResp
// @Failure      500  {object}  response.ErrorResp
func (h *Handler) DeleteComingTable(c *gin.Context) {
	id := c.Param("id")

	resp, err := h.storage.Coming_Table().DeleteComingTable(&models.ComingTableIdRequest{Id: id})
	if err != nil {
		h.log.Error("error deleting ComingTable:", logger.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete ComingTable"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ComingTable successfully deleted", "id": resp})
}
