package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type EntryHandler struct {
	Service *EntryService
}

func NewEntryHandler(service *EntryService) *EntryHandler {
	return &EntryHandler{Service: service}
}

// input structure for create/update
type entryInput struct {
	Ticket string  `json:"ticket"`
	Label  string  `json:"label"`
	Hours  float64 `json:"hours"`
	Date   string  `json:"date"` // "YYYY-MM-DD"
}

type ticketSummary struct {
	Ticket     string  `json:"ticket"`
	Label      string  `json:"label"`
	TotalHours float64 `json:"total_hours"`
}

// ---------- LIST (with optional date range) ----------

func (h *EntryHandler) List(c *gin.Context) {
	entries, err := h.Service.ListEntries(c.Query("start"), c.Query("end"))
	if err != nil {
		if err == ErrBadDateRange {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start or end date, expected YYYY-MM-DD"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch entries"})
		return
	}

	c.JSON(http.StatusOK, entries)
}

// ---------- SUMMARY BY TICKET ----------

func (h *EntryHandler) Summary(c *gin.Context) {
	start := c.Query("start")
	end := c.Query("end")
	if start == "" || end == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start and end are required (YYYY-MM-DD)"})
		return
	}

	summaries, err := h.Service.Summary(start, end)
	if err != nil {
		if err == ErrBadDateRange {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start or end date, expected YYYY-MM-DD"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch summary"})
		return
	}

	c.JSON(http.StatusOK, summaries)
}

// ---------- MONTHLY REPORT (for invoices) ----------

func (h *EntryHandler) MonthlyReport(c *gin.Context) {
	month := c.Query("month")
	report, err := h.Service.MonthlySummary(month)
	if err != nil {
		if err == ErrBadMonth {
			c.JSON(http.StatusBadRequest, gin.H{"error": "month is required and must be in YYYY-MM format"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to build monthly summary"})
		return
	}

	c.JSON(http.StatusOK, report)
}

// ---------- GET BY ID ----------

func (h *EntryHandler) Get(c *gin.Context) {
	entry, err := h.Service.GetEntry(c.Param("id"))
	if err != nil {
		if err == ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "entry not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch entry"})
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

	entry, created, err := h.Service.CreateOrSum(input)
	if err != nil {
		switch err {
		case ErrInvalidHours:
			c.JSON(http.StatusBadRequest, gin.H{"error": "ticket is required and hours must be in 15-minute increments between 0 and 24"})
		case ErrInvalidLabel:
			c.JSON(http.StatusBadRequest, gin.H{"error": h.Service.AllowedLabelsError()})
		case ErrInvalidDate:
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date, expected YYYY-MM-DD"})
		case ErrNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "entry not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save entry"})
		}
		return
	}

	if created {
		c.JSON(http.StatusCreated, entry)
		return
	}
	c.JSON(http.StatusOK, entry)
}

// ---------- UPDATE BY ID ----------

func (h *EntryHandler) Update(c *gin.Context) {
	var input entryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	entry, err := h.Service.UpdateEntry(c.Param("id"), input)
	if err != nil {
		switch err {
		case ErrInvalidHours:
			c.JSON(http.StatusBadRequest, gin.H{"error": "hours must be in 15-minute increments between 0 and 24"})
		case ErrInvalidLabel:
			c.JSON(http.StatusBadRequest, gin.H{"error": h.Service.AllowedLabelsError()})
		case ErrInvalidDate:
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date, expected YYYY-MM-DD"})
		case ErrNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "entry not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update entry"})
		}
		return
	}

	c.JSON(http.StatusOK, entry)
}

// ---------- DELETE BY ID ----------

func (h *EntryHandler) Delete(c *gin.Context) {
	if err := h.Service.DeleteEntry(c.Param("id")); err != nil {
		if err == ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "entry not found"})
			return
		}
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
	if err := h.Service.DeleteByTicket(ticket); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete ticket entries"})
		return
	}
	c.Status(http.StatusNoContent)
}
