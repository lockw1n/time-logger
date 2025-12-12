package entry

import (
	"errors"

	entrydto "github.com/lockw1n/time-logger/internal/dto/entry"
	entrymapper "github.com/lockw1n/time-logger/internal/mapper/entry"
	"github.com/lockw1n/time-logger/internal/models"
	assignmentrepo "github.com/lockw1n/time-logger/internal/repository/consultantassignment"
	entryrepo "github.com/lockw1n/time-logger/internal/repository/entry"
	ticketrepo "github.com/lockw1n/time-logger/internal/repository/ticket"
	assignmentservice "github.com/lockw1n/time-logger/internal/service/consultantassignment"
)

type service struct {
	entryRepo      entryrepo.Repository
	assignmentRepo assignmentrepo.Repository
	ticketRepo     ticketrepo.Repository
}

func NewService(
	entryRepo entryrepo.Repository,
	assignmentRepo assignmentrepo.Repository,
	ticketRepo ticketrepo.Repository,
) Service {
	return &service{
		entryRepo:      entryRepo,
		assignmentRepo: assignmentRepo,
		ticketRepo:     ticketRepo,
	}
}

func (s *service) Create(req entrydto.Request) (*entrydto.Response, error) {
	if err := validateDurationMinutesRange(req.DurationMinutes); err != nil {
		return nil, err
	}
	if err := validateDurationMinutesQuarter(req.DurationMinutes); err != nil {
		return nil, err
	}
	if err := validateDateFormat(req.Date); err != nil {
		return nil, err
	}

	assignment, err := s.assignmentRepo.FindByPair(req.ConsultantID, req.CompanyID)
	if errors.Is(err, assignmentrepo.ErrNotFound) {
		return nil, assignmentservice.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	ticket, err := s.ticketRepo.FindByCode(req.CompanyID, req.TicketCode)
	if errors.Is(err, ticketrepo.ErrNotFound) {
		ticket = &models.Ticket{
			Code:      req.TicketCode,
			CompanyID: req.CompanyID,
		}

		_, err := s.ticketRepo.Create(ticket)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	model, err := entrymapper.ToModel(req)
	if err != nil {
		return nil, err
	}

	model.ConsultantAssignmentID = assignment.ID
	model.HourlyRateSnapshot = &assignment.HourlyRate
	model.CurrencySnapshot = &assignment.Currency
	model.TicketID = &ticket.ID

	created, err := s.entryRepo.Create(model)

	if err != nil {
		return nil, err
	}

	out := entrymapper.ToResponse(created)
	return &out, nil
}

func (s *service) Update(id uint64, req entrydto.Request) (*entrydto.Response, error) {
	if err := validateDurationMinutesRange(req.DurationMinutes); err != nil {
		return nil, err
	}
	if err := validateDurationMinutesQuarter(req.DurationMinutes); err != nil {
		return nil, err
	}
	if err := validateDateFormat(req.Date); err != nil {
		return nil, err
	}

	existing, err := s.entryRepo.FindByID(id)

	if err != nil {
		if errors.Is(err, entryrepo.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	model, err := entrymapper.ToModel(req)

	if err != nil {
		return nil, err
	}

	model.ID = existing.ID
	updated, err := s.entryRepo.Update(model)

	if err != nil {
		return nil, err
	}

	out := entrymapper.ToResponse(updated)
	return &out, nil
}

func (s *service) Delete(id uint64) error {
	err := s.entryRepo.Delete(id)
	if err != nil {
		if errors.Is(err, entryrepo.ErrNotFound) {
			return ErrNotFound
		}
		return err
	}
	return nil
}

func (s *service) Get(id uint64) (*entrydto.Response, error) {
	model, err := s.entryRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, entryrepo.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	out := entrymapper.ToResponse(model)
	return &out, nil
}

func (s *service) ListByCompany(companyID uint64) ([]entrydto.Response, error) {
	list, err := s.entryRepo.FindByCompany(companyID)
	if err != nil {
		return nil, err
	}

	result := make([]entrydto.Response, len(list))
	for i := range list {
		result[i] = entrymapper.ToResponse(&list[i])
	}
	return result, nil
}

func (s *service) ListByConsultant(consultantID uint64) ([]entrydto.Response, error) {
	list, err := s.entryRepo.FindByConsultant(consultantID)
	if err != nil {
		return nil, err
	}

	result := make([]entrydto.Response, len(list))
	for i := range list {
		result[i] = entrymapper.ToResponse(&list[i])
	}
	return result, nil
}

func (s *service) ListAll() ([]entrydto.Response, error) {
	list, err := s.entryRepo.ListAll()
	if err != nil {
		return nil, err
	}

	result := make([]entrydto.Response, len(list))
	for i := range list {
		result[i] = entrymapper.ToResponse(&list[i])
	}
	return result, nil
}
