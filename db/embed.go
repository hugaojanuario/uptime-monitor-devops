// Package db embute os arquivos de migration no binário, para que o schema
// viaje junto com a versão da aplicação e não dependa do filesystem do host.
package db

import "embed"

//go:embed migrations/*.sql
var Migrations embed.FS
