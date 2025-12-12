package label

import (
	labeldto "github.com/lockw1n/time-logger/internal/dto/label"
	"github.com/lockw1n/time-logger/internal/models"
)

/*************** DTO → Model ***************/

func ToModel(d labeldto.Request) *models.Label {
	return &models.Label{
		CompanyID: d.CompanyID,
		Name:      d.Name,
		Color:     d.Color,
	}
}

/*************** Model → DTO ***************/

func ToResponse(m *models.Label) labeldto.Response {
	if m == nil {
		color := "gray"
		return labeldto.Response{
			ID:    0,
			Name:  "(missing label)",
			Color: &color,
		}
	}

	return labeldto.Response{
		ID:        m.ID,
		CompanyID: m.CompanyID,
		Name:      m.Name,
		Color:     m.Color,
	}
}

func ToResponses(list []models.Label) []labeldto.Response {
	out := make([]labeldto.Response, len(list))
	for i := range list {
		out[i] = ToResponse(&list[i])
	}
	return out
}
