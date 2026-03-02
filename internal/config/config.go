package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type BasaltConfig struct {
	Project   string `json:"project"`
	Database  string `json:"database"`
	Framework string `json:"framework"`
	Auth      string `json:"auth"`
	Language  string `json:"language"`
}

func WriteConfig(cfg BasaltConfig) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile("basalt.config.json", data, 0644)
}

func WriteStarterDSL(cfg BasaltConfig) error {
	dsl := fmt.Sprintf(`hclconfig {
  db:        %s
  auth:      %s
  framework: %s
  lang:      %s
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
`, stringsToLower(cfg.Database), stringsToLower(cfg.Auth), stringsToLower(cfg.Framework), stringsToLower(cfg.Language))

	return os.WriteFile("main.bs", []byte(dsl), 0644)
}

func stringsToLower(s string) string {
	return strings.ToLower(strings.ReplaceAll(s, " ", ""))
}
