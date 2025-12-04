package handlers

import (
	"errors"
	"math"
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/lockw1n/time-logger/internal/config"
	"github.com/lockw1n/time-logger/internal/models"
)

var (
	ErrBadDateRange = errors.New("bad date range")
	ErrBadMonth     = errors.New("bad month")
	ErrInvalidHours = errors.New("invalid hours")
	ErrInvalidLabel = errors.New("invalid label")
	ErrInvalidDate  = errors.New("invalid date")
	ErrNotFound     = errors.New("not found")
)

const maxDailyHours = 24.0

type EntryService struct {
	DB            *gorm.DB
	allowedLabels []string
	allowedSet    map[string]struct{}
}

type MonthlyReport struct {
	Month      string         `json:"month"`
	Start      string         `json:"start"`
	End        string         `json:"end"`
	Items      []labelSummary `json:"items"`
	TotalHours float64        `json:"total_hours"`
}

type labelSummary struct {
	Label      string  `json:"label"`
	TotalHours float64 `json:"total_hours"`
}

func NewEntryService(db *gorm.DB) *EntryService {
	labels := config.AllowedLabels()
	set := make(map[string]struct{})
	for _, l := range labels {
		set[l] = struct{}{}
	}
	return &EntryService{
		DB:            db,
		allowedLabels: labels,
		allowedSet:    set,
	}
}

func (s *EntryService) AllowedLabelsError() string {
	return "label must be one of: " + strings.Join(s.allowedLabels, ", ")
}

// Helper validation/parsing
func isQuarterHour(hours float64) bool {
	quarters := hours * 4
	return math.Abs(quarters-math.Round(quarters)) < 1e-9
}

func (s *EntryService) validateHours(hours float64) error {
	if hours < 0 || hours > maxDailyHours || !isQuarterHour(hours) {
		return ErrInvalidHours
	}
	return nil
}

func (s *EntryService) validateLabel(label string) error {
	if label == "" {
		return ErrInvalidLabel
	}
	_, ok := s.allowedSet[strings.ToLower(strings.TrimSpace(label))]
	if !ok {
		return ErrInvalidLabel
	}
	return nil
}

func parseDateOnly(dateStr string) (time.Time, error) {
	parsed, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return time.Time{}, ErrInvalidDate
	}
	return time.Date(parsed.Year(), parsed.Month(), parsed.Day(), 0, 0, 0, 0, time.UTC), nil
}

func todayUTCDate() time.Time {
	now := time.Now().UTC()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
}

func parseDateRange(startStr, endStr string) (time.Time, time.Time, error) {
	start, err1 := time.Parse("2006-01-02", startStr)
	end, err2 := time.Parse("2006-01-02", endStr)
	if err1 != nil || err2 != nil {
		return time.Time{}, time.Time{}, ErrBadDateRange
	}
	end = end.Add(24 * time.Hour) // include end date
	return start, end, nil
}

func parseMonth(monthStr string) (time.Time, time.Time, error) {
	monthStr = strings.TrimSpace(monthStr)
	if monthStr == "" {
		return time.Time{}, time.Time{}, ErrBadMonth
	}

	parsed, err := time.Parse("2006-01", monthStr)
	if err != nil {
		return time.Time{}, time.Time{}, ErrBadMonth
	}
	start := time.Date(parsed.Year(), parsed.Month(), 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0) // first day of next month
	return start, end, nil
}

// Service methods

func (s *EntryService) ListEntries(startStr, endStr string) ([]models.Entry, error) {
	var entries []models.Entry
	query := s.DB.Model(&models.Entry{})

	if startStr != "" && endStr != "" {
		start, end, err := parseDateRange(startStr, endStr)
		if err != nil {
			return nil, err
		}
		query = query.Where("date >= ? AND date < ?", start, end)
	}

	if err := query.Order("ticket asc, date asc").Find(&entries).Error; err != nil {
		return nil, err
	}
	return entries, nil
}

func (s *EntryService) Summary(startStr, endStr string) ([]ticketSummary, error) {
	start, end, err := parseDateRange(startStr, endStr)
	if err != nil {
		return nil, err
	}

	var summaries []ticketSummary
	if err := s.DB.Model(&models.Entry{}).
		Select("ticket, COALESCE(MAX(label), '') as label, SUM(hours) as total_hours").
		Where("date >= ? AND date < ?", start, end).
		Group("ticket").
		Order("ticket asc").
		Scan(&summaries).Error; err != nil {
		return nil, err
	}

	return summaries, nil
}

func (s *EntryService) MonthlySummary(monthStr string) (*MonthlyReport, error) {
	start, end, err := parseMonth(monthStr)
	if err != nil {
		return nil, err
	}

	var items []labelSummary
	if err := s.DB.Model(&models.Entry{}).
		Select("COALESCE(label, '') as label, SUM(hours) as total_hours").
		Where("date >= ? AND date < ?", start, end).
		Group("label").
		Order("label asc").
		Scan(&items).Error; err != nil {
		return nil, err
	}

	total := 0.0
	for _, item := range items {
		total += item.TotalHours
	}

	return &MonthlyReport{
		Month:      start.Format("2006-01"),
		Start:      start.Format("2006-01-02"),
		End:        end.AddDate(0, 0, -1).Format("2006-01-02"),
		Items:      items,
		TotalHours: total,
	}, nil
}

func (s *EntryService) GetEntry(id string) (*models.Entry, error) {
	var entry models.Entry
	if err := s.DB.First(&entry, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &entry, nil
}

// CreateOrSum creates a new entry or sums hours if one already exists for the same ticket/date.
// Returns created=true for new entry, created=false for merged entry.
func (s *EntryService) CreateOrSum(input entryInput) (*models.Entry, bool, error) {
	if input.Ticket == "" {
		return nil, false, ErrInvalidHours
	}
	if err := s.validateHours(input.Hours); err != nil {
		return nil, false, err
	}
	if err := s.validateLabel(input.Label); err != nil {
		return nil, false, err
	}

	date := todayUTCDate()
	if input.Date != "" {
		var err error
		date, err = parseDateOnly(input.Date)
		if err != nil {
			return nil, false, err
		}
	}

	var existing models.Entry
	err := s.DB.Where("ticket = ? AND date = ?", input.Ticket, date).First(&existing).Error
	if err == nil {
		newHours := existing.Hours + input.Hours
		if err := s.validateHours(newHours); err != nil {
			return nil, false, err
		}
		existing.Hours = newHours
		if input.Label != "" {
			existing.Label = input.Label
		}
		existing.UpdatedAt = time.Now().UTC()
		if err := s.DB.Save(&existing).Error; err != nil {
			return nil, false, err
		}
		return &existing, false, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, false, err
	}

	entry := models.Entry{
		Ticket:    input.Ticket,
		Label:     input.Label,
		Hours:     input.Hours,
		Date:      date,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	if err := s.DB.Create(&entry).Error; err != nil {
		return nil, false, err
	}
	return &entry, true, nil
}

func (s *EntryService) UpdateEntry(id string, input entryInput) (*models.Entry, error) {
	entry, err := s.GetEntry(id)
	if err != nil {
		return nil, err
	}

	if input.Ticket != "" {
		entry.Ticket = input.Ticket
	}
	if input.Label != "" {
		if err := s.validateLabel(input.Label); err != nil {
			return nil, err
		}
		entry.Label = input.Label
	}

	if err := s.validateHours(input.Hours); err != nil {
		return nil, err
	}
	entry.Hours = input.Hours

	if input.Date != "" {
		parsed, err := parseDateOnly(input.Date)
		if err != nil {
			return nil, err
		}
		entry.Date = parsed
	}

	if err := s.DB.Save(entry).Error; err != nil {
		return nil, err
	}

	return entry, nil
}

func (s *EntryService) DeleteEntry(id string) error {
	if err := s.DB.Delete(&models.Entry{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (s *EntryService) DeleteByTicket(ticket string) error {
	return s.DB.Where("ticket = ?", ticket).Delete(&models.Entry{}).Error
}
