package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	_ "time/tzdata" // base de timezones embutida no binário

	"github.com/hugaojanuario/uptime-monitor-devops/internal/healthcheck"
	"github.com/hugaojanuario/uptime-monitor-devops/internal/http/handler"
	"github.com/hugaojanuario/uptime-monitor-devops/internal/http/router"
	"github.com/hugaojanuario/uptime-monitor-devops/internal/repository"
	"github.com/hugaojanuario/uptime-monitor-devops/internal/services"
	"github.com/hugaojanuario/uptime-monitor-devops/pkg/config"
	"github.com/hugaojanuario/uptime-monitor-devops/pkg/database"
)

// @title			Uptime Monitor API
// @version		1.0
// @description	API para cadastrar domínios/urls e verificar o status http de cada um. Cada verificação também é gravada no arquivo de resultados.
// @host			localhost:8080
// @BasePath		/
func main() {
	fmt.Println("Uptime Monitor starting...")

	cfg := config.LoadDotEnv()

	loc, err := time.LoadLocation(cfg.TZ)
	if err != nil {
		log.Fatalf("timezone inválida (%s): %v", cfg.TZ, err)
	}
	time.Local = loc

	db, err := database.Conn(database.Config{
		DB_HOST:     cfg.DB_HOST,
		DB_PORT:     cfg.DB_PORT,
		DB_USER:     cfg.DB_USER,
		DB_PASSWORD: cfg.DB_PASSWORD,
		DB_NAME:     cfg.DB_NAME,
		DB_SSLMODE:  cfg.DB_SSLMODE,
		TZ:          cfg.TZ,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	checker := healthcheck.NewChecker(cfg.CHECK_TIMEOUT, cfg.RESULTS_FILE)
	repo := repository.NewRepository(db)
	serv := services.NewService(repo, checker)
	url := handler.NewURLController(serv)
	router := router.SetupRouter(url)

	srv := &http.Server{
		Addr:         ":" + cfg.PORT,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("erro no servidor: %v", err)
		}
	}()
	log.Printf("Uptime Monitor ouvindo na porta %s (resultados em %s)", cfg.PORT, cfg.RESULTS_FILE)

	<-ctx.Done()
	log.Println("desligando servidor...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("erro no shutdown: %v", err)
	}

	log.Println("servidor encerrado")
}
