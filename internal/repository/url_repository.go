package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/hugaojanuario/uptime-monitor-devops/internal/models"
	"github.com/lib/pq"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateURL(req models.CreateURLRequest) (*models.URL, error) {
	query := `
        INSERT INTO urls (name, address)
        VALUES ($1, $2)
        RETURNING id, name, address, created_at
    `
	url := &models.URL{}
	err := r.db.QueryRow(query, req.Name, req.URL).
		Scan(&url.ID, &url.Name, &url.Address, &url.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("erro ao cadastrar a nova url: %w", err)
	}

	return url, nil
}

func (r *Repository) ListURLs() ([]models.URL, error) {
	query := `SELECT id, name, address, created_at
				FROM urls
				ORDER BY created_at`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar as urls: %w", err)
	}
	defer rows.Close()

	urls := []models.URL{}
	for rows.Next() {
		url := models.URL{}
		if err := rows.Scan(&url.ID, &url.Name, &url.Address, &url.CreatedAt); err != nil {
			return nil, fmt.Errorf("erro ao ler a url: %w", err)
		}
		urls = append(urls, url)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao listar as urls: %w", err)
	}

	return urls, nil
}

func (r *Repository) FindURLByID(id string) (*models.URL, error) {
	query := `SELECT id, name, address, created_at
				FROM urls
				WHERE id = $1`

	url := &models.URL{}
	err := r.db.QueryRow(query, id).
		Scan(&url.ID, &url.Name, &url.Address, &url.CreatedAt)
	if err == sql.ErrNoRows || isInvalidUUID(err) {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("erro ao buscar a url: %w", err)
	}

	return url, nil
}

// isInvalidUUID indica um id fora do formato uuid, tratado como url inexistente.
func isInvalidUUID(err error) bool {
	var pqErr *pq.Error
	return errors.As(err, &pqErr) && pqErr.Code == "22P02"
}
