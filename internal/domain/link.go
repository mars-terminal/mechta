package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrBadURL       = errors.New("bad origin link URL")
	ErrBadShortLink = errors.New("bad short link")
	ErrNotFound     = errors.New("not found")
	ErrLinkDeleted  = errors.New("link is deleted")
)

type LinkID string

func NewLinkID() LinkID {
	return LinkID(uuid.NewString())
}

func (id LinkID) String() string {
	return string(id)
}

func ParseLinkID(id string) (LinkID, error) {
	_, err := uuid.Parse(id)
	if err != nil {
		return "", err
	}
	return LinkID(id), nil
}

type Link struct {
	ID          LinkID
	TargetUrl   string
	ShortLink   string
	LastAccess  *time.Time
	AccessCount uint64
	CreatedAt   time.Time
	ExpireAt    time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}
