package parser

import (
	"errors"
	"strings"
)

// Type d'opérateur pour les conditions
type Operator string

const (
	EQ  Operator = "="
	NEQ Operator = "!="
	GT  Operator = ">"
	LT  Operator = "<"
	GTE Operator = ">="
	LTE Operator = "<="
)

// Condition représente une condition dans une clause WHERE
type Condition struct {
	Column    string
	Operator  Operator
	Value     interface{}
}

// SelectQuery représente une requête SELECT
type SelectQuery struct {
	Table      string
	Columns    []string
	Conditions []Condition
}

// InsertQuery représente une requête INSERT
type InsertQuery struct {
	Table   string
	Columns []string
	Values  []interface{}
}

// ParseSelect parse une requête SELECT
func ParseSelect(query string) (*SelectQuery, error) {
	query = strings.TrimSpace(query)
	if !strings.HasPrefix(strings.ToUpper(query), "SELECT") {
		return nil, errors.New("not a SELECT query")
	}

	// Enlever le "SELECT" initial
	query = strings.TrimSpace(query[6:])
	
	// Parser les colonnes
	parts := strings.Split(query, "FROM")
	if len(parts) != 2 {
		return nil, errors.New("invalid SELECT syntax: missing FROM clause")
	}

	columnsStr := strings.TrimSpace(parts[0])
	columns := strings.Split(columnsStr, ",")
	for i := range columns {
		columns[i] = strings.TrimSpace(columns[i])
	}

	// Parser la table et les conditions
	fromPart := strings.TrimSpace(parts[1])
	whereParts := strings.Split(fromPart, "WHERE")
	
	table := strings.TrimSpace(whereParts[0])
	conditions := []Condition{}

	if len(whereParts) > 1 {
		whereClause := strings.TrimSpace(whereParts[1])
		// TODO: Parser les conditions WHERE
		// Pour l'instant, on ne gère que les conditions simples
		if strings.Contains(whereClause, "=") {
			parts := strings.Split(whereClause, "=")
			if len(parts) == 2 {
				conditions = append(conditions, Condition{
					Column:   strings.TrimSpace(parts[0]),
					Operator: EQ,
					Value:    strings.TrimSpace(parts[1]),
				})
			}
		}
	}

	return &SelectQuery{
		Table:      table,
		Columns:    columns,
		Conditions: conditions,
	}, nil
}

// ParseInsert parse une requête INSERT
func ParseInsert(query string) (*InsertQuery, error) {
	query = strings.TrimSpace(query)
	if !strings.HasPrefix(strings.ToUpper(query), "INSERT INTO") {
		return nil, errors.New("not an INSERT query")
	}

	// Enlever le "INSERT INTO" initial
	query = strings.TrimSpace(query[11:])
	
	// Trouver la position de la première parenthèse ouvrante
	openParen := strings.Index(query, "(")
	if openParen == -1 {
		return nil, errors.New("invalid INSERT syntax: missing opening parenthesis")
	}

	// Extraire le nom de la table
	table := strings.TrimSpace(query[:openParen])
	
	// Trouver la position de la parenthèse fermante correspondante
	closeParen := strings.LastIndex(query, ")")
	if closeParen == -1 {
		return nil, errors.New("invalid INSERT syntax: missing closing parenthesis")
	}

	// Extraire la partie entre parenthèses
	columnsPart := query[openParen+1:closeParen]
	
	// Séparer les colonnes et les valeurs
	parts := strings.Split(columnsPart, ") VALUES (")
	if len(parts) != 2 {
		return nil, errors.New("invalid INSERT syntax: missing VALUES clause")
	}

	// Parser les colonnes
	columns := strings.Split(parts[0], ",")
	for i := range columns {
		columns[i] = strings.TrimSpace(columns[i])
	}

	// Parser les valeurs
	valuesParts := strings.Split(parts[1], ",")
	values := make([]interface{}, len(valuesParts))
	
	for i, v := range valuesParts {
		v = strings.TrimSpace(v)
		// Nettoyer les valeurs (enlever les guillemets si présents)
		if strings.HasPrefix(v, "'") && strings.HasSuffix(v, "'") {
			v = v[1 : len(v)-1]
		}
		values[i] = v
	}

	return &InsertQuery{
		Table:   table,
		Columns: columns,
		Values:  values,
	}, nil
}