package shortener

import (
	"context"
	"errors"
	"testing"
	"time"
	"unicode/utf8"

	"github.com/jackc/pgx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/mars-terminal/mechta/internal/domain"
	"github.com/mars-terminal/mechta/internal/service"
	"github.com/mars-terminal/mechta/internal/storage"
)

const baseURL = "https://example.com"

func TestService_CreateShortLink(t *testing.T) {
	t.Parallel()

	type args service.CreateLinkCMD

	type result struct {
		want *domain.Link
		err  error
	}

	tests := map[string]struct {
		setup  func() storage.Shortener
		args   args
		result result
	}{
		"happy path": {
			setup: func() storage.Shortener {
				shortenerStorage := storage.NewMockShortener(gomock.NewController(t))

				shortenerStorage.EXPECT().
					CreateLink(gomock.Any(), gomock.AssignableToTypeOf(storage.CreateLinkCMD{})).
					DoAndReturn(func(ctx context.Context, cmd storage.CreateLinkCMD) (domain.Link, error) {
						if cmd.ShortLink == "" {
							return domain.Link{}, errors.New("shortener short url is required")
						}
						if _, err := domain.ParseLinkID(cmd.ID.String()); err != nil {
							return domain.Link{}, err
						}

						if cmd.TargetURL != "https://google.com/1" {
							return domain.Link{}, errors.New("target url does not match")
						}

						if result := time.Now().AddDate(0, 0, 30).Sub(cmd.ExpireAt); result > time.Hour {
							return domain.Link{}, errors.New("expire days too far")
						}

						return domain.Link{
							ID:          "1",
							TargetUrl:   "https://google.com/1",
							ShortLink:   "short-url",
							LastAccess:  nil,
							AccessCount: 0,
							CreatedAt:   time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
							ExpireAt:    time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
							UpdatedAt:   time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
							DeletedAt:   nil,
						}, nil
					})

				return shortenerStorage
			},
			args: args{
				URL:        "https://google.com/1",
				ExpireDays: 30,
			},
			result: result{
				want: &domain.Link{
					ID:          "1",
					TargetUrl:   "https://google.com/1",
					ShortLink:   "short-url",
					LastAccess:  nil,
					AccessCount: 0,
					CreatedAt:   time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
					ExpireAt:    time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt:   time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
					DeletedAt:   nil,
				},
				err: nil,
			},
		},
		"max times already exists": {
			setup: func() storage.Shortener {
				shortenerStorage := storage.NewMockShortener(gomock.NewController(t))

				shortenerStorage.EXPECT().
					CreateLink(gomock.Any(), gomock.AssignableToTypeOf(storage.CreateLinkCMD{})).
					Times(maxRetries).
					DoAndReturn(func(ctx context.Context, cmd storage.CreateLinkCMD) (domain.Link, error) {
						return domain.Link{}, storage.ErrDuplicateShortURL
					})

				return shortenerStorage
			},
			args: args{
				URL:        "https://google.com/1",
				ExpireDays: 30,
			},
			result: result{
				want: &domain.Link{},
				err:  service.ErrMaxRetriesReachedOnCreateLink,
			},
		},
		"max - 3 times already exists": {
			setup: func() storage.Shortener {
				shortenerStorage := storage.NewMockShortener(gomock.NewController(t))

				shortenerStorage.EXPECT().
					CreateLink(gomock.Any(), gomock.AssignableToTypeOf(storage.CreateLinkCMD{})).
					Times(maxRetries - 3).
					DoAndReturn(func(ctx context.Context, cmd storage.CreateLinkCMD) (domain.Link, error) {
						return domain.Link{}, storage.ErrDuplicateShortURL
					})

				shortenerStorage.EXPECT().
					CreateLink(gomock.Any(), gomock.AssignableToTypeOf(storage.CreateLinkCMD{})).
					DoAndReturn(func(ctx context.Context, cmd storage.CreateLinkCMD) (domain.Link, error) {
						return domain.Link{
							ID:          "1",
							TargetUrl:   "https://google.com/1",
							ShortLink:   "short-url",
							LastAccess:  nil,
							AccessCount: 0,
							CreatedAt:   time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
							ExpireAt:    time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
							UpdatedAt:   time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
							DeletedAt:   nil,
						}, nil
					})

				return shortenerStorage
			},
			args: args{
				URL:        "https://google.com/1",
				ExpireDays: 30,
			},
			result: result{
				want: &domain.Link{
					ID:          "1",
					TargetUrl:   "https://google.com/1",
					ShortLink:   "short-url",
					LastAccess:  nil,
					AccessCount: 0,
					CreatedAt:   time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
					ExpireAt:    time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt:   time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
					DeletedAt:   nil,
				},
				err: nil,
			},
		},
		"failed to validate url": {
			setup: func() storage.Shortener {
				shortenerStorage := storage.NewMockShortener(gomock.NewController(t))

				return shortenerStorage
			},
			args: args{
				URL:        "htts://google./me/.com/1",
				ExpireDays: 30,
			},
			result: result{
				want: &domain.Link{},
				err:  domain.ErrBadURL,
			},
		},
	}

	for nn, tc := range tests {
		nn, tc := nn, tc

		t.Run(nn, func(t *testing.T) {
			t.Parallel()

			s := NewService(baseURL, tc.setup())

			link, err := s.CreateShortLink(context.Background(), service.CreateLinkCMD(tc.args))
			if tc.result.err == nil {
				require.NoError(t, err)
			} else {
				require.ErrorIs(t, err, tc.result.err)
			}

			if tc.result.want == nil {
				assert.Nil(t, link)
			} else {
				assert.Equal(t, *tc.result.want, link)
			}
		})
	}
}

func TestService_DeleteLink(t *testing.T) {
	t.Parallel()

	type result struct {
		err error
	}

	tests := map[string]struct {
		setup  func() storage.Shortener
		args   string
		result result
	}{
		"happy path": {
			setup: func() storage.Shortener {
				shortenerStorage := storage.NewMockShortener(gomock.NewController(t))

				shortenerStorage.EXPECT().
					DeleteLinkByShortUrl(gomock.Any(), "test_url").
					Return(nil)

				return shortenerStorage
			},
			args: "test_url",
			result: result{
				err: nil,
			},
		},
		"not found": {
			setup: func() storage.Shortener {
				shortenerStorage := storage.NewMockShortener(gomock.NewController(t))

				shortenerStorage.EXPECT().
					DeleteLinkByShortUrl(gomock.Any(), "test_url").
					Return(domain.ErrNotFound)

				return shortenerStorage
			},
			args: "test_url",
			result: result{
				err: domain.ErrNotFound,
			},
		},
		"deleted": {
			setup: func() storage.Shortener {
				shortenerStorage := storage.NewMockShortener(gomock.NewController(t))

				shortenerStorage.EXPECT().
					DeleteLinkByShortUrl(gomock.Any(), "test_url").
					Return(domain.ErrLinkDeleted)

				return shortenerStorage
			},
			args: "test_url",
			result: result{
				err: domain.ErrLinkDeleted,
			},
		},
		"bad url": {
			setup: func() storage.Shortener {
				shortenerStorage := storage.NewMockShortener(gomock.NewController(t))

				shortenerStorage.EXPECT().
					DeleteLinkByShortUrl(gomock.Any(), gomock.Any()).
					Return(domain.ErrBadURL)

				return shortenerStorage
			},
			args: "htts://mamska11.az",
			result: result{
				err: domain.ErrBadURL,
			},
		},
	}

	for nn, tc := range tests {
		nn, tc := nn, tc

		t.Run(nn, func(t *testing.T) {
			t.Parallel()
			s := NewService(baseURL, tc.setup())

			err := s.DeleteLink(context.Background(), tc.args)
			if tc.result.err == nil {
				require.NoError(t, err)
			} else {
				require.ErrorIs(t, err, tc.result.err)
			}
		})
	}
}

func TestService_GetLinkStatistics(t *testing.T) {
	t.Parallel()

	type result struct {
		want *domain.Link
		err  error
	}

	tests := map[string]struct {
		setup  func() storage.Shortener
		args   string
		result result
	}{
		"happy path": {
			setup: func() storage.Shortener {
				shortenerStorage := storage.NewMockShortener(gomock.NewController(t))

				shortenerStorage.EXPECT().GetRawLinkByShortLink(gomock.Any(), "12345678").
					DoAndReturn(func(ctx context.Context, shortLink string) (domain.Link, error) {
						if shortLink != "12345678" {
							return domain.Link{}, errors.New("short url does not match")
						}

						return domain.Link{
							ID:          "1",
							TargetUrl:   "https://google.com/1",
							ShortLink:   shortLink,
							LastAccess:  nil,
							AccessCount: 0,
							CreatedAt:   time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
							ExpireAt:    time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
							UpdatedAt:   time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
							DeletedAt:   nil,
						}, nil
					})

				return shortenerStorage
			},
			args: "12345678",
			result: result{
				want: &domain.Link{
					ID:          "1",
					TargetUrl:   "https://google.com/1",
					ShortLink:   "12345678",
					LastAccess:  nil,
					AccessCount: 0,
					CreatedAt:   time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
					ExpireAt:    time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt:   time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
					DeletedAt:   nil,
				},
				err: nil,
			},
		},
		"not found": {
			setup: func() storage.Shortener {
				shortenerStorage := storage.NewMockShortener(gomock.NewController(t))

				shortenerStorage.EXPECT().GetRawLinkByShortLink(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, shortLink string) (domain.Link, error) {
						return domain.Link{}, domain.ErrNotFound
					})

				return shortenerStorage
			},
			args: "12345678",
			result: result{
				want: &domain.Link{},
				err:  domain.ErrNotFound,
			},
		},
		"bad url": {
			setup: func() storage.Shortener {
				shortenerStorage := storage.NewMockShortener(gomock.NewController(t))

				return shortenerStorage
			},
			args: "htttps:/ada1231.akz",
			result: result{
				want: &domain.Link{},
				err:  domain.ErrBadShortLink,
			},
		},
	}

	for nn, tc := range tests {
		nn, tc := nn, tc

		t.Run(nn, func(t *testing.T) {
			t.Parallel()

			s := NewService(baseURL, tc.setup())

			link, err := s.GetLinkStatistics(context.Background(), tc.args)
			if tc.result.err == nil {
				require.NoError(t, err)
			} else {
				require.ErrorIs(t, err, tc.result.err)
			}

			if tc.result.want == nil {
				assert.Nil(t, link)
			} else {
				assert.Equal(t, *tc.result.want, link)
			}
		})
	}
}

func TestService_GetLinks(t *testing.T) {
	t.Parallel()

	type result struct {
		want *[]domain.Link
		err  error
	}

	tests := map[string]struct {
		setup  func() storage.Shortener
		result result
	}{
		"happy path": {
			setup: func() storage.Shortener {
				shortenerStorage := storage.NewMockShortener(gomock.NewController(t))

				shortenerStorage.EXPECT().GetLinks(gomock.Any()).
					DoAndReturn(func(ctx context.Context) ([]domain.Link, error) {
						return []domain.Link{
							{
								ID:          "1",
								TargetUrl:   "https://google.com/1",
								ShortLink:   "short-url",
								LastAccess:  nil,
								AccessCount: 0,
								CreatedAt:   time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
								ExpireAt:    time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
								UpdatedAt:   time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
								DeletedAt:   nil,
							},
						}, nil
					})

				return shortenerStorage
			},
			result: result{
				want: &[]domain.Link{
					{
						ID:          "1",
						TargetUrl:   "https://google.com/1",
						ShortLink:   "short-url",
						LastAccess:  nil,
						AccessCount: 0,
						CreatedAt:   time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
						ExpireAt:    time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
						UpdatedAt:   time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
						DeletedAt:   nil,
					},
				},
				err: nil,
			},
		},
		"database not active": {
			setup: func() storage.Shortener {
				shortenerStorage := storage.NewMockShortener(gomock.NewController(t))

				shortenerStorage.EXPECT().GetLinks(gomock.Any()).
					DoAndReturn(func(ctx context.Context) ([]domain.Link, error) {
						return nil, pgx.ErrDeadConn
					})

				return shortenerStorage
			},
			result: result{
				want: nil,
				err:  pgx.ErrDeadConn,
			},
		},
	}

	for nn, tc := range tests {
		nn, tc := nn, tc

		t.Run(nn, func(t *testing.T) {
			t.Parallel()

			s := NewService(baseURL, tc.setup())

			link, err := s.GetLinks(context.Background())
			if tc.result.err == nil {
				require.NoError(t, err)
			} else {
				require.ErrorIs(t, err, tc.result.err)
			}

			if tc.result.want == nil {
				assert.Nil(t, link)
			} else {
				assert.Equal(t, *tc.result.want, link)
			}
		})
	}
}

func TestService_RedirectLink(t *testing.T) {
	t.Parallel()

	type args struct {
		link string
	}

	type result struct {
		want *domain.Link
		err  error
	}

	tests := map[string]struct {
		setup  func() storage.Shortener
		args   args
		result result
	}{
		"happy path": {
			setup: func() storage.Shortener {
				shortenerStorage := storage.NewMockShortener(gomock.NewController(t))

				shortenerStorage.EXPECT().GetLinkByShortLink(gomock.Any(), "12345678").
					DoAndReturn(func(ctx context.Context, shortLink string) (domain.Link, error) {
						if shortLink != "12345678" {
							return domain.Link{}, errors.New("short url does not match")
						}

						return domain.Link{
							ID:          "1",
							TargetUrl:   "https://google.com/1",
							ShortLink:   shortLink,
							LastAccess:  nil,
							AccessCount: 0,
							CreatedAt:   time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
							ExpireAt:    time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
							UpdatedAt:   time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
							DeletedAt:   nil,
						}, nil
					})

				shortenerStorage.EXPECT().
					UpdateLinkByShortUrl(gomock.Any(), gomock.AssignableToTypeOf(storage.UpdateLinkCMD{})).
					DoAndReturn(func(ctx context.Context, cmd storage.UpdateLinkCMD) error {
						if cmd.ID != "1" {
							return errors.New("short url does not match")
						}

						if date := cmd.LastAccess.Sub(time.Now()); date > time.Millisecond {
							return errors.New("last access date is null")
						}

						if cmd.AccessCount <= 0 {
							return errors.New("access count is negative value")
						}

						return nil
					})

				return shortenerStorage
			},
			args: args{
				link: "12345678",
			},
			result: result{
				want: &domain.Link{
					ID:          "1",
					TargetUrl:   "https://google.com/1",
					ShortLink:   "12345678",
					LastAccess:  nil,
					AccessCount: 0,
					CreatedAt:   time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
					ExpireAt:    time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt:   time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
					DeletedAt:   nil,
				},
				err: nil,
			},
		},
		"not found": {
			setup: func() storage.Shortener {
				shortenerStorage := storage.NewMockShortener(gomock.NewController(t))

				shortenerStorage.EXPECT().GetLinkByShortLink(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, shortLink string) (domain.Link, error) {
						return domain.Link{}, domain.ErrNotFound
					})

				return shortenerStorage
			},
			args: args{
				link: "12345678",
			},
			result: result{
				want: &domain.Link{},
				err:  domain.ErrNotFound,
			},
		},
		"deleted": {
			setup: func() storage.Shortener {
				shortenerStorage := storage.NewMockShortener(gomock.NewController(t))

				shortenerStorage.EXPECT().GetLinkByShortLink(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, shortLink string) (domain.Link, error) {
						return domain.Link{}, domain.ErrLinkDeleted
					})

				return shortenerStorage
			},
			args: args{
				link: "12345678",
			},
			result: result{
				want: &domain.Link{},
				err:  domain.ErrLinkDeleted,
			},
		},
	}

	for nn, tc := range tests {
		nn, tc := nn, tc

		t.Run(nn, func(t *testing.T) {
			t.Parallel()

			s := NewService(baseURL, tc.setup())

			link, err := s.RedirectLink(context.Background(), tc.args.link)
			if tc.result.err == nil {
				require.NoError(t, err)
			} else {
				require.ErrorIs(t, err, tc.result.err)
			}

			if tc.result.want == nil {
				assert.Nil(t, link)
			} else {
				assert.Equal(t, *tc.result.want, link)
			}
		})
	}
}

func Test_validateURL(t *testing.T) {
	tests := map[string]struct {
		err bool
	}{
		"https":                  {err: true},
		"https://":               {err: true},
		"":                       {err: true},
		"http://www":             {err: true},
		"http://www.example.com": {err: false},
		"http://example.com":     {err: false},
		"https://example.com":    {err: false},
		"htts://example.com":     {err: true},
		"ftp://example.com":      {err: true},
		"dns://example.com":      {err: true},
		"https://www.example.com?somecoolquery=1": {err: false},
		"https://www.example.com:443":             {err: false},
		"/testing-path":                           {err: true},
		"testing-path":                            {err: true},
		"alskjff#?asf//dfas":                      {err: true},
	}
	for nn, tc := range tests {
		nn, tc := nn, tc

		t.Run(nn, func(t *testing.T) {
			t.Parallel()

			err := validateURL(nn)
			switch tc.err {
			case true:
				require.Error(t, err)
			case false:
				require.NoError(t, err)
			}
		})
	}
}

func Test_validateShortLink(t *testing.T) {
	type args struct {
		shortLink string
	}
	tests := map[string]struct {
		args  args
		error bool
	}{
		"happy path": {
			args: args{
				shortLink: "12345678",
			},
			error: false,
		},
		"length error": {
			args: args{
				shortLink: "shortLink909",
			},
			error: true,
		},
		"space error": {
			args: args{
				shortLink: "0 491 213",
			},
			error: true,
		},
		"empty error": {
			args: args{
				shortLink: "",
			},
			error: true,
		},
		"uri": {
			args: args{
				shortLink: "https://example.com",
			},
			error: true,
		},
	}
	for nn, tc := range tests {
		nn, tc := nn, tc

		t.Run(nn, func(t *testing.T) {
			t.Parallel()

			err := validateShortLink(tc.args.shortLink)
			switch tc.error {
			case true:
				require.Error(t, err)
			case false:
				require.NoError(t, err)
			}
		})
	}
}

func Test_createShortUrl(t *testing.T) {
	t.Parallel()

	type args struct {
		uuid string
	}
	tests := map[string]struct {
		args args
		want string
	}{
		"happy path": {
			args: args{
				uuid: "550e8400-e29b-41d4-a716-446655440000",
			},
			want: "o6nh7Zc0",
		},
		"uuid 1": {
			args: args{
				uuid: "1b315fea-e802-48c2-80fc-808953f2ce04",
			},
			want: "gXJ8jsO4",
		},
		"uuid 2": {
			args: args{
				uuid: "80d572c6-fc98-4bb0-aa02-0bb4ecf1dde2",
			},
			want: "6lelAOF2",
		},
		"uuid 3": {
			args: args{
				uuid: "bb609228-f91d-46fa-9add-0c2a4e20a5c0",
			},
			want: "uNREqbz0",
		},
	}
	for nn, tc := range tests {
		nn, tc := nn, tc

		t.Run(nn, func(t *testing.T) {
			t.Parallel()

			result := createShortUrl(tc.args.uuid)

			require.Equal(t, utf8.RuneCountInString(tc.want), shortLinkMaxChars)
			require.Equal(t, tc.want, result)
		})
	}
}
