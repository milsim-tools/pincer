package helpers

import (
	"encoding/base64"
	"encoding/json"
)

type PaginationCursor struct {
	CreatedAt string
}

func (pc PaginationCursor) String() (string, error) {
	bytes, err := json.Marshal(pc)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

func CursorFromString(s string) (PaginationCursor, error) {
	var pc PaginationCursor
	bytes, err := base64.URLEncoding.DecodeString(s)
	if err != nil {
		return pc, err
	}
	err = json.Unmarshal(bytes, &pc)
	return pc, err
}
