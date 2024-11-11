package shortener

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	api "github.com/mars-terminal/mechta/api/gen"
	"github.com/mars-terminal/mechta/internal/domain"
	"github.com/mars-terminal/mechta/internal/service"
)

func TestHandlers_DeleteLink(t *testing.T) {
	t.Parallel()

	type result struct {
		want api.DeleteLinkResponseObject
		err  error
	}

	tests := map[string]struct {
		setup  func() service.Shortener
		result result
	}{
		"happy path": {
			setup: func() service.Shortener {
				shortenerService := service.NewMockShortener(gomock.NewController(t))

				shortenerService.EXPECT().
					DeleteLink(gomock.Any(), "short-url").
					DoAndReturn(func(ctx context.Context, shortURL string) error {
						return nil
					})

				return shortenerService
			},
			result: result{
				want: api.DeleteLink200JSONResponse{
					Code:    http.StatusOK,
					Message: "success",
				},
				err: nil,
			},
		},
		"deleted": {
			setup: func() service.Shortener {
				shortenerService := service.NewMockShortener(gomock.NewController(t))

				shortenerService.EXPECT().
					DeleteLink(gomock.Any(), "short-url").
					DoAndReturn(func(ctx context.Context, shortURL string) error {
						return domain.ErrLinkDeleted
					})

				return shortenerService
			},
			result: result{
				want: api.DeleteLink404JSONResponse{
					Code:    http.StatusNotFound,
					Message: domain.ErrNotFound.Error(),
				},
				err: nil,
			},
		},
		"not found": {
			setup: func() service.Shortener {
				shortenerService := service.NewMockShortener(gomock.NewController(t))

				shortenerService.EXPECT().
					DeleteLink(gomock.Any(), "short-url").
					DoAndReturn(func(ctx context.Context, shortURL string) error {
						return domain.ErrNotFound
					})

				return shortenerService
			},
			result: result{
				want: api.DeleteLink404JSONResponse{
					Code:    http.StatusNotFound,
					Message: domain.ErrNotFound.Error(),
				},
				err: nil,
			},
		},
		"internal server error": {
			setup: func() service.Shortener {
				shortenerService := service.NewMockShortener(gomock.NewController(t))

				shortenerService.EXPECT().
					DeleteLink(gomock.Any(), "short-url").
					DoAndReturn(func(ctx context.Context, shortURL string) error {
						return fmt.Errorf("internal server error")
					})

				return shortenerService
			},
			result: result{
				want: api.DeleteLink500JSONResponse{
					Code:    http.StatusInternalServerError,
					Message: "internal server error",
				},
				err: nil,
			},
		},
	}

	for nn, tc := range tests {
		nn, tc := nn, tc

		t.Run(nn, func(t *testing.T) {
			t.Parallel()

			s := NewHandlers(tc.setup())

			link, err := s.DeleteLink(context.Background(), api.DeleteLinkRequestObject{
				Link: "short-url",
			})
			if tc.result.err == nil {
				require.NoError(t, err)
			} else {
				require.ErrorIs(t, err, tc.result.err)
			}

			if tc.result.want == nil {
				assert.Nil(t, link)
			} else {
				assert.Equal(t, tc.result.want, link)
			}
		})
	}
}

func TestHandlers_GetLink(t *testing.T) {
	t.Parallel()

	type result struct {
		want api.GetLinkResponseObject
		err  error
	}

	tests := map[string]struct {
		setup  func() service.Shortener
		result result
	}{
		"happy path": {
			setup: func() service.Shortener {
				shortenerService := service.NewMockShortener(gomock.NewController(t))

				shortenerService.EXPECT().
					RedirectLink(gomock.Any(), "short-url").
					DoAndReturn(func(ctx context.Context, shortURL string) (domain.Link, error) {
						if shortURL != "short-url" {
							return domain.Link{}, errors.New("url does not match")
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

				return shortenerService
			},
			result: result{
				want: api.GetLink302Response{
					Headers: api.GetLink302ResponseHeaders{
						Location: "https://google.com/1",
					},
				},
				err: nil,
			},
		},
		"not found": {
			setup: func() service.Shortener {
				shortenerService := service.NewMockShortener(gomock.NewController(t))

				shortenerService.EXPECT().
					RedirectLink(gomock.Any(), "short-url").
					DoAndReturn(func(ctx context.Context, shortURL string) (domain.Link, error) {
						return domain.Link{}, domain.ErrNotFound
					})

				return shortenerService
			},
			result: result{
				want: api.GetLink404JSONResponse{
					Code:    http.StatusNotFound,
					Message: domain.ErrNotFound.Error(),
				},
				err: nil,
			},
		},
		"internal server error": {
			setup: func() service.Shortener {
				shortenerService := service.NewMockShortener(gomock.NewController(t))

				shortenerService.EXPECT().
					RedirectLink(gomock.Any(), "short-url").
					DoAndReturn(func(ctx context.Context, shortURL string) (domain.Link, error) {
						return domain.Link{}, fmt.Errorf("internal server error")
					})

				return shortenerService
			},
			result: result{
				want: api.GetLink500JSONResponse{
					Code:    http.StatusInternalServerError,
					Message: "internal server error",
				},
				err: nil,
			},
		},
	}

	for nn, tc := range tests {
		nn, tc := nn, tc

		t.Run(nn, func(t *testing.T) {
			t.Parallel()

			s := NewHandlers(tc.setup())

			link, err := s.GetLink(context.Background(), api.GetLinkRequestObject{
				Link: "short-url",
			})
			if tc.result.err == nil {
				require.NoError(t, err)
			} else {
				require.ErrorIs(t, err, tc.result.err)
			}

			if tc.result.want == nil {
				assert.Nil(t, link)
			} else {
				assert.Equal(t, tc.result.want, link)
			}
		})
	}
}

func TestHandlers_GetShortener(t *testing.T) {
	t.Parallel()

	type result struct {
		want api.GetShortenerResponseObject
		err  error
	}

	tests := map[string]struct {
		setup  func() service.Shortener
		result result
	}{
		"happy path": {
			setup: func() service.Shortener {
				shortenerService := service.NewMockShortener(gomock.NewController(t))

				shortenerService.EXPECT().
					GetLinks(gomock.Any()).
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

				return shortenerService
			},
			result: result{
				want: api.GetShortener200JSONResponse{
					api.LinkItem{
						Id:          "1",
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
		"not found": {
			setup: func() service.Shortener {
				shortenerService := service.NewMockShortener(gomock.NewController(t))

				shortenerService.EXPECT().
					GetLinks(gomock.Any()).
					DoAndReturn(func(ctx context.Context) ([]domain.Link, error) {
						return []domain.Link{}, domain.ErrNotFound
					})

				return shortenerService
			},
			result: result{
				want: api.GetShortener404JSONResponse{
					Code:    http.StatusNotFound,
					Message: domain.ErrNotFound.Error(),
				},
				err: nil,
			},
		},
		"internal server error": {
			setup: func() service.Shortener {
				shortenerService := service.NewMockShortener(gomock.NewController(t))

				shortenerService.EXPECT().
					GetLinks(gomock.Any()).
					DoAndReturn(func(ctx context.Context) ([]domain.Link, error) {
						return []domain.Link{}, fmt.Errorf("internal server error")
					})

				return shortenerService
			},
			result: result{
				want: api.GetShortener500JSONResponse{
					Code:    http.StatusInternalServerError,
					Message: "internal server error",
				},
				err: nil,
			},
		},
	}

	for nn, tc := range tests {
		nn, tc := nn, tc

		t.Run(nn, func(t *testing.T) {
			t.Parallel()

			s := NewHandlers(tc.setup())

			link, err := s.GetShortener(context.Background(), api.GetShortenerRequestObject{})
			if tc.result.err == nil {
				require.NoError(t, err)
			} else {
				require.ErrorIs(t, err, tc.result.err)
			}

			if tc.result.want == nil {
				assert.Nil(t, link)
			} else {
				assert.Equal(t, tc.result.want, link)
			}
		})
	}
}

func TestHandlers_GetStatsLink(t *testing.T) {
	t.Parallel()

	type result struct {
		want api.GetStatsLinkResponseObject
		err  error
	}

	tests := map[string]struct {
		setup  func() service.Shortener
		result result
	}{
		"happy path": {
			setup: func() service.Shortener {
				shortenerService := service.NewMockShortener(gomock.NewController(t))

				shortenerService.EXPECT().
					GetLinkStatistics(gomock.Any(), "short-url").
					DoAndReturn(func(ctx context.Context, shortURL string) (domain.Link, error) {
						if shortURL != "short-url" {
							return domain.Link{}, errors.New("url does not match")
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

				return shortenerService
			},
			result: result{
				want: api.GetStatsLink200JSONResponse{
					Id:          "1",
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
		"not found": {
			setup: func() service.Shortener {
				shortenerService := service.NewMockShortener(gomock.NewController(t))

				shortenerService.EXPECT().
					GetLinkStatistics(gomock.Any(), "short-url").
					DoAndReturn(func(ctx context.Context, shortURL string) (domain.Link, error) {
						return domain.Link{}, domain.ErrNotFound
					})

				return shortenerService
			},
			result: result{
				want: api.GetStatsLink404JSONResponse{
					Code:    http.StatusNotFound,
					Message: domain.ErrNotFound.Error(),
				},
				err: nil,
			},
		},
		"internal server error": {
			setup: func() service.Shortener {
				shortenerService := service.NewMockShortener(gomock.NewController(t))

				shortenerService.EXPECT().
					GetLinkStatistics(gomock.Any(), "short-url").
					DoAndReturn(func(ctx context.Context, shortURL string) (domain.Link, error) {
						return domain.Link{}, fmt.Errorf("internal server error")
					})

				return shortenerService
			},
			result: result{
				want: api.GetStatsLink500JSONResponse{
					Code:    http.StatusInternalServerError,
					Message: "internal server error",
				},
				err: nil,
			},
		},
	}

	for nn, tc := range tests {
		nn, tc := nn, tc

		t.Run(nn, func(t *testing.T) {
			t.Parallel()

			s := NewHandlers(tc.setup())

			link, err := s.GetStatsLink(context.Background(), api.GetStatsLinkRequestObject{
				Link: "short-url",
			})
			if tc.result.err == nil {
				require.NoError(t, err)
			} else {
				require.ErrorIs(t, err, tc.result.err)
			}

			if tc.result.want == nil {
				assert.Nil(t, link)
			} else {
				assert.Equal(t, tc.result.want, link)
			}
		})
	}
}

func TestHandlers_PostShortener(t *testing.T) {
	t.Parallel()

	type args service.CreateLinkCMD

	type result struct {
		want api.PostShortenerResponseObject
		err  error
	}

	tests := map[string]struct {
		setup  func() service.Shortener
		args   args
		result result
	}{
		"happy path": {
			setup: func() service.Shortener {
				shortenerService := service.NewMockShortener(gomock.NewController(t))

				shortenerService.EXPECT().
					CreateShortLink(gomock.Any(), gomock.AssignableToTypeOf(service.CreateLinkCMD{})).
					DoAndReturn(func(ctx context.Context, cmd service.CreateLinkCMD) (domain.Link, error) {
						if cmd.URL != "https://google.com/1" {
							return domain.Link{}, errors.New("url does not match")
						}

						if cmd.ExpireDays < 0 {
							return domain.Link{}, errors.New("expire days is negative")
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

				return shortenerService
			},
			args: args{
				URL:        "https://google.com/1",
				ExpireDays: 30,
			},
			result: result{
				want: api.PostShortener200JSONResponse{
					ShortLink: "short-url",
				},
				err: nil,
			},
		},
		"bad url": {
			setup: func() service.Shortener {
				shortenerService := service.NewMockShortener(gomock.NewController(t))

				shortenerService.EXPECT().
					CreateShortLink(gomock.Any(), gomock.AssignableToTypeOf(service.CreateLinkCMD{})).
					DoAndReturn(func(ctx context.Context, cmd service.CreateLinkCMD) (domain.Link, error) {
						return domain.Link{}, domain.ErrBadURL
					})

				return shortenerService
			},
			args: args{
				URL:        "asfdjklfasd;ljafsd;jklafsd;jl",
				ExpireDays: 30,
			},
			result: result{
				want: api.PostShortener400JSONResponse{
					Code:    http.StatusBadRequest,
					Message: domain.ErrBadURL.Error(),
				},
				err: nil,
			},
		},
		"bad url 2": {
			setup: func() service.Shortener {
				shortenerService := service.NewMockShortener(gomock.NewController(t))

				shortenerService.EXPECT().
					CreateShortLink(gomock.Any(), gomock.AssignableToTypeOf(service.CreateLinkCMD{})).
					DoAndReturn(func(ctx context.Context, cmd service.CreateLinkCMD) (domain.Link, error) {
						return domain.Link{}, domain.ErrBadURL
					})

				return shortenerService
			},
			args: args{
				URL:        "ftp://google.com",
				ExpireDays: 30,
			},
			result: result{
				want: api.PostShortener400JSONResponse{
					Code:    http.StatusBadRequest,
					Message: domain.ErrBadURL.Error(),
				},
				err: nil,
			},
		},
		"bad url 3": {
			setup: func() service.Shortener {
				shortenerService := service.NewMockShortener(gomock.NewController(t))

				shortenerService.EXPECT().
					CreateShortLink(gomock.Any(), gomock.AssignableToTypeOf(service.CreateLinkCMD{})).
					DoAndReturn(func(ctx context.Context, cmd service.CreateLinkCMD) (domain.Link, error) {
						return domain.Link{}, domain.ErrBadURL
					})

				return shortenerService
			},
			args: args{
				URL:        "https://testcontainer s.com/guides/getting-started-with-testcontainers-for-go/",
				ExpireDays: 30,
			},
			result: result{
				want: api.PostShortener400JSONResponse{
					Code:    http.StatusBadRequest,
					Message: domain.ErrBadURL.Error(),
				},
				err: nil,
			},
		},
		"service internal error": {
			setup: func() service.Shortener {
				shortenerService := service.NewMockShortener(gomock.NewController(t))

				shortenerService.EXPECT().
					CreateShortLink(gomock.Any(), gomock.AssignableToTypeOf(service.CreateLinkCMD{})).
					DoAndReturn(func(ctx context.Context, cmd service.CreateLinkCMD) (domain.Link, error) {
						return domain.Link{}, fmt.Errorf("some internal error")
					})

				return shortenerService
			},
			args: args{
				URL:        "https://google.com/1",
				ExpireDays: 30,
			},
			result: result{
				want: api.PostShortener500JSONResponse{
					Code:    http.StatusInternalServerError,
					Message: "internal server error",
				},
				err: nil,
			},
		},
	}

	for nn, tc := range tests {
		nn, tc := nn, tc

		t.Run(nn, func(t *testing.T) {
			t.Parallel()

			s := NewHandlers(tc.setup())

			link, err := s.PostShortener(context.Background(), api.PostShortenerRequestObject{
				Body: &api.PostShortenerJSONRequestBody{
					ExpireDays: 30,
					Url:        "https://google.com/1",
				},
			})
			if tc.result.err == nil {
				require.NoError(t, err)
			} else {
				require.ErrorIs(t, err, tc.result.err)
			}

			if tc.result.want == nil {
				assert.Nil(t, link)
			} else {
				assert.Equal(t, tc.result.want, link)
			}
		})
	}
}
