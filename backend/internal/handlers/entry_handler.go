package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/lockw1n/time-logger/internal/models"
)

type EntryHandler struct {
	DB *gorm.DB
}

func NewEntryHandler(db *gorm.DB) *EntryHandler {
	return &EntryHandler{DB: db}
}

// input structure for create/update
type entryInput struct {
	Ticket string  `json:"ticket"`
	Label  string  `json:"label"`
	Hours  float64 `json:"hours"`
	Date   string  `json:"date"` // "YYYY-MM-DD"
}

// ---------- LIST (with optional date range) ----------

func (h *EntryHandler) List(c *gin.Context) {
	var entries []models.Entry

	startStr := c.Query("start")
	endStr := c.Query("end")

	query := h.DB.Model(&models.Entry{})

	if startStr != "" && endStr != "" {
		start, err1 := time.Parse("2006-01-02", startStr)
		end, err2 := time.Parse("2006-01-02", endStr)
		if err1 != nil || err2 != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start or end date, expected YYYY-MM-DD"})
			return
		}
		// include end date
		end = end.Add(24 * time.Hour)
		query = query.Where("date >= ? AND date < ?", start, end)
	}

	if err := query.Order("ticket asc, date asc").Find(&entries).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch entries"})
		return
	}

	c.JSON(http.StatusOK, entries)
}

// ---------- GET BY ID ----------

func (h *EntryHandler) Get(c *gin.Context) {
	var entry models.Entry
	if err := h.DB.First(&entry, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "entry not found"})
		return
	}
	c.JSON(http.StatusOK, entry)
}

// ---------- CREATE ----------

func (h *EntryHandler) Create(c *gin.Context) {
	var input entryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if input.Ticket == "" || input.Hours < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ticket and non-negative hours are required"})
		return
	}

	var date time.Time
	var err error

	if input.Date == "" {
		// Default: current UTC time
		date = time.Now().UTC()
	} else {
		// Parse ISO-8601 UTC timestamp (e.g. "2025-11-11T18:22:30Z")
		date, err = time.Parse(time.RFC3339, input.Date)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date, expected RFC3339 format"})
			return
		}
		date = date.UTC()
	}

	entry := models.Entry{
		Ticket:    input.Ticket,
		Label:     input.Label,
		Hours:     input.Hours,
		Date:      date,
		CreatedAt: time.Now().UTC(),
	}

	if err := h.DB.Create(&entry).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save entry"})
		return
	}
	c.JSON(http.StatusCreated, entry)
}

// ---------- UPDATE BY ID ----------

func (h *EntryHandler) Update(c *gin.Context) {
	var entry models.Entry
	if err := h.DB.First(&entry, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "entry not found"})
		return
	}

	var input entryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	if input.Ticket != "" {
		entry.Ticket = input.Ticket
	}
	entry.Label = input.Label
	if input.Hours >= 0 {
		entry.Hours = input.Hours
	}

	if input.Date != "" {
		parsed, err := time.Parse(time.RFC3339, input.Date)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date"})
			return
		}
		entry.Date = parsed.UTC() // ✅ always stored as UTC
	} else {
		entry.Date = time.Now().UTC()
	}

	if err := h.DB.Save(&entry).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update entry"})
		return
	}
	c.JSON(http.StatusOK, entry)
}

// ---------- DELETE BY ID ----------

func (h *EntryHandler) Delete(c *gin.Context) {
	if err := h.DB.Delete(&models.Entry{}, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete entry"})
		return
	}
	c.Status(http.StatusNoContent)
}

// ---------- DELETE BY TICKET (delete row) ----------

func (h *EntryHandler) DeleteByTicket(c *gin.Context) {
	ticket := c.Param("ticket")
	if ticket == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ticket is required"})
		return
	}
	if err := h.DB.Where("ticket = ?", ticket).Delete(&models.Entry{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete ticket entries"})
		return
	}
	c.Status(http.StatusNoContent)
}
