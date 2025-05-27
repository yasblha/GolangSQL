package sql

// Type représente un type de données SQL
type Type string

const (
	// Types de données supportés
	TypeInteger Type = "INTEGER"
	TypeText    Type = "TEXT"
	TypeBlob    Type = "BLOB"
	TypeReal    Type = "REAL"
	TypeNull    Type = "NULL"
)

// Column représente une colonne dans une table
type Column struct {
	Name     string
	Type     Type
	Nullable bool
	Default  interface{}
}

// Table représente une table dans la base de données
type Table struct {
	Name    string
	Columns []Column
}

// Query représente une requête SQL générique
type Query interface {
	GetType() string
	GetTable() string
}

// SelectQuery représente une requête SELECT
type SelectQuery struct {
	Columns    []string
	Table      string
	Conditions []Condition
	OrderBy    string
	Limit      int
	Offset     int
}

// InsertQuery représente une requête INSERT
type InsertQuery struct {
	Table      string
	Columns    []string
	Values     []interface{}
	Conditions []Condition
	OnConflict string
}

// Condition représente une condition WHERE
type Condition struct {
	Column    string
	Operator  string
	Value     interface{}
	IsNull    bool
	IsNotNull bool
}

// GetType implémente l'interface Query
func (q *SelectQuery) GetType() string { return "SELECT" }
func (q *SelectQuery) GetTable() string { return q.Table }

// GetType implémente l'interface Query
func (q *InsertQuery) GetType() string { return "INSERT" }
func (q *InsertQuery) GetTable() string { return q.Table } 