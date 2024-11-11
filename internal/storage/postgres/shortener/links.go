package shortener

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx"

	"github.com/mars-terminal/mechta/internal/domain"
	"github.com/mars-terminal/mechta/internal/storage"
)

type link struct {
	ID          domain.LinkID `db:"id"`
	TargetUrl   string        `db:"target_url"`
	ShortLink   string        `db:"short_link"`
	LastAccess  *time.Time    `db:"last_access"`
	AccessCount uint64        `db:"access_count"`
	CreatedAt   time.Time     `db:"created_at"`
	ExpireAt    time.Time     `db:"expire_at"`
	UpdatedAt   time.Time     `db:"updated_at"`
	DeletedAt   *time.Time    `db:"deleted_at"`
}

func (s *Storage) CreateLink(ctx context.Context, cmd storage.CreateLinkCMD) (domain.Link, error) {
	row := s.storage.QueryRowxContext(
		ctx,
		`INSERT INTO 
    		   		links
    		   		(id, target_url, short_link, expire_at)
			   VALUES
			        ($1, $2, $3, $4)
			   RETURNING id, short_link
	        `,
		cmd.ID,
		cmd.TargetURL,
		cmd.ShortLink,
		cmd.ExpireAt,
	)
	if err := row.Err(); err != nil {
		var e pgx.PgError
		if errors.As(err, &e) && e.Code == pgerrcode.UniqueViolation && e.ColumnName == "short_link" {
			return domain.Link{}, fmt.Errorf("short url is already exists: %w", storage.ErrDuplicateShortURL)
		}
		return domain.Link{}, fmt.Errorf("failed to insert link: %w", err)
	}

	var result link
	if err := row.StructScan(&result); err != nil {
		return domain.Link{}, err
	}

	return mapLinkToDomain(result), nil
}

func (s *Storage) GetLinkByShortLink(ctx context.Context, shortLink string) (domain.Link, error) {
	link, err := s.GetRawLinkByShortLink(ctx, shortLink)
	if err != nil {
		return domain.Link{}, err
	}

	if link.DeletedAt != nil {
		return domain.Link{}, domain.ErrLinkDeleted
	}

	return link, nil
}

func (s *Storage) GetRawLinkByShortLink(ctx context.Context, shortLink string) (domain.Link, error) {
	row := s.storage.QueryRowxContext(
		ctx,
		`select * from links where short_link = $1`,
		shortLink,
	)
	if err := row.Err(); err != nil {
		return domain.Link{}, fmt.Errorf("failed to get rows: %w", err)
	}

	var result link
	if err := row.StructScan(&result); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Link{}, fmt.Errorf("no rows: %w", domain.ErrNotFound)
		}
		return domain.Link{}, fmt.Errorf("failed to scan: %w", err)
	}

	return mapLinkToDomain(result), nil
}

func (s *Storage) GetLinks(ctx context.Context) ([]domain.Link, error) {
	rows, err := s.storage.QueryxContext(ctx, `select * from links order by created_at desc`)
	if err != nil {
		return nil, fmt.Errorf("failed to get rows: %w", err)
	}

	var result = make([]domain.Link, 0)
	for rows.Next() {
		var l link
		if err := rows.StructScan(&l); err != nil {
			return nil, fmt.Errorf("failed to scan: %w", err)
		}

		result = append(result, mapLinkToDomain(l))
	}

	if err := rows.Close(); err != nil {
		return nil, fmt.Errorf("failed to close rows: %w", err)
	}

	return result, nil
}

func (s *Storage) UpdateLinkByShortUrl(ctx context.Context, cmd storage.UpdateLinkCMD) error {
	_, err := s.storage.ExecContext(
		ctx,
		`update links set last_access = $1, access_count = $2 where id = $3 and deleted_at is null`,
		cmd.LastAccess,
		cmd.AccessCount,
		cmd.ID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("failed to get rows: %w", domain.ErrNotFound)
		}
		return fmt.Errorf("failed to delete row: %w", err)
	}

	return nil
}

func (s *Storage) DeleteLinkByShortUrl(ctx context.Context, shortLink string) error {
	_, err := s.storage.ExecContext(
		ctx,
		`update links set deleted_at = now() where short_link = $1 and deleted_at is null`,
		shortLink,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("failed to get rows: %w", domain.ErrNotFound)
		}
		return fmt.Errorf("failed to delete row: %w", err)
	}

	return nil
}

func mapLinkToDomain(l link) domain.Link {
	return domain.Link{
		ID:          l.ID,
		TargetUrl:   l.TargetUrl,
		ShortLink:   l.ShortLink,
		LastAccess:  l.LastAccess,
		AccessCount: l.AccessCount,
		CreatedAt:   l.CreatedAt,
		ExpireAt:    l.ExpireAt,
		UpdatedAt:   l.UpdatedAt,
		DeletedAt:   l.DeletedAt,
	}
}
