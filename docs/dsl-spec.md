# Basalt DSL Specification

Basalt uses a custom HCL-like DSL to define backend schemas and resources.

## Config Block

The `hclconfig` block defines the project-wide settings.

- `db`: The database provider (`supabase`, `postgresql`, `mysql`).
- `auth`: Authentication method (`jwt`, `apikey`, `none`).
- `framework`: Web framework (`express`, `fastify`, `hono`).
- `lang`: Output language (`typescript`, `javascript`).

## Table Block

Defines a database table.

- `id: uuid @primary`
- `email: string @unique`
- `created_at: timestamp @default(now)`

## Resource Block

Maps a table to API endpoints.

```hcl
resource "users" {
  table: users
  endpoints {
    GET    /users
    POST   /users
  }
}
```
