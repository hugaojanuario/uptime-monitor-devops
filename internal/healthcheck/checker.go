package healthcheck

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/hugaojanuario/uptime-monitor-devops/internal/models"
)

type Checker struct {
	client      *http.Client
	resultsFile string
	mu          sync.Mutex
}

func NewChecker(timeout time.Duration, resultsFile string) *Checker {
	return &Checker{
		client:      &http.Client{Timeout: timeout},
		resultsFile: resultsFile,
	}
}

// Check faz um GET na url e grava o retorno no arquivo de resultados.
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

	if err := c.writeResult(result); err != nil {
		log.Printf("[checker] erro ao gravar o resultado em %s: %v", c.resultsFile, err)
	}

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

func (c *Checker) writeResult(result models.CheckResponse) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	file, err := os.OpenFile(c.resultsFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("erro ao abrir o arquivo de resultados: %w", err)
	}
	defer file.Close()

	status := fmt.Sprintf("%d", result.StatusCode)
	if result.Error != "" {
		status = "ERROR: " + result.Error
	}

	line := fmt.Sprintf("%s\tid=%s\turl=%s\tstatus=%s\tduration=%dms\n",
		result.CheckedAt.Format(time.RFC3339), result.ID, result.Address, status, result.DurationMs)

	if _, err := file.WriteString(line); err != nil {
		return fmt.Errorf("erro ao gravar o resultado: %w", err)
	}

	return nil
}
