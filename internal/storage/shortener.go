package storage

import (
	"context"
	"errors"
	"time"

	"github.com/mars-terminal/mechta/internal/domain"
)

var ErrDuplicateShortURL = errors.New("short url already exists")

type CreateLinkCMD struct {
	ID        domain.LinkID
	TargetURL string
	ShortLink string
	ExpireAt  time.Time
}

type UpdateLinkCMD struct {
	ID          domain.LinkID
	LastAccess  time.Time
	AccessCount uint64
}

//go:generate mockgen -source=shortener.go -destination shortener_mock.gen.go -package storage
type Shortener interface {
	CreateLink(ctx context.Context, cmd CreateLinkCMD) (domain.Link, error)

	GetLinks(ctx context.Context) ([]domain.Link, error)

	GetLinkByShortLink(ctx context.Context, shortURL string) (domain.Link, error)

	GetRawLinkByShortLink(ctx context.Context, shortURL string) (domain.Link, error)

	UpdateLinkByShortUrl(ctx context.Context, cmd UpdateLinkCMD) error

	DeleteLinkByShortUrl(ctx context.Context, shortURL string) error
}
