package shortener

import (
	"github.com/mars-terminal/mechta/internal/storage"
)

type Service struct {
	baseURL string
	storage storage.Shortener
}

func NewService(baseURL string, storage storage.Shortener) *Service {
	return &Service{
		baseURL: baseURL,
		storage: storage,
	}
}
