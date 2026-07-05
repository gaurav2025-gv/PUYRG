# PUYRG Backend

Initial Go API scaffold for PUYRG.

Run:

```bash
go run ./cmd/api
```

Endpoints:

- `GET /api/health`
- `GET /api/dashboard`
- `GET /api/revisions/today`
- `POST /api/ledger/attempts`

This version uses the Go standard library so the API is immediately buildable. Gin, GORM, PostgreSQL, Redis, and JWT will be added when the database schema is implemented.

