# DistriShop

## Overview
Distributed shop platform demonstrating a microservices architecture in Go. Provides product catalog, account management, (extensible) order service and a GraphQL API Gateway for unified client access.

## Architecture
- Services: Account, Catalog, Order (scaffold), GraphQL Gateway
- Communication: gRPC between services, GraphQL to clients
- Data isolation: each service owns its persistence layer (repository pattern)
- Containerization: Docker + docker-compose for local orchestration
- Contracts: Protocol Buffers (.proto) for gRPC; schema.graphql for GraphQL

## Tech Stack
- Go (gRPC, protobuf-generated stubs)
- gqlgen for GraphQL schema-first generation
- Docker / docker-compose for environment parity
- Clean layering: transport -> service -> repository

## Repository Structure (excerpt)
```
catalog/        # Catalog service (gRPC server, repo, client)
account/        # Account service
graphql/        # GraphQL gateway (resolvers, schema, generated code)
order/          # (placeholder for order service implementation)
proto (*.proto) # gRPC service contracts per bounded context
```

## Running Locally
1. Generate protobuf & GraphQL (if modified):
   - Protobuf: `protoc --go_out=. --go-grpc_out=. <service>.proto` (handled via tools.go if configured)
   - GraphQL: `go run github.com/99designs/gqlgen generate`
2. Build & run all services:
   - `docker compose up --build`
3. Access:
   - GraphQL Playground: http://localhost:8080 (adjust if different)

## Development
- Hot reload: rebuild specific service container after code change
- Add a new RPC: update `<service>.proto`, regenerate stubs, implement in `server.go` & expose via client
- Add a GraphQL field: update `schema.graphql`, run gqlgen generate, implement resolver

## gRPC Pattern
- Thin client wrappers (e.g., `catalog/client.go`) provide typed methods & translation to internal models
- Pagination supported via `Skip` / `Take` in GetProducts

## GraphQL
- Aggregates data across microservices
- Centralized error mapping and type-safe resolvers

## Future Enhancements
- Implement Order service logic & persistence
- Add observability: OpenTelemetry tracing + Prometheus metrics
- Integrate centralized config & service discovery (Consul / etcd) beyond compose DNS
- CI pipeline (lint, tests, security scan)

## License
Internal / Educational sample (add chosen license).
