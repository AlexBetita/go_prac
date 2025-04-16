### Run

```bash
go run ./cmd/server/main.go
```

### Seed

#### Posts

```bash
SEED_MODE=seed go run internal/db/seed_main/main.go
```

#### Post Embeddings  
⚠️ Requires initially seeded posts

```bash
SEED_MODE=embed go run internal/db/seed_main/main.go
```

### Tests

```bash
go test ./tests -v
```
