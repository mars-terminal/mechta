package service

import (
	"context"
	"errors"

	"github.com/mars-terminal/mechta/internal/domain"
)

var ErrMaxRetriesReachedOnCreateLink = errors.New("max retries reached on create link")

type CreateLinkCMD struct {
	URL        string
	ExpireDays int
}

//go:generate mockgen -source=shortener.go -destination shortener_mock.gen.go -package service
type Shortener interface {
	CreateShortLink(ctx context.Context, cmd CreateLinkCMD) (domain.Link, error)

	GetLinks(ctx context.Context) ([]domain.Link, error)
	
	GetLinkStatistics(ctx context.Context, shortLink string) (domain.Link, error)

	RedirectLink(ctx context.Context, shortLink string) (domain.Link, error)

	DeleteLink(ctx context.Context, shortLink string) error
}
