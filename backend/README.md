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

Roughly 2.098K input tokens for 33 requests < 0.01$ for ada v2

```bash
SEED_MODE=embed go run internal/db/seed_main/main.go
```

### Tests

```bash
go test ./tests -v
```
