# 🔧 Migrations (Go Runner)

All migrations are located in the `migrations/` folder. Run them using your custom Go migration command.

---

### Run all migrations (up):

```bash
go run cmd/migrate/main.go -direction up
```

### Run a specific number of steps (up):

```bash
go run cmd/migrate/main.go -direction up -steps 1
```

### Rollback migrations (down):

```bash
go run cmd/migrate/main.go -direction down
```

### Rollback a specific number of steps (down):

```bash
go run cmd/migrate/main.go -direction down -steps 1
```

> ⚠️ Make sure your `DATABASE_URL` is set correctly in `.env`.
