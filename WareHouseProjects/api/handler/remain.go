package handler

import (
	"WareHouseProjects/models"
	"WareHouseProjects/pkg/logger"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateRemain godoc
// @Router       /do_income/{coming_table_id} [POST]
// @Summary      CREATE Remain
// @Description adds Remain data to db based on given id
// @Tags         remain
// @Accept       json
// @Produce      json
// @Param        coming_table_id path string true "Coming Table ID"
// @Success      200  {string}  string
// @Failure      400  {object}  response.ErrorResp
// @Failure      404  {object}  response.ErrorResp
// @Failure      500  {object}  response.ErrorResp
func (h *Handler) CreateRemain(c *gin.Context) {
	comingTableID := c.Param("coming_table_id")
	var remain models.CreateRemain

	// Check status
	comingTableIDRequest := models.ComingTableIdRequest{Id: comingTableID}
	branchID, err := h.storage.Coming_Table().GetStatus(&comingTableIDRequest)
	if err != nil {
		h.log.Error("error while getting coming table status", logger.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	remain.Branch_id = branchID

	comingIDRequest := models.ComingTableProductIdRequest{Id: comingTableID}
	comingTableData, err := h.storage.Coming_TableProduct().GetComingTableById(&comingIDRequest)
	if err != nil {
		h.log.Error("error while getting coming table data details:", logger.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Not Found Coming_table_data with that coming_id"})
		return
	}

	remain.Barcode = comingTableData.Barcode
	remain.Category_id = comingTableData.Category_id
	remain.Price = comingTableData.Price
	remain.Count = comingTableData.Count
	remain.TotalPrice = comingTableData.TotalPrice
	remain.Name = comingTableData.Name

	checkRemainRequest := models.CheckRemain{Branch_id: branchID, Barcode: remain.Barcode}
	id, err := h.storage.Remaining().CheckRemain(&checkRemainRequest)
	if err != nil {
		h.log.Info("remaining not found, creating new remaining", logger.Error(err))
		resp, err := h.storage.Remaining().CreateRemain(&remain)
		if err != nil {
			h.log.Error("error creating remaining:", logger.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"code": http.StatusCreated, "message": "added new remaining", "resp": resp})
		h.storage.Coming_Table().UpdateStatus(&comingTableIDRequest)
		return
	}

	updatingData := models.UpdateRemain{
		ID:          id,
		Branch_id:   remain.Branch_id,
		Category_id: remain.Category_id,
		Name:        remain.Name,
		Price:       remain.Price,
		Barcode:     remain.Barcode,
		Count:       remain.Count,
		TotalPrice:  remain.TotalPrice,
	}
	r, err := h.storage.Remaining().UpdateIdAviable(&updatingData)
	if err != nil {
		h.log.Info("error updating remaining:", logger.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "updated existing remaining table", "resp": r})

	// If everything is ok, change status to finished
	h.storage.Coming_Table().UpdateStatus(&comingTableIDRequest)
}

// GetRemain godoc
// @Router       /remain/{id} [GET]
// @Summary      GET BY ID
// @Description  gets Remain by ID
// @Tags         remain
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Remain ID" format(uuid)
// @Success      200  {object}  models.Remain
// @Failure      400  {object}  response.ErrorResp
// @Failure      404  {object}  response.ErrorResp
// @Failure      500  {object}  response.ErrorResp
func (h *Handler) GetRemain(c *gin.Context) {
	id := c.Param("id")

	resp, err := h.storage.Remaining().GetRemain(&models.RemainIdRequest{Id: id})
	if err != nil {
		h.log.Error("error get Remain:", logger.Error(err))
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetALlRemains godoc
// @Router       /remain [GET]
// @Summary      LIST Remain
// @Description  gets all Remain based on limit, page and search by name
// @Tags         remain
// @Accept       json
// @Produce      json
// @Param  		 limit         query     int        false  "limit"          minimum(1)     default(10)
// @Param  		 page          query     int        false  "page"           minimum(1)     default(1)
// @Param   	 branch_id        query     string     false  "branch_id"
// @Param   	 category_id        query     string     false  "category_id"
// @Param   	 barcode        query     string     false  "barcode"
// @Success      200  {object}  models.GetAllRemainRequest
// @Failure      400  {object}  response.ErrorResp
// @Failure      404  {object}  response.ErrorResp
// @Failure      500  {object}  response.ErrorResp
func (h *Handler) GetAllRemain(c *gin.Context) {
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

	resp, err := h.storage.Remaining().GetAllRemain(&models.GetAllRemainRequest{
		Page:        page,
		Limit:       limit,
		Branch_id:   c.Query("search"),
		Category_id: c.Query("search"),
		Barcode:     c.Query("search"),
	})
	if err != nil {
		h.log.Error("error Remain GetAllRemain:", logger.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// UpdateRemain godoc
// @Router       /remain/{id} [PUT]
// @Summary      UPDATE Remain
// @Description  UPDATES Remain BASED ON GIVEN DATA AND ID
// @Tags         remain
// @Accept       json
// @Produce      json
// @Param        id    path     string  true  "id of Remain" format(uuid)
// @Param        data  body      models.UpdateRemain  true  "Remain data"
// @Success      200  {string}  string
// @Failure      400  {object}  response.ErrorResp
// @Failure      404  {object}  response.ErrorResp
// @Failure      500  {object}  response.ErrorResp
func (h *Handler) UpdateRemain(c *gin.Context) {
	var Remain models.UpdateRemain

	err := c.ShouldBind(&Remain)
	if err != nil {
		h.log.Error("error while binding:", logger.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	Remain.ID = c.Param("id")
	resp, err := h.storage.Remaining().UpdateRemain(&Remain)
	if err != nil {
		h.log.Error("error Remain update:", logger.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "success", "resp": resp})
}

// DeleteRemain godoc
// @Router       /remain/{id} [DELETE]
// @Summary      DELETE Remain BY ID
// @Description  deletes Remain by id
// @Tags         remain
// @Accept       json
// @Produce      json
// @Param        id    path     string  true  "id of Remain" format(uuid)
// @Success      200  {string}  string
// @Failure      400  {object}  response.ErrorResp
// @Failure      404  {object}  response.ErrorResp
// @Failure      500  {object}  response.ErrorResp
func (h *Handler) DeleteRemain(c *gin.Context) {
	id := c.Param("id")

	resp, err := h.storage.Remaining().DeleteRemain(&models.RemainIdRequest{Id: id})
	if err != nil {
		h.log.Error("error deleting Remain:", logger.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete Remain"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Remain successfully deleted", "id": resp})
}
