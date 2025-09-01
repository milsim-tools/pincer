package helpers

import (
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/milsim-tools/pincer/internal/models"
	"gorm.io/gorm"
)

const (
	DefaultPageSize = 50
	MaxPageSize     = 100
)

type PaginationCursor struct {
	CreatedAt time.Time
}

func (pc PaginationCursor) String() (string, error) {
	bytes, err := json.Marshal(pc)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// GetPageLimit ensures the page size is within acceptable bounds.
func GetPageLimit(pageSize int) int {
	limit := int(pageSize)
	if limit <= 0 {
		limit = DefaultPageSize
	} else if limit > MaxPageSize {
		limit = MaxPageSize
	}
	return limit
}

type Model interface {
	models.Model
}

func ApplyPageLimit[T any](iface gorm.ChainInterface[T], pageSize int) gorm.ChainInterface[T] {
	limit := GetPageLimit(pageSize)
	qb := iface.Limit(limit)
	return qb
}

func ApplyCursor[T any](iface gorm.ChainInterface[T], pc *PaginationCursor) gorm.ChainInterface[T] {
	qb := iface.Where("created_at < ?", pc.CreatedAt)
	return qb
}

func GenerateCursor(items []models.Model) *PaginationCursor {
	if len(items) == 0 {
		return nil
	}
	pc := &PaginationCursor{
		CreatedAt: items[len(items)-1].CreatedAt,
	}
	return pc
}

func GenerateCursorString(items []models.Model) string {
	pc := GenerateCursor(items)
	if val, err := pc.String(); err != nil {
		return ""
	} else {
		return val
	}
}

func CursorFromString(s string) (*PaginationCursor, error) {
	if s == "" {
		return nil, nil
	}

	var pc PaginationCursor
	bytes, err := base64.URLEncoding.DecodeString(s)
	if err != nil {
		return &pc, err
	}
	err = json.Unmarshal(bytes, &pc)
	return &pc, err
}
