package shortener

import (
	"context"
	"errors"
	"net/http"

	api "github.com/mars-terminal/mechta/api/gen"
	"github.com/mars-terminal/mechta/internal/domain"
	"github.com/mars-terminal/mechta/internal/service"
)

var _ api.StrictServerInterface = (*Handlers)(nil)

type Handlers struct {
	service service.Shortener
}

func NewHandlers(service service.Shortener) *Handlers {
	return &Handlers{service: service}
}

func (h *Handlers) GetShortener(ctx context.Context, request api.GetShortenerRequestObject) (api.GetShortenerResponseObject, error) {
	links, err := h.service.GetLinks(ctx)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			return api.GetShortener404JSONResponse{
				Code:    http.StatusNotFound,
				Message: domain.ErrNotFound.Error(),
			}, nil
		}

		return api.GetShortener500JSONResponse{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		}, nil
	}

	var result = make(api.GetShortener200JSONResponse, len(links))
	for i := range links {
		result[i] = api.LinkItem{
			AccessCount: int(links[i].AccessCount),
			CreatedAt:   links[i].CreatedAt,
			DeletedAt:   links[i].DeletedAt,
			ExpireAt:    links[i].ExpireAt,
			Id:          links[i].ID.String(),
			LastAccess:  links[i].LastAccess,
			ShortLink:   links[i].ShortLink,
			TargetUrl:   links[i].TargetUrl,
			UpdatedAt:   links[i].UpdatedAt,
		}
	}

	return result, nil
}

func (h *Handlers) PostShortener(ctx context.Context, request api.PostShortenerRequestObject) (api.PostShortenerResponseObject, error) {
	link, err := h.service.CreateShortLink(ctx, service.CreateLinkCMD{
		URL:        request.Body.Url,
		ExpireDays: request.Body.ExpireDays,
	})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrBadURL):
			return api.PostShortener400JSONResponse{
				Code:    http.StatusBadRequest,
				Message: domain.ErrBadURL.Error(),
			}, nil

		case errors.Is(err, service.ErrMaxRetriesReachedOnCreateLink):
			return api.PostShortener400JSONResponse{
				Code:    http.StatusBadRequest,
				Message: service.ErrMaxRetriesReachedOnCreateLink.Error(),
			}, nil
		}

		return api.PostShortener500JSONResponse{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		}, nil
	}

	return api.PostShortener200JSONResponse{
		ShortLink: link.ShortLink,
	}, nil
}

func (h *Handlers) GetStatsLink(ctx context.Context, request api.GetStatsLinkRequestObject) (api.GetStatsLinkResponseObject, error) {
	link, err := h.service.GetLinkStatistics(ctx, request.Link)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			return api.GetStatsLink404JSONResponse{
				Code:    http.StatusNotFound,
				Message: domain.ErrNotFound.Error(),
			}, nil
		}

		return api.GetStatsLink500JSONResponse{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		}, nil
	}

	return api.GetStatsLink200JSONResponse{
		AccessCount: int(link.AccessCount),
		CreatedAt:   link.CreatedAt,
		DeletedAt:   link.DeletedAt,
		ExpireAt:    link.ExpireAt,
		Id:          link.ID.String(),
		LastAccess:  link.LastAccess,
		ShortLink:   link.ShortLink,
		TargetUrl:   link.TargetUrl,
		UpdatedAt:   link.UpdatedAt,
	}, nil
}

func (h *Handlers) DeleteLink(ctx context.Context, request api.DeleteLinkRequestObject) (api.DeleteLinkResponseObject, error) {
	if err := h.service.DeleteLink(ctx, request.Link); err != nil {
		switch {
		case errors.Is(err, domain.ErrLinkDeleted):
			return api.DeleteLink404JSONResponse{
				Code:    http.StatusNotFound,
				Message: domain.ErrNotFound.Error(),
			}, nil
		case errors.Is(err, domain.ErrNotFound):
			return api.DeleteLink404JSONResponse{
				Code:    http.StatusNotFound,
				Message: domain.ErrNotFound.Error(),
			}, nil
		}

		return api.DeleteLink500JSONResponse{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		}, nil
	}
	return api.DeleteLink200JSONResponse{
		Code:    http.StatusOK,
		Message: "success",
	}, nil
}

func (h *Handlers) GetLink(ctx context.Context, request api.GetLinkRequestObject) (api.GetLinkResponseObject, error) {
	link, err := h.service.RedirectLink(ctx, request.Link)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			return api.GetLink404JSONResponse{
				Code:    http.StatusNotFound,
				Message: domain.ErrNotFound.Error(),
			}, nil
		}

		return api.GetLink500JSONResponse{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		}, nil
	}

	return api.GetLink302Response{
		Headers: api.GetLink302ResponseHeaders{
			Location: link.TargetUrl,
		},
	}, nil
}
