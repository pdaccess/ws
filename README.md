# PDAccess

**Privileged Access Management System**

PDAccess is a centralized PAM (Privileged Access Management) platform that provides secure management of services, groups, users, policies, and audit activities through a well-defined REST API.

---

## Overview

PDAccess exposes a versioned REST API (`/v1`) documented via OpenAPI 3.0.3. The full specification lives at:

```
resources/corews-api.yaml
```

All endpoints require **Bearer token authentication** (JWT). The base URL for the production API is:

```
https://app.pdaccess.com/api/v1/ws
```

For local development:

```
http://localhost:80/v1
```

---

## API Reference

The API is organized into eight functional areas:

### Services
Manage privileged services (databases, applications, infrastructure, etc.).

| Method | Path | Description |
|--------|------|-------------|
| `POST` | `/service` | Create a new service |
| `GET` | `/service/{serviceId}` | Get service by ID |
| `PUT` | `/service/{serviceId}` | Update service |
| `DELETE` | `/service/{serviceId}` | Delete service |

### Groups
Organize services and users into logical groups with hierarchical support.

| Method | Path | Description |
|--------|------|-------------|
| `POST` | `/group` | Create a new group |
| `GET` | `/group/{groupId}` | Get group by ID |
| `DELETE` | `/group/{groupId}` | Delete group |
| `GET` | `/group/{groupId}/members` | List group members |
| `POST` | `/group/{groupId}/members` | Add user to group |
| `DELETE` | `/group/{groupId}/members/{userId}` | Remove user from group |
| `GET` | `/group/{groupId}/policy` | Get policies assigned to group |
| `POST` | `/group/{groupId}/policy` | Assign policy to group |
| `DELETE` | `/group/{groupId}/policy/{policyId}` | Remove policy from group |

### Policies
Define and apply access, security, and network policies to groups.

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/policies` | List all policies (paginated) |
| `POST` | `/policies` | Create a policy |
| `GET` | `/policies/{policyId}` | Get policy by ID |
| `PUT` | `/policies/{policyId}` | Update policy |
| `DELETE` | `/policies/{policyId}` | Delete policy |

Policy types: `Access`, `Security`, `Network`

### Admin
User management, audit logs, system health, and configuration.

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/admin/users` | List users (filterable by role/status) |
| `POST` | `/admin/users` | Create a user |
| `GET` | `/admin/users/{userId}` | Get user by ID |
| `PUT` | `/admin/users/{userId}` | Update user |
| `DELETE` | `/admin/users/{userId}` | Delete user |
| `PUT` | `/admin/users/{userId}/status` | Change user status |
| `GET` | `/admin/audit-logs` | Retrieve audit logs |
| `GET` | `/admin/system-health` | System health metrics |
| `GET` | `/admin/settings` | Get system settings |
| `PUT` | `/admin/settings` | Update system settings |

User statuses: `active`, `inactive`, `suspended`

### Activities
Immutable activity log for compliance and audit purposes.

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/activities` | List activities (filterable by severity, service, group, session) |
| `GET` | `/activities/{activityId}` | Get activity by ID |

### Alarms
Real-time alarm and alert management.

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/alarms` | List alarms (filterable by severity) |
| `GET` | `/alarms/{alarmId}` | Get alarm by ID |
| `POST` | `/alarms/{alarmId}/acknowledge` | Acknowledge an alarm |

### Paste
Secure, ephemeral snippet sharing with burn-after-read support.

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/paste` | List pastes |
| `POST` | `/paste` | Create a paste (optional expiration, burn-after-read) |
| `GET` | `/paste/{pasteId}` | Read paste (consumed if burn-after-read is enabled) |
| `DELETE` | `/paste/{pasteId}` | Delete paste |

### Search
Cross-resource full-text and semantic search.

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/search` | Search across services, groups, users, and activities |

Search supports pagination, sorting, and filtering by resource type (`Service`, `Group`, `User`, `Activity`).

---

## Authentication

All requests must include a valid JWT in the `Authorization` header:

```
Authorization: Bearer <token>
```

Requests without a valid token will receive a `401 Unauthorized` response.

---

## OpenAPI Specification

The canonical API specification is maintained at `resources/corews-api.yaml`. The server implementation is code-generated from this file — **always edit the spec first**, then regenerate the server code:

```sh
make generate
```

The generated server code lives at `internal/platform/handlers/external/server.gen.go` and must not be edited by hand.

---

## Running Locally

```sh
make local
```

For containerized development:

```sh
make build
```

To run the full integration test suite:

```sh
make cicd-tests
```

---