# Project Structure

```
pdaccess/ws
├── cmd/                 # Application entrypoints
│   ├── main.go         # CLI entrypoint
│   └── app/            # Application logic
│       ├── server.go   # HTTP server setup
│       ├── config.go   # Configuration
│       └── healthcheck.go
├── internal/           # Private application code
│   ├── core/           # Domain logic (ports, services, domains)
│   │   ├── domain/     # Domain models (assets, errors)
│   │   ├── ports/      # Interface definitions (service.go, vector.go)
│   │   └── service/    # Business logic (service.go, builder.go)
│   ├── database/       # Database adapters (postgres)
│   │   ├── database.go # DB init, migrations
│   │   ├── assets.go   # Unified AssetRepository
│   │   └── audit.go    # AuditRepository
│   ├── adapters/       # External integrations
│   │   └── embed.go    # VectorGenerator (ONNX embeddings)
│   └── platform/       # Framework code
│       ├── handlers/   # HTTP handlers
│       │   ├── external/  # Generated OpenAPI server
│       │   ├── custom/    # Custom handlers (health)
│       │   └── http.go    # Handler implementation
│       └── servers/    # Server setup (http.go)
├── cicd/               # Integration tests
│   └── tests/          # Ginkgo test suite
├── pkg/                # Public packages
│   └── http/           # Generated HTTP client
├── resources/          # API specs and configs
│   ├── corews-api.yaml # OpenAPI 3.0.3 spec (corews-api2)
│   ├── api-config-client.yaml
│   └── api-config-server.yaml
├── lib/                # External libraries (native deps)
├── Makefile            # Build targets
├── Dockerfile          # Container build
├── docker-compose.yml  # Local dev environment
└── go.mod / go.sum    # Go dependencies
```

## Hexagonal Architecture (Ports & Adapters)

```
┌──────────────────────────────────────────────────────────┐
│                     Inbound Adapters                       │
│  internal/platform/handlers/external/  — Generated HTTP   │
│  internal/platform/handlers/custom/    — Custom handlers  │
│  internal/platform/handlers/http.go    — Handler impl     │
└────────────────────────┬─────────────────────────────────┘
                         │  calls
┌────────────────────────▼─────────────────────────────────┐
│                      Ports                                │
│  internal/core/ports/service.go  — Service interface     │
│                     AssetOperations, IdentityOperations,  │
│                     VaultOperations, AdminConfigOperations│
│                     AuditOperations                       │
│  internal/core/ports/vector.go   — VectorGenerator        │
└────────────────────────┬─────────────────────────────────┘
                         │  implements
┌────────────────────────▼─────────────────────────────────┐
│                   Service/Business Logic                   │
│  internal/core/service/  — Impl struct                    │
│                     service.go   (Asset CRUD, embeddings) │
│                     builder.go   (Wiring / DI)            │
└────────────────────────┬─────────────────────────────────┘
                         │  calls
┌────────────────────────▼─────────────────────────────────┐
│                   Domain Models                            │
│  internal/core/domain/  — Pure Go structs & errors        │
│                     assets.go   (Asset, User, UserGroup,  │
│                                  VaultMembership,         │
│                                  AuditEntry, AdminConfig, │
│                                  Vector)                  │
│                     errors.go   (Domain error types)      │
└────────────────────────┬─────────────────────────────────┘
                         │  calls
┌────────────────────────▼─────────────────────────────────┐
│                    Outbound Adapters                       │
│  internal/database/  — PostgreSQL implementations         │
│                     assets.go   (AssetRepository)         │
│                     audit.go    (AuditRepository)         │
│  internal/adapters/  — External integrations              │
│                     embed.go    (VectorGenerator)         │
└──────────────────────────────────────────────────────────┘
```

## Domain Models
- `assets.go` — Asset (unified: service, service_group, vault, secret, policy, paste, jit), User, UserGroup, GroupMember, VaultMembership, AuditEntry, AdminConfig, Vector
- `errors.go` — ValidationError, NotFoundError, InvalidIDError, InternalError

## API Endpoints (corews-api.yaml)

| Path | Method | Handler | Purpose |
|------|--------|---------|---------|
| `/admin/config` | GET | `GetAdminConfig` | List admin config |
| `/admin/config` | PATCH | `PatchAdminConfig` | Update admin config |
| `/assets` | GET | `GetAssets` | Hybrid search (query via `?q=`) |
| `/assets` | POST | `PostAssets` | Create any asset (type in body) |
| `/assets/{id}/actions` | POST | `PostAssetsIdActions` | Polymorphic actions |
| `/audit/logs` | GET | `GetAuditLogs` | Query audit trail |
| `/identity/users` | GET | `GetIdentityUsers` | List users |
| `/identity/users` | POST | `PostIdentityUsers` | Onboard user |
| `/identity/groups` | GET | `GetIdentityGroups` | List groups (paginated) |
| `/identity/groups` | POST | `PostIdentityGroups` | Create group |
| `/identity/groups/{groupId}/memberships` | GET | `GetIdentityGroupsGroupIdMemberships` | List group members (paginated) |
| `/identity/groups/{groupId}/memberships` | POST | `PostIdentityGroupsGroupIdMemberships` | Assign user to group |
| `/identity/groups/{groupId}/memberships` | DELETE | `DeleteIdentityGroupsGroupIdMemberships` | Remove user from group |
| `/vaults/{vaultId}/memberships` | GET | `GetVaultsVaultIdMemberships` | List vault members |
| `/vaults/{vaultId}/memberships` | POST | `PostVaultsVaultIdMemberships` | Add vault member |

## Database Tables

| Table | Columns | Description |
|-------|---------|-------------|
| `ws_assets` | id, type, name, owner_id, parent_id (UUID), spec (JSONB), embedding (FLOAT8[]), created_at | Unified asset storage |
| `ws_users` | id, username, email, status, created_at | Identity users |
| `ws_groups` | id, name, description, created_at | Identity groups |
| `ws_group_members` | group_id, user_id, created_at | Group memberships |
| `ws_vault_memberships` | vault_id, member_id, member_type, role, created_at | Vault access |
| `ws_admin_config` | key, value, updated_at | Key-value settings |
| `ws_audit_logs` | id, timestamp, actor_id, action, resource_id, success | Audit trail |

## Embedding / Vector Search

Assets have a `embedding DOUBLE PRECISION[]` column storing vector embeddings generated by the ONNX-based VectorGenerator. When an asset is created, the service layer generates an embedding from the asset name and stores it. The VectorGenerator is optional — if ONNX/CGO is unavailable, a stub is used and embedding generation is silently skipped.

## Service Layer

The service `Impl` struct takes:
- `AssetRepository` — all CRUD for assets, identity, vault memberships, admin config
- `AuditRepository` — audit log search/create
- `VectorGenerator` (optional) — generates embeddings for asset names

## Makefile Targets

| Target | Description |
|--------|-------------|
| `make format` | Format code (tidy, vet, gofmt, fix) |
| `make unit-tests` | Run unit tests on internal packages |
| `make cicd-tests` | Run integration tests (ginkgo) |
| `make build` | Build Docker image |
| `make generate` | Regenerate OpenAPI server/client code |
| `make local` | Build and run locally with CGO |
| `make clean` | Remove binary artifacts |
| `make install-onnx` | Download ONNX Runtime libraries |

## Development Workflow
1. Edit `resources/corews-api.yaml` to modify API
2. Run `make generate` to regenerate server/client code
3. Implement/customize handlers in `internal/platform/handlers/http.go`
4. Run `make local` or `make cicd-tests` to test
