package shortener

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/phuslu/log"

	"github.com/mars-terminal/mechta/internal/domain"
	"github.com/mars-terminal/mechta/internal/service"
	"github.com/mars-terminal/mechta/internal/shared/ctx_tools"
	"github.com/mars-terminal/mechta/internal/storage"
)

const (
	defaultExpireDays = 30

	maxRetries = 10

	shortLinkMaxChars = 8
)

func (s *Service) CreateShortLink(ctx context.Context, cmd service.CreateLinkCMD) (domain.Link, error) {
	if err := validateURL(cmd.URL); err != nil {
		return domain.Link{}, fmt.Errorf("%w: %w", err, domain.ErrBadURL)
	}

	if cmd.ExpireDays <= 0 {
		cmd.ExpireDays = defaultExpireDays
	}

	for retries := 0; retries < maxRetries; retries++ {
		link, err := s.storage.CreateLink(ctx, storage.CreateLinkCMD{
			ID:        domain.NewLinkID(),
			TargetURL: cmd.URL,
			ShortLink: createShortUrl(uuid.NewString()),
			ExpireAt:  time.Now().AddDate(0, 0, cmd.ExpireDays),
		})
		if err != nil && !errors.Is(err, storage.ErrDuplicateShortURL) {
			return domain.Link{}, fmt.Errorf("failed to create link, %w", err)
		}

		if err == nil {
			link.ShortLink = s.baseURL + "/" + link.ShortLink
			return link, nil
		}

		ctx_tools.GetLogger(ctx, log.Info()).Err(err).Msg("there is collisions")
	}

	return domain.Link{}, service.ErrMaxRetriesReachedOnCreateLink
}

func (s *Service) GetLinks(ctx context.Context) ([]domain.Link, error) {
	links, err := s.storage.GetLinks(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get links: %w", err)
	}

	for i := range links {
		links[i].ShortLink = s.baseURL + "/" + links[i].ShortLink
	}

	return links, nil
}

func (s *Service) GetLinkStatistics(ctx context.Context, shortLink string) (domain.Link, error) {
	if err := validateShortLink(shortLink); err != nil {
		return domain.Link{}, fmt.Errorf("%w: %w", err, domain.ErrBadShortLink)
	}

	link, err := s.storage.GetRawLinkByShortLink(ctx, shortLink)
	if err != nil {
		return domain.Link{}, fmt.Errorf("failed to get link by short url: %w", err)
	}

	link.ShortLink = s.baseURL + "/" + link.ShortLink

	return link, nil
}

func (s *Service) RedirectLink(ctx context.Context, shortLink string) (domain.Link, error) {
	if err := validateShortLink(shortLink); err != nil {
		return domain.Link{}, fmt.Errorf("%w: %w", err, domain.ErrBadShortLink)
	}

	link, err := s.storage.GetLinkByShortLink(ctx, shortLink)
	if err != nil {
		return domain.Link{}, err
	}

	if link.DeletedAt != nil {
		return domain.Link{}, domain.ErrLinkDeleted
	}

	if err := s.storage.UpdateLinkByShortUrl(ctx, storage.UpdateLinkCMD{
		ID:          link.ID,
		LastAccess:  time.Now(),
		AccessCount: link.AccessCount + 1,
	}); err != nil {
		return domain.Link{}, err
	}

	return link, nil
}

func (s *Service) DeleteLink(ctx context.Context, shortURL string) error {
	if shortURL == "" {
		return domain.ErrBadURL
	}

	return s.storage.DeleteLinkByShortUrl(ctx, shortURL)
}

func validateURL(sourceURL string) error {
	if len(strings.TrimSpace(sourceURL)) == 0 {
		return fmt.Errorf("url cannot be empty")
	}

	u, err := url.ParseRequestURI(sourceURL)
	if err != nil {
		return fmt.Errorf("invalid url: %w", err)
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("url scheme must be http or https")
	}

	if u.Hostname() == "" {
		return fmt.Errorf("url hostname is empty")
	}

	if len(strings.Split(u.Hostname(), ".")) == 1 {
		return fmt.Errorf("url cannot be root domain")
	}

	return nil
}

func validateShortLink(shortLink string) error {
	if shortLink == "" {
		return fmt.Errorf("link cannot be empty, [%s]", shortLink)
	}

	if strings.Contains(shortLink, " ") {
		return fmt.Errorf("invalid short link has space, [%s]", shortLink)
	}

	if utf8.RuneCountInString(shortLink) != shortLinkMaxChars {
		return fmt.Errorf("short link length must be %d characters", shortLinkMaxChars)
	}

	if validateURL(shortLink) == nil {
		return fmt.Errorf("must be not URL: %s", shortLink)
	}

	return nil
}

func createShortUrl(uuid string) string {
	hash := sha256.Sum256([]byte(uuid))
	return base64.URLEncoding.EncodeToString(hash[:])[:shortLinkMaxChars-1] + string(uuid[len(uuid)-1])
}
