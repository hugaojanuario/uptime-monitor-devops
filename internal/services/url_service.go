package services

import (
	"context"
	"errors"
	"net/url"
	"strings"

	"github.com/hugaojanuario/uptime-monitor-devops/internal/healthcheck"
	"github.com/hugaojanuario/uptime-monitor-devops/internal/models"
	"github.com/hugaojanuario/uptime-monitor-devops/internal/repository"
)

var (
	ErrURLNotFound      = errors.New("url não encontrada")
	ErrInvalidURL       = errors.New("url inválida, use http:// ou https://")
	ErrURLAlreadyExists = errors.New("url já cadastrada")
)

type Service struct {
	r *repository.Repository
	c *healthcheck.Checker
}

func NewService(r *repository.Repository, c *healthcheck.Checker) *Service {
	return &Service{r: r, c: c}
}

func (s *Service) CreateURL(req models.CreateURLRequest) (*models.URL, error) {
	req.Name = strings.TrimSpace(req.Name)
	req.URL = strings.TrimSpace(req.URL)

	if !isValidURL(req.URL) {
		return nil, ErrInvalidURL
	}

	created, err := s.r.CreateURL(req)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return nil, ErrURLAlreadyExists
		}
		return nil, err
	}

	return created, nil
}

// CheckAllURLs verifica todas as urls cadastradas.
func (s *Service) CheckAllURLs(ctx context.Context) ([]models.CheckResponse, error) {
	urls, err := s.r.ListURLs()
	if err != nil {
		return nil, err
	}

	return s.c.CheckAll(ctx, urls), nil
}

// CheckURLByID verifica apenas a url do id informado.
func (s *Service) CheckURLByID(ctx context.Context, id string) (*models.CheckResponse, error) {
	existing, err := s.r.FindURLByID(id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, ErrURLNotFound
	}

	result := s.c.Check(ctx, *existing)

	return &result, nil
}

func isValidURL(raw string) bool {
	parsed, err := url.Parse(raw)
	if err != nil {
		return false
	}

	return (parsed.Scheme == "http" || parsed.Scheme == "https") && parsed.Host != ""
}
