## Project Structure (Ride service)

```
ride/
├── api/                        # API contract + generator
│   ├── openapi.yaml            # OpenAPI spec
│   ├── oapi-types.cfg.yaml     # Generator config for models
│   ├── oapi-server.cfg.yaml    # Generator config for Gin server
│   └── generator.go            # go:generate commands
├── internal/
│   ├── adapters/
│   │   └── http/
│   │       ├── api/            # Generated OpenAPI code (DO NOT EDIT)
│   │       │   ├── types.gen.go
│   │       │   └── server.gen.go
│   │       └── converter/      # Manual converters between API and domain
│   │           └── rest_converter.go
│   ├── models/                 # Pure business entities
│   │   └── ride.go
│   └── ports/                  # Interfaces (ports)
│       └── ride_service.go
```

Generated code (`*.gen.go`) is never edited by hand — regenerate with `go generate ./...` after changing `api/openapi.yaml`.
