package entry

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	entrydto "github.com/lockw1n/time-logger/internal/dto/entry"
	entryservice "github.com/lockw1n/time-logger/internal/service/entry"
)

type Handler struct {
	service entryservice.Service
}

func NewEntryHandler(service entryservice.Service) *Handler {
	return &Handler{service: service}
}

/*
|--------------------------------------------------------------------------
| CREATE ENTRY
|--------------------------------------------------------------------------
*/

func (h *Handler) Create(c *gin.Context) {
	var req entrydto.Request

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	resp, err := h.service.Create(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

/*
|--------------------------------------------------------------------------
| UPDATE ENTRY
|--------------------------------------------------------------------------
*/

func (h *Handler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 64)

	var req entrydto.Request

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	resp, err := h.service.Update(id, req)
	if err != nil {
		if errors.Is(err, entryservice.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "entry not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

/*
|--------------------------------------------------------------------------
| DELETE ENTRY
|--------------------------------------------------------------------------
*/

func (h *Handler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 64)

	err := h.service.Delete(id)
	if err != nil {
		if errors.Is(err, entryservice.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "entry not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
