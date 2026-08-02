package healthcheck

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/hugaojanuario/uptime-monitor-devops/internal/models"
)

type Checker struct {
	client *http.Client
}

func NewChecker(timeout time.Duration) *Checker {
	return &Checker{
		client: &http.Client{Timeout: timeout},
	}
}

// Check faz um GET na url e emite o retorno no log.
func (c *Checker) Check(ctx context.Context, url models.URL) models.CheckResponse {
	result := models.CheckResponse{
		ID:        url.ID,
		Name:      url.Name,
		Address:   url.Address,
		CheckedAt: time.Now(),
	}

	start := time.Now()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url.Address, nil)
	if err != nil {
		result.Error = err.Error()
	} else {
		resp, err := c.client.Do(req)
		if err != nil {
			result.Error = err.Error()
		} else {
			resp.Body.Close()
			result.StatusCode = resp.StatusCode
		}
	}
	result.DurationMs = time.Since(start).Milliseconds()

	logResult(result)

	return result
}

// CheckAll verifica todas as urls em paralelo, mantendo a ordem da entrada.
func (c *Checker) CheckAll(ctx context.Context, urls []models.URL) []models.CheckResponse {
	log.Printf("[checker] verificando %d urls", len(urls))

	results := make([]models.CheckResponse, len(urls))
	sem := make(chan struct{}, 5)
	var wg sync.WaitGroup

	for i, url := range urls {
		wg.Add(1)
		go func(i int, url models.URL) {
			defer wg.Done()

			select {
			case <-ctx.Done():
				return
			case sem <- struct{}{}:
			}
			defer func() { <-sem }()

			results[i] = c.Check(ctx, url)
		}(i, url)
	}

	wg.Wait()

	return results
}

// logResult emite o resultado do check no log do processo. Quem coleta, roteia
// e retém esse log é a plataforma, não a aplicação.
func logResult(result models.CheckResponse) {
	status := fmt.Sprintf("%d", result.StatusCode)
	if result.Error != "" {
		status = "ERROR: " + result.Error
	}

	log.Printf("[checker] id=%s url=%s status=%s duration=%dms",
		result.ID, result.Address, status, result.DurationMs)
}
