package members

import (
	"github.com/milsim-tools/pincer/internal/models"
	membersv1 "github.com/milsim-tools/pincer/pkg/api/gen/milsimtools/members/v1"
)

type MembersUnitMember struct {
	models.Model

	UnitID      string `gorm:"notNull"`
	UserID      string `gorm:"notNull"`
	Permissions int32  `gorm:"notNull"`
	Status      int32  `gorm:"notNull;default=1"`
}

func (u MembersUnitMember) Proto() *membersv1.UnitMember {
	return &membersv1.UnitMember{
		UnitId:      u.UnitID,
		UserId:      u.UserID,
		Permissions: u.Permissions,
	}
}
