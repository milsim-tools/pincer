package users

import (
	"github.com/milsim-tools/pincer/internal/models"
	usersv1 "github.com/milsim-tools/pincer/pkg/api/gen/milsimtools/users/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UsersUser struct {
	models.Model

	SSOID       string `gorm:"notNull;uniqueIndex"`
	DisplayName string `gorm:"notNull"`
	Username    string `gorm:"notNull;uniqueIndex"`
	Email       string `gorm:"notNull;uniqueIndex"`
	Bio         string `gorm:"type:text"`
	AvatarURL   string
}

func (u UsersUser) Proto() *usersv1.User {
	return &usersv1.User{
		Id:          u.ID,
		DisplayName: u.DisplayName,
		Username:    u.Username,
		Email:       u.Email,
		Bio:         u.Bio,
		AvatarUrl:   u.AvatarURL,
		CreatedAt:   timestamppb.New(u.CreatedAt),
		UpdatedAt:   timestamppb.New(u.UpdatedAt),
	}
}
