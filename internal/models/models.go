package models

import "time"

type Model struct {
	ID string `gorm:"primaryKey"`
  CreatedAt    time.Time
  UpdatedAt    time.Time
}
