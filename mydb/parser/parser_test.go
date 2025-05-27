package parser

import (
	"testing"
)

func TestParseSelect(t *testing.T) {
	tests := []struct {
		name    string
		query   string
		wantErr bool
	}{
		{
			name:    "Simple SELECT *",
			query:   "SELECT * FROM products",
			wantErr: false,
		},
		{
			name:    "SELECT with WHERE",
			query:   "SELECT name FROM apples WHERE color = 'red'",
			wantErr: false,
		},
		{
			name:    "SELECT with multiple conditions",
			query:   "SELECT * FROM users WHERE age > 18 AND name = 'John'",
			wantErr: false,
		},
		{
			name:    "Invalid SELECT syntax",
			query:   "SELECT FROM products",
			wantErr: true,
		},
		{
			name:    "SELECT with NULL value",
			query:   "SELECT * FROM users WHERE name IS NULL",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseSelect(tt.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseSelect() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestParseInsert(t *testing.T) {
	tests := []struct {
		name    string
		query   string
		wantErr bool
	}{
		{
			name:    "Simple INSERT",
			query:   "INSERT INTO oranges (name, description) VALUES ('Orange', 'Agrume')",
			wantErr: false,
		},
		{
			name:    "INSERT with multiple values",
			query:   "INSERT INTO users (id, name, age) VALUES (1, 'John', 25)",
			wantErr: false,
		},
		{
			name:    "INSERT with NULL value",
			query:   "INSERT INTO users (name, description) VALUES ('John', NULL)",
			wantErr: false,
		},
		{
			name:    "Invalid INSERT syntax",
			query:   "INSERT INTO users VALUES (1, 'John')",
			wantErr: true,
		},
		{
			name:    "INSERT with special characters",
			query:   "INSERT INTO products (name) VALUES ('Product''s name')",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseInsert(tt.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseInsert() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
