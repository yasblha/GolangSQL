package parser

import (
	"testing"
)

func TestParseSelect(t *testing.T) {
	tests := []struct {
		name    string
		query   string
		want    *SelectQuery
		wantErr bool
	}{
		{
			name:  "Simple SELECT",
			query: "SELECT * FROM users",
			want: &SelectQuery{
				Columns:    []string{"*"},
				Table:      "users",
				Conditions: []Condition{},
			},
			wantErr: false,
		},
		{
			name:  "SELECT with WHERE",
			query: "SELECT name, age FROM users WHERE age > 25",
			want: &SelectQuery{
				Columns: []string{"name", "age"},
				Table:   "users",
				Conditions: []Condition{
					{
						Column:   "age",
						Operator: ">",
						Value:    "25",
					},
				},
			},
			wantErr: false,
		},
		{
			name:  "SELECT with multiple conditions",
			query: "SELECT * FROM users WHERE age > 25 AND name LIKE 'J%'",
			want: &SelectQuery{
				Columns: []string{"*"},
				Table:   "users",
				Conditions: []Condition{
					{
						Column:   "age",
						Operator: ">",
						Value:    "25",
					},
					{
						Column:   "name",
						Operator: "LIKE",
						Value:    "J%",
					},
				},
			},
			wantErr: false,
		},
		{
			name:  "SELECT with IS NULL",
			query: "SELECT * FROM users WHERE email IS NULL",
			want: &SelectQuery{
				Columns: []string{"*"},
				Table:   "users",
				Conditions: []Condition{
					{
						Column: "email",
						IsNull: true,
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseSelect(tt.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseSelect() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !compareSelectQueries(got, tt.want) {
				t.Errorf("ParseSelect() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseInsert(t *testing.T) {
	tests := []struct {
		name    string
		query   string
		want    *InsertQuery
		wantErr bool
	}{
		{
			name:  "Simple INSERT",
			query: "INSERT INTO users (name, age) VALUES ('John', 25)",
			want: &InsertQuery{
				Table:   "users",
				Columns: []string{"name", "age"},
				Values:  []interface{}{"John", 25.0},
			},
			wantErr: false,
		},
		{
			name:  "INSERT OR REPLACE",
			query: "INSERT OR REPLACE INTO users (id, name) VALUES (1, 'John')",
			want: &InsertQuery{
				Table:      "users",
				Columns:    []string{"id", "name"},
				Values:     []interface{}{1.0, "John"},
				OnConflict: "REPLACE",
			},
			wantErr: false,
		},
		{
			name:  "INSERT with NULL",
			query: "INSERT INTO users (name, email) VALUES ('John', NULL)",
			want: &InsertQuery{
				Table:   "users",
				Columns: []string{"name", "email"},
				Values:  []interface{}{"John", nil},
			},
			wantErr: false,
		},
		{
			name:  "INSERT with WHERE",
			query: "INSERT INTO users (name, age) VALUES ('John', 25) WHERE NOT EXISTS (SELECT 1 FROM users WHERE name = 'John')",
			want: &InsertQuery{
				Table:   "users",
				Columns: []string{"name", "age"},
				Values:  []interface{}{"John", 25.0},
				Conditions: []Condition{
					{
						Column:   "name",
						Operator: "=",
						Value:    "John",
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseInsert(tt.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseInsert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !compareInsertQueries(got, tt.want) {
				t.Errorf("ParseInsert() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Fonctions utilitaires pour comparer les requêtes
func compareSelectQueries(a, b *SelectQuery) bool {
	if a == nil || b == nil {
		return a == b
	}
	if a.Table != b.Table {
		return false
	}
	if len(a.Columns) != len(b.Columns) {
		return false
	}
	for i := range a.Columns {
		if a.Columns[i] != b.Columns[i] {
			return false
		}
	}
	if len(a.Conditions) != len(b.Conditions) {
		return false
	}
	for i := range a.Conditions {
		if !compareConditions(a.Conditions[i], b.Conditions[i]) {
			return false
		}
	}
	return true
}

func compareInsertQueries(a, b *InsertQuery) bool {
	if a == nil || b == nil {
		return a == b
	}
	if a.Table != b.Table {
		return false
	}
	if len(a.Columns) != len(b.Columns) {
		return false
	}
	for i := range a.Columns {
		if a.Columns[i] != b.Columns[i] {
			return false
		}
	}
	if len(a.Values) != len(b.Values) {
		return false
	}
	for i := range a.Values {
		if a.Values[i] != b.Values[i] {
			return false
		}
	}
	if a.OnConflict != b.OnConflict {
		return false
	}
	return true
}

func compareConditions(a, b Condition) bool {
	if a.Column != b.Column {
		return false
	}
	if a.Operator != b.Operator {
		return false
	}
	if a.Value != b.Value {
		return false
	}
	if a.IsNull != b.IsNull {
		return false
	}
	if a.IsNotNull != b.IsNotNull {
		return false
	}
	return true
}
