package models

import "time"

type URL struct {
	ID        string    `json:"id" example:"84c142da-dd0f-4664-bd43-0e32ff1c5ef7"`
	Name      string    `json:"name" example:"google"`
	Address   string    `json:"url" example:"https://www.google.com"`
	CreatedAt time.Time `json:"created_at" example:"2026-07-20T01:16:23Z"`
}

// DTOs
type CreateURLRequest struct {
	Name string `json:"name" example:"google"`
	URL  string `json:"url" example:"https://www.google.com"`
}

type CheckResponse struct {
	ID         string    `json:"id" example:"84c142da-dd0f-4664-bd43-0e32ff1c5ef7"`
	Name       string    `json:"name" example:"google"`
	Address    string    `json:"url" example:"https://www.google.com"`
	StatusCode int       `json:"status_code" example:"200"`
	Error      string    `json:"error,omitempty" example:""`
	DurationMs int64     `json:"duration_ms" example:"478"`
	CheckedAt  time.Time `json:"checked_at" example:"2026-07-20T01:17:58Z"`
}

type ErrorResponse struct {
	Error string `json:"error" example:"url não encontrada"`
}

type HealthResponse struct {
	Status string `json:"status" example:"ok"`
}
