package database

import (
	"database/sql"
	"errors"
	"fmt"
	"io/fs"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

// Migrate aplica as migrations pendentes usando a conexão já aberta. O driver
// do postgres pega um advisory lock, então réplicas subindo ao mesmo tempo não
// disputam o schema: uma migra e as outras esperam.
func Migrate(db *sql.DB, migrations fs.FS, dir string) error {
	source, err := iofs.New(migrations, dir)
	if err != nil {
		return fmt.Errorf("erro ao ler as migrations embutidas: %w", err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("erro ao preparar o driver de migration: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", source, "postgres", driver)
	if err != nil {
		return fmt.Errorf("erro ao inicializar o migrate: %w", err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("[migrate] schema já está atualizado")
			return nil
		}
		return fmt.Errorf("erro ao aplicar as migrations: %w", err)
	}

	version, _, err := m.Version()
	if err != nil {
		return fmt.Errorf("erro ao ler a versão do schema: %w", err)
	}

	log.Printf("[migrate] migrations aplicadas, schema na versão %d", version)

	return nil
}
