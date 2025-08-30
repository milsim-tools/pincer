package units

import (
	"github.com/milsim-tools/pincer/internal/models"
	unitsv1 "github.com/milsim-tools/pincer/pkg/api/gen/milsimtools/units/v1"
)

type UnitsUnit struct {
	models.Model

	DisplayName string `gorm:"notNull"`
	Slug        string `gorm:"notNull;uniqueIndex"`
	Description string `gorm:"notNull"`
	OwnerID     string `gorm:"notNull"`
}

func (u UnitsUnit) Proto() *unitsv1.Unit {
	return &unitsv1.Unit{
		Id:          u.ID,
		DisplayName: u.DisplayName,
		Slug:        u.Slug,
		Description: u.Description,
		OwnerId:     u.OwnerID,
	}
}
