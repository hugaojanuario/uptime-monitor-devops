# uptime-monitor-devops

> **O foco deste repositório é DevOps, não a aplicação.**

Projeto de portfólio para construir e demonstrar toda a esteira de infraestrutura e operação de um serviço
real: **Terraform** (provisionamento), **Ansible** (configuração), **Kubernetes** (orquestração),
**AWS** (cloud) e **CI/CD**, além de observabilidade, segurança e automação.

A API de monitoramento de domínios existe para ser a *workload* dessa esteira — um serviço pequeno,
com estado (Postgres) e dependência de rede externa, o suficiente para exercitar deploy, escalonamento,
health checks, secrets, persistência e monitoramento de verdade. Ela é o meio, não o fim.

## A aplicação

API REST em Go (Gin) para cadastrar URLs e verificar o status HTTP delas. Postgres roda via Docker Compose
e as migrations são aplicadas pelo `migrate/migrate`. Cada verificação também é gravada em um arquivo `.txt`.

## Subindo

```bash
cp .env.example .env   # ajuste DB_USER / DB_PASSWORD
docker compose up -d --build
```

- API: http://localhost:8080
- Postgres: localhost:5432
- Resultados: `./data/results.txt` (montado como `/data` no container)

## Endpoints

| Método | Rota         | O que faz                                                        |
|--------|--------------|------------------------------------------------------------------|
| GET    | `/health`    | Healthcheck                                                       |
| POST   | `/urls`      | Cadastra uma URL                                                  |
| GET    | `/urls`      | Verifica **todas** as URLs cadastradas e retorna os status HTTP   |
| GET    | `/urls/:id`  | Verifica **apenas** a URL do id informado                         |
| GET    | `/swagger/index.html` | Documentação interativa (Swagger UI)                    |

Todo `GET` de verificação escreve as linhas correspondentes no `results.txt`.

### Swagger

Com a stack no ar, a documentação fica em **http://localhost:8080/swagger/index.html**
(spec em `/swagger/doc.json`). Os arquivos gerados ficam em `docs/` e são compilados no binário.

Depois de mexer nas anotações dos controllers ou nos models, regere:

```bash
go install github.com/swaggo/swag/cmd/swag@latest   # uma vez
swag init -g cmd/api/main.go
```

### Exemplos

```bash
# cadastrar
curl -X POST localhost:8080/urls -d '{"name":"google","url":"https://www.google.com"}'
# {"id":"84c142da-...","name":"google","url":"https://www.google.com","created_at":"..."}

# verificar todas
curl localhost:8080/urls
# [{"id":"84c142da-...","status_code":200,"duration_ms":478,"checked_at":"..."}]

# verificar uma
curl localhost:8080/urls/84c142da-dd0f-4664-bd43-0e32ff1c5ef7
```

URLs que falham (DNS, timeout, conexão) voltam com `status_code: 0` e o campo `error` preenchido.

### Formato do results.txt

```
2026-07-19T22:48:45-03:00	id=84c142da-...	url=https://www.google.com	status=200	duration=478ms
2026-07-19T22:48:45-03:00	id=21882adc-...	url=https://x.invalido	status=ERROR: dial tcp: lookup ...	duration=277ms
```

## Variáveis de ambiente

Ver `.env.example`: `PORT`, `DB_HOST`, `DB_PORT`, `DB_NAME`, `DB_USER`, `DB_PASSWORD`, `DB_SSLMODE`,
`RESULTS_FILE`, `CHECK_TIMEOUT`, `TZ`.

Os horários (`created_at`, `checked_at` e o `results.txt`) usam o fuso de **Brasília**
(`TZ=America/Sao_Paulo`, padrão). A base de timezones vem embutida no binário (`time/tzdata`),
então não depende do sistema operacional da imagem.

## Rodando local (sem container da API)

```bash
docker compose up -d postgres migrate
go run ./cmd/api
```

## Estrutura

```
cmd/api/                 entrypoint (config, conexão, servidor, graceful shutdown)
pkg/config/              carga da .env
pkg/database/            conexão com o Postgres
internal/models/         entidade URL e DTOs
internal/repository/     queries SQL
internal/services/       regras de negócio (validação, orquestração dos checks)
internal/healthcheck/    execução dos checks HTTP e escrita do .txt
internal/http/handler/   controllers Gin
internal/http/router/    rotas e CORS
internal/utils/          tradução de erro para status HTTP
db/migrations/           migrations (golang-migrate)
docs/                    spec do Swagger gerada pelo swag
```

## Roadmap DevOps

A parte que realmente importa aqui:

- [ ] **Terraform** — VPC, subnets, security groups, EKS/EC2, RDS, ECR, IAM (state remoto no S3 + lock)
- [ ] **Ansible** — provisionamento e hardening dos nós, instalação de agentes
- [ ] **Kubernetes** — Deployment, Service, Ingress, ConfigMap/Secret, HPA, probes, Job de migration
- [ ] **AWS** — ECR, EKS, RDS, ALB, Route53, CloudWatch
- [ ] **CI/CD** — build, testes, lint, scan de imagem e deploy automatizado
- [ ] **Observabilidade** — Prometheus, Grafana e logs centralizados
- [ ] **Segurança** — gestão de secrets, least privilege no IAM, imagem sem root

