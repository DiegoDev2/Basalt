# Basalt DSL Specification v1.0

**Status:** Draft  
**Version:** 1.0.0  
**Created:** 2026-03-02  

---

## Abstract

Basalt is an open standard for describing backend services using a declarative DSL (Domain-Specific Language). Any AI model or tool that understands this spec can generate valid `.bs` files. Any Basalt-compatible generator can then compile those files into production-ready backend code.

The goal is to separate **intent** (what the backend should do) from **implementation** (how it is generated), enabling any AI to describe backends in a language-agnostic, framework-agnostic way.

---

## 1. File Format

Basalt files use the `.bs` extension. A `.bs` file is a UTF-8 encoded text file composed of four top-level block types:

- `config` — global project settings (required, exactly once)
- `table` — database table definitions (one or more)
- `resource` — API resource definitions (one or more)
- `role` — permission roles (optional)

Blocks are defined using curly braces `{}`. Keys and values are separated by `:`. Strings use double quotes `"`. Comments use `#`.

---

## 2. Config Block

The `config` block defines global settings for the project. It must appear once and before any other block.

```hcl
config {
  db:        supabase         # Database provider
  auth:      jwt              # Authentication method
  framework: express          # Target framework
  lang:      typescript       # Output language
  pagination: true            # Enable automatic pagination (optional, default: false)
}
```

### 2.1 Config Fields

| Field | Required | Values | Default |
|-------|----------|--------|---------|
| `db` | yes | `supabase`, `postgres`, `mysql` | — |
| `auth` | yes | `jwt`, `apikey`, `none` | — |
| `framework` | no | `express`, `fastify`, `hono` | `express` |
| `lang` | no | `typescript`, `javascript` | `typescript` |
| `pagination` | no | `true`, `false` | `false` |

---

## 3. Table Block

A `table` block defines a database table and its fields.

```hcl
table "users" {
  id:         uuid      @primary
  email:      string    @unique @min(5) @max(255)
  name:       string    @min(2) @max(100)
  age:        number    @min(0) @max(120)
  role:       string    @default("user")
  created_at: timestamp @default(now)
}
```

### 3.1 Field Types

| DSL Type | TypeScript | SQL |
|----------|-----------|-----|
| `uuid` | `string` | `uuid` |
| `string` | `string` | `text` |
| `number` | `number` | `integer` |
| `float` | `number` | `float8` |
| `boolean` | `boolean` | `boolean` |
| `timestamp` | `string` | `timestamptz` |
| `json` | `Record<string, any>` | `jsonb` |

### 3.2 Decorators

Decorators modify the behavior of a field. Multiple decorators can be applied to a single field.

#### Structural Decorators

| Decorator | Description | SQL Effect |
|-----------|-------------|------------|
| `@primary` | Primary key | `PRIMARY KEY` |
| `@unique` | Unique constraint | `UNIQUE` |
| `@default(value)` | Default value | `DEFAULT value` |
| `@default(now)` | Current timestamp | `DEFAULT now()` |
| `@ref(table.field)` | Foreign key relation | `REFERENCES table(field)` |
| `@nullable` | Allow null values | column is nullable |

#### Validation Decorators

| Decorator | Applies to | Description |
|-----------|-----------|-------------|
| `@min(n)` | `string`, `number`, `float` | Min length (string) or min value (number) |
| `@max(n)` | `string`, `number`, `float` | Max length (string) or max value (number) |
| `@regex("pattern")` | `string` | Must match regex pattern |
| `@email` | `string` | Must be valid email format |
| `@url` | `string` | Must be valid URL format |

### 3.3 Relations

Foreign keys are defined using `@ref(table.field)`:

```hcl
table "posts" {
  id:      uuid   @primary
  title:   string @min(1) @max(200)
  user_id: uuid   @ref(users.id)
}
```

---

## 4. Role Block

Roles define named permission groups that can be applied to endpoints.

```hcl
role "admin" {
  description: "Full access to all resources"
}

role "user" {
  description: "Standard authenticated user"
}

role "guest" {
  description: "Unauthenticated access"
}
```

---

## 5. Resource Block

A `resource` block defines an API resource, its endpoints, and access control.

```hcl
resource "users" {
  table: users
  paginate: true       # Override global pagination setting (optional)

  endpoints {
    GET    /users               roles: [guest, user, admin]
    POST   /users               roles: [admin]
    GET    /users/:id           roles: [user, admin]
    PUT    /users/:id           roles: [user, admin]
    DELETE /users/:id           roles: [admin]
  }
}
```

### 5.1 Endpoint Syntax

```
METHOD  /path   roles: [role1, role2]
```

- **METHOD** — `GET`, `POST`, `PUT`, `PATCH`, `DELETE`
- **/path** — must start with `/` and match the resource name as prefix
- **roles** — optional list of roles allowed to access this endpoint. If omitted, uses the global `auth` setting.

### 5.2 Pagination

When pagination is enabled (globally or per resource), `GET` list endpoints automatically support:

- `?page=1` — page number (default: 1)
- `?limit=20` — items per page (default: 20, max: 100)

Response format:
```json
{
  "data": [...],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 100,
    "pages": 5
  }
}
```

### 5.3 Multiple Databases

When using multiple databases, each resource can specify its own `db`:

```hcl
resource "users" {
  table: users
  db: supabase
  endpoints { ... }
}

resource "logs" {
  table: logs
  db: postgres
  endpoints { ... }
}
```

---

## 6. Complete Example

```hcl
# Blog API — Basalt DSL v1.0

config {
  db:        supabase
  auth:      jwt
  framework: express
  lang:      typescript
  pagination: true
}

role "admin" {
  description: "Full access"
}

role "user" {
  description: "Authenticated user"
}

role "guest" {
  description: "Public access"
}

table "users" {
  id:         uuid      @primary
  email:      string    @unique @email
  name:       string    @min(2) @max(100)
  password:   string    @min(8)
  role:       string    @default("user")
  created_at: timestamp @default(now)
}

table "posts" {
  id:         uuid      @primary
  title:      string    @min(1) @max(200)
  body:       string
  user_id:    uuid      @ref(users.id)
  published:  boolean   @default(false)
  created_at: timestamp @default(now)
}

table "comments" {
  id:         uuid      @primary
  body:       string    @min(1) @max(1000)
  user_id:    uuid      @ref(users.id)
  post_id:    uuid      @ref(posts.id)
  created_at: timestamp @default(now)
}

resource "users" {
  table: users
  endpoints {
    GET    /users           roles: [admin]
    POST   /users           roles: [guest]
    GET    /users/:id       roles: [user, admin]
    PUT    /users/:id       roles: [user, admin]
    DELETE /users/:id       roles: [admin]
  }
}

resource "posts" {
  table: posts
  paginate: true
  endpoints {
    GET    /posts           roles: [guest, user, admin]
    POST   /posts           roles: [user, admin]
    GET    /posts/:id       roles: [guest, user, admin]
    PUT    /posts/:id       roles: [user, admin]
    DELETE /posts/:id       roles: [admin]
  }
}

resource "comments" {
  table: comments
  paginate: true
  endpoints {
    GET    /comments        roles: [guest, user, admin]
    POST   /comments        roles: [user, admin]
    DELETE /comments/:id    roles: [admin]
  }
}
```

---

## 7. Grammar (EBNF)

```ebnf
file        = config { table | resource | role }

config      = "config" "{" config_field* "}"
config_field = ident ":" value

table       = "table" string "{" field* "}"
field       = ident ":" type decorator*
type        = "uuid" | "string" | "number" | "float" | "boolean" | "timestamp" | "json"
decorator   = "@" ident [ "(" arg ")" ]
arg         = string | number | "now"

role        = "role" string "{" role_field* "}"
role_field  = ident ":" string

resource    = "resource" string "{" resource_field* endpoints_block "}"
resource_field = ident ":" value
endpoints_block = "endpoints" "{" endpoint* "}"
endpoint    = method path [ "roles:" "[" ident { "," ident } "]" ]
method      = "GET" | "POST" | "PUT" | "PATCH" | "DELETE"
path        = "/" { ident | ":" ident | "/" }

string      = '"' { char } '"'
ident       = letter { letter | digit | "_" }
value       = string | ident | "true" | "false"
number      = digit { digit } [ "." digit { digit } ]
```

---

## 8. AI Generation Guidelines

This section is specifically for AI models generating `.bs` files.

### Rules

1. Always include exactly one `config` block at the top
2. Every `resource` must reference a `table` that is defined in the same file
3. Every `@ref` must point to a field that exists in the referenced table
4. Every role used in endpoints must be declared in a `role` block
5. Paths must start with `/` followed by the resource name
6. Never generate fields named `password` without `@min(8)`
7. Always include `id`, `created_at` in every table
8. Never add markdown, explanations, or code fences — output only the `.bs` content

### Recommended System Prompt for AI Models

```
You are a Basalt DSL generator. Given a description of a backend API,
you respond ONLY with a valid .bs file following the Basalt DSL Specification v1.0.

Rules:
- Output only the .bs file content, no explanations, no markdown fences
- Always include config, tables, roles, and resources
- Every resource must have a matching table
- Use @ref for all foreign key relationships
- Follow the grammar exactly as specified

Basalt DSL Specification: https://github.com/DiegoDev2/basalt/spec/RFC-0001.md
```

---

## 9. Versioning

This spec follows semantic versioning. Breaking changes increment the major version. The spec version must be compatible with the Basalt generator version being used.

---

## 10. Contributing

The Basalt spec is open. Proposals for new features must be submitted as RFCs at `github.com/DiegoDev2/basalt/spec`.
