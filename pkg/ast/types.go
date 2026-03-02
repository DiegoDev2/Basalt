package ast

// File represents the root of the Basalt DSL AST.
type File struct {
	Config    *Config
	Tables    []*Table
	Resources []*Resource
}

// Config represents the HCL-like configuration block.
type Config struct {
	DB        string
	Auth      string
	Framework string
	Lang      string
}

// Table represents a database table definition.
type Table struct {
	Name   string
	Fields []*Field
}

// Field represents a column in a database table.
type Field struct {
	Name       string
	Type       string
	Decorators []*Decorator
}

// Decorator represents a field attribute like @primary or @unique.
type Decorator struct {
	Name string
	Arg  string
}

// Resource represents an API resource mapping a table to endpoints.
type Resource struct {
	Name      string
	Table     string
	Endpoints []*Endpoint
}

// Endpoint represents an HTTP endpoint definition.
type Endpoint struct {
	Method string
	Path   string
}
