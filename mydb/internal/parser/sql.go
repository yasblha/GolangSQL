package parser

import (
	"fmt"
	"strconv"
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

// Condition représente une condition WHERE
type Condition struct {
	Column    string
	Operator  string
	Value     interface{}
	IsNull    bool
	IsNotNull bool
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
	OnConflict string // "REPLACE" ou "IGNORE"
}

// ParseSelect analyse une requête SELECT
func ParseSelect(query string) (*SelectQuery, error) {
	query = strings.TrimSpace(query)
	if !strings.HasPrefix(strings.ToUpper(query), "SELECT") {
		return nil, fmt.Errorf("not a SELECT query")
	}

	// Extraire les colonnes
	parts := strings.SplitN(query, "FROM", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid SELECT syntax")
	}

	columns := strings.TrimSpace(parts[0][6:]) // Enlever "SELECT"
	colList := strings.Split(columns, ",")
	for i := range colList {
		colList[i] = strings.TrimSpace(colList[i])
	}

	// Extraire la table et les conditions
	whereParts := strings.SplitN(parts[1], "WHERE", 2)
	table := strings.TrimSpace(whereParts[0])

	var conditions []Condition
	if len(whereParts) > 1 {
		whereClause := strings.TrimSpace(whereParts[1])
		var err error
		conditions, err = parseWhereConditions(whereClause)
		if err != nil {
			return nil, err
		}
	}

	return &SelectQuery{
		Columns:    colList,
		Table:      table,
		Conditions: conditions,
	}, nil
}

// parseWhereConditions analyse les conditions WHERE
func parseWhereConditions(whereClause string) ([]Condition, error) {
	var conditions []Condition
	
	// Gérer les conditions AND
	andParts := strings.Split(whereClause, "AND")
	for _, part := range andParts {
		part = strings.TrimSpace(part)
		
		// Vérifier pour IS NULL
		if strings.HasSuffix(strings.ToUpper(part), "IS NULL") {
			col := strings.TrimSpace(part[:len(part)-8])
			conditions = append(conditions, Condition{
				Column: col,
				IsNull: true,
			})
			continue
		}

		// Vérifier pour IS NOT NULL
		if strings.HasSuffix(strings.ToUpper(part), "IS NOT NULL") {
			col := strings.TrimSpace(part[:len(part)-12])
			conditions = append(conditions, Condition{
				Column:    col,
				IsNotNull: true,
			})
			continue
		}

		// Gérer les opérateurs de comparaison
		operators := []string{"=", ">", "<", ">=", "<=", "!=", "LIKE"}
		for _, op := range operators {
			if strings.Contains(part, op) {
				parts := strings.Split(part, op)
				if len(parts) != 2 {
					return nil, fmt.Errorf("invalid condition syntax: %s", part)
				}

				col := strings.TrimSpace(parts[0])
				val := strings.TrimSpace(parts[1])

				// Gérer les chaînes entre guillemets
				if strings.HasPrefix(val, "'") && strings.HasSuffix(val, "'") {
					val = val[1 : len(val)-1]
				}

				conditions = append(conditions, Condition{
					Column:   col,
					Operator: op,
					Value:    val,
				})
				break
			}
		}
	}

	return conditions, nil
}

// ParseInsert analyse une requête INSERT
func ParseInsert(query string) (*InsertQuery, error) {
	query = strings.TrimSpace(query)
	if !strings.HasPrefix(strings.ToUpper(query), "INSERT") {
		return nil, fmt.Errorf("not an INSERT query")
	}

	// Gérer INSERT OR REPLACE/IGNORE
	onConflict := ""
	if strings.Contains(strings.ToUpper(query), "OR REPLACE") {
		onConflict = "REPLACE"
		query = strings.Replace(query, "OR REPLACE", "", 1)
	} else if strings.Contains(strings.ToUpper(query), "OR IGNORE") {
		onConflict = "IGNORE"
		query = strings.Replace(query, "OR IGNORE", "", 1)
	}

	// Extraire la table
	parts := strings.SplitN(query[6:], "INTO", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid INSERT syntax")
	}

	// Extraire les colonnes et les valeurs
	tablePart := strings.TrimSpace(parts[1])
	valuesIndex := strings.Index(strings.ToUpper(tablePart), "VALUES")
	if valuesIndex == -1 {
		return nil, fmt.Errorf("invalid INSERT syntax: missing VALUES")
	}

	tableDef := strings.TrimSpace(tablePart[:valuesIndex])
	valuesPart := strings.TrimSpace(tablePart[valuesIndex:])

	// Parser les colonnes
	var columns []string
	if strings.Contains(tableDef, "(") {
		colPart := tableDef[strings.Index(tableDef, "(")+1:strings.Index(tableDef, ")")]
		columns = strings.Split(colPart, ",")
		for i := range columns {
			columns[i] = strings.TrimSpace(columns[i])
		}
		tableDef = strings.TrimSpace(tableDef[:strings.Index(tableDef, "(")])
	}

	// Parser les valeurs
	valuesStr := valuesPart[6:strings.Index(valuesPart, ")")+1] // Enlever "VALUES" et garder les parenthèses
	values, err := parseValues(valuesStr)
	if err != nil {
		return nil, err
	}

	// Parser les conditions WHERE si présentes
	var conditions []Condition
	if whereIndex := strings.Index(strings.ToUpper(valuesPart), "WHERE"); whereIndex != -1 {
		whereClause := strings.TrimSpace(valuesPart[whereIndex+5:])
		conditions, err = parseWhereConditions(whereClause)
		if err != nil {
			return nil, err
		}
	}

	return &InsertQuery{
		Table:      tableDef,
		Columns:    columns,
		Values:     values,
		Conditions: conditions,
		OnConflict: onConflict,
	}, nil
}

// parseValues analyse les valeurs d'une requête INSERT
func parseValues(valuesStr string) ([]interface{}, error) {
	var values []interface{}
	
	// Enlever les parenthèses externes
	valuesStr = strings.TrimSpace(valuesStr)
	if !strings.HasPrefix(valuesStr, "(") || !strings.HasSuffix(valuesStr, ")") {
		return nil, fmt.Errorf("invalid VALUES syntax")
	}
	valuesStr = valuesStr[1 : len(valuesStr)-1]

	// Parser les valeurs individuelles
	parts := strings.Split(valuesStr, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		
		// Gérer les chaînes entre guillemets
		if strings.HasPrefix(part, "'") && strings.HasSuffix(part, "'") {
			values = append(values, part[1:len(part)-1])
			continue
		}

		// Gérer les nombres
		if num, err := strconv.ParseFloat(part, 64); err == nil {
			values = append(values, num)
			continue
		}

		// Gérer NULL
		if strings.ToUpper(part) == "NULL" {
			values = append(values, nil)
			continue
		}

		values = append(values, part)
	}

	return values, nil
}