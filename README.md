# 🔧 Migrations & Development Setup

All migrations are located in the `migrations/` folder. You can run them either directly with `go run` or using the `Makefile` shortcuts.

---

## 🚀 Development

Start the dev server with hot-reload (using [Air](https://github.com/air-verse/air)):

```bash
make dev
```

---

## 🗄️ Database (Docker Compose)

Spin up Postgres + pgAdmin (if configured in `docker-compose.yml`):

```bash
make db-up
```

Stop the database:

```bash
make db-down
```

---

## 📜 Migrations (Go Runner)

### Run all migrations (up):

```bash
make migrate-up
# or
go run cmd/migrate/main.go -direction up
```

### Run a specific number of steps (up):

```bash
go run cmd/migrate/main.go -direction up -steps 1
```

### Rollback migrations (down):

```bash
make migrate-down
# or
go run cmd/migrate/main.go -direction down
```

### Rollback a specific number of steps (down):

```bash
go run cmd/migrate/main.go -direction down -steps 1
```

---

> ⚠️ Make sure your `DATABASE_URL` is set correctly in `.env` before running migrations or starting the server.
