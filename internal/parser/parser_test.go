package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParser(t *testing.T) {
	input := `
hclconfig {
  db:        supabase
  auth:      jwt
  framework: express
  lang:      typescript
}

table "users" {
  id:         uuid      @primary
  email:      string    @unique
  name:       string
  created_at: timestamp @default(now)
}

resource "users" {
  table: users
  endpoints {
    GET    /users
    POST   /users
    GET    /users/:id
    PUT    /users/:id
    DELETE /users/:id
  }
}
`
	l := NewLexer(input)
	p := NewParser(l)
	file, err := p.ParseFile()

	assert.Nil(t, err)
	assert.NotNil(t, file.Config)
	assert.Equal(t, "supabase", file.Config.DB)
	assert.Equal(t, "jwt", file.Config.Auth)
	assert.Equal(t, "express", file.Config.Framework)
	assert.Equal(t, "typescript", file.Config.Lang)

	assert.Equal(t, 1, len(file.Tables))
	assert.Equal(t, "users", file.Tables[0].Name)
	assert.Equal(t, 4, len(file.Tables[0].Fields))

	assert.Equal(t, 1, len(file.Resources))
	assert.Equal(t, "users", file.Resources[0].Name)
	assert.Equal(t, "users", file.Resources[0].Table)
	assert.Equal(t, 5, len(file.Resources[0].Endpoints))
}
