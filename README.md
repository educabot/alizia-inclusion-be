# Alizia Inclusion Backend

API REST para la plataforma de inclusion educativa Alizia. Gestiona dispositivos adaptativos, perfiles de estudiantes, adaptaciones curriculares y recomendaciones de IA para docentes.

## Arquitectura

Clean Architecture con 4 capas y dependency injection manual:

```
cmd/                            Bootstrap, DI wiring, HTTP server
  |
src/entrypoints/                Handlers HTTP (Gin) + Middleware (JWT, Tenant)
  |
src/core/usecases/              Logica de negocio (patron Execute)
  |
src/core/providers/             Interfaces (contratos)
  |                               ^
src/repositories/               Implementaciones (GORM, Azure OpenAI)
  |
src/core/entities/              Modelos de dominio
```

Multi-tenancy: todas las queries estan scoped por `organization_id` (UUID), extraido del JWT.

## Tech Stack

- **Go 1.26** con Gin (HTTP) y GORM (ORM)
- **PostgreSQL 16** con migraciones SQL
- **Azure OpenAI** para recomendaciones de IA
- **JWT RS256** para autenticacion (delegada a auth-service)
- **Docker Compose** para desarrollo local

## Prerequisitos

- Go 1.26.x
- Docker y Docker Compose
- `GITHUB_TOKEN` con acceso a repos de Educabot (modulo privado `team-ai-toolkit`)
- `GOPRIVATE=github.com/educabot/*`

## Setup

```bash
# 1. Levantar DB y correr migraciones
docker compose up -d

# 2. Seed data (opcional)
make seed

# 3. Correr el servidor (requiere 'air' para hot-reload)
make run

# O build y ejecutar directamente
make build
./alizia-inclusion-api
```

## Variables de Entorno

| Variable | Requerida | Default | Descripcion |
|---|---|---|---|
| `DATABASE_URL` | Si | - | Connection string de PostgreSQL |
| `PORT` | No | `8080` | Puerto del servidor HTTP |
| `ENV` | No | `production` | Entorno (`local`, `test`, `production`) |
| `JWT_PUBLIC_KEY` | Si (prod) | - | Clave publica RSA en PEM para validar JWTs |
| `ALLOWED_ORIGINS` | No | `*` | Origenes CORS permitidos (comma-separated) |
| `AZURE_OPENAI_API_KEY` | No | - | API key de Azure OpenAI (sin key se usa stub) |
| `AZURE_OPENAI_ENDPOINT` | No | - | Endpoint de Azure OpenAI |
| `AZURE_OPENAI_MODEL` | No | `gpt-4o-mini` | Deployment name de Azure OpenAI |
| `DB_MAX_OPEN_CONNS` | No | `25` | Max conexiones abiertas al DB |
| `DB_MAX_IDLE_CONNS` | No | `10` | Max conexiones idle |
| `DB_CONN_MAX_LIFETIME_MIN` | No | `30` | Lifetime maximo de conexion (minutos) |
| `DB_CONN_MAX_IDLE_TIME_MIN` | No | `5` | Idle time maximo de conexion (minutos) |

## API Endpoints

Todos bajo `/api/v1`, requieren JWT Bearer token.

| Grupo | Endpoints |
|---|---|
| Auth | `GET /auth/me` |
| Classrooms | `GET/POST /classrooms`, `GET/PUT/DELETE /classrooms/:id`, `GET /classrooms/:id/students` |
| Teachers | `GET /teachers` |
| Students | `GET/POST /students`, `GET/PUT/DELETE /students/:id`, `GET/PUT /students/:id/profile` |
| Catalog | `GET /ramps`, `GET /ramps/:id`, `GET /devices`, `GET /devices/:id` |
| Adaptations | `GET/POST /adaptations`, `GET/PUT/DELETE /adaptations/:id`, `GET /adaptations/:id/resources` |
| Chat | `GET /chat/history/:contextId` |
| AI | `POST /inclusion/recommend`, `POST /inclusion/assist` |
| Dashboard | `GET /dashboard/metrics` |
| Health | `GET /health` (sin auth) |

## Desarrollo

```bash
make lint         # golangci-lint
make vet          # go vet
make test         # go test -race ./...
make test-cover   # test + coverage report
make build        # build binario
make docker       # docker compose up -d
```

## Deployment

El proyecto se despliega en Railway via auto-deploy desde `master`. Usa el `Dockerfile` multi-stage (builder + alpine runtime). Railway verifica salud via `GET /health`.
