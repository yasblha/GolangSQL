package sql

import (
	"fmt"
	"strings"
	"mydb/pkg/types"
)

// Parser représente un parser SQL
type Parser interface {
	Parse(query string) (types.Query, error)
}

// Parse analyse une requête SQL et retourne une Query
func Parse(query string) (types.Query, error) {
	// Déterminer le type de requête
	queryType := getQueryType(query)

	switch queryType {
	case "SELECT":
		return parseSelect(query)
	case "INSERT":
		return parseInsert(query)
	case "CREATE":
		return parseCreate(query)
	default:
		return nil, fmt.Errorf("type de requête non supporté: %s", queryType)
	}
}

// getQueryType détermine le type de requête SQL
func getQueryType(query string) string {
	// Normaliser la requête
	query = normalizeQuery(query)

	// Extraire le premier mot en ignorant les espaces en début
	words := strings.Fields(query)
	if len(words) == 0 {
		return ""
	}

	return words[0]
}

// normalizeQuery normalise une requête SQL
func normalizeQuery(query string) string {
	// Supprimer les espaces en début et fin
	query = strings.TrimSpace(query)

	// Convertir en majuscules
	query = strings.ToUpper(query)

	return query
}

// splitWords divise une requête en mots
func splitWords(query string) []string {
	// Supprimer les points-virgules en fin de requête
	query = strings.TrimRight(query, ";")

	// Diviser la requête en mots
	return strings.Fields(query)
}

// parseSelect analyse une requête SELECT
func parseSelect(query string) (types.Query, error) {
	// Supprimer les espaces en début et fin
	query = strings.TrimSpace(query)
	
	// Extraire la partie SELECT
	selectPart := strings.SplitN(query, "FROM", 2)
	if len(selectPart) != 2 {
		return nil, fmt.Errorf("syntaxe SELECT invalide")
	}
	
	// Extraire les colonnes
	columns := strings.Split(selectPart[0][6:], ",") // Enlever "SELECT "
	for i, col := range columns {
		columns[i] = strings.TrimSpace(col)
	}

	// Extraire la table et la condition WHERE
	tablePart := selectPart[1]
	table := ""
	conditions := make([]types.Condition, 0)
	if strings.Contains(tablePart, "WHERE") {
		// Séparer la table de la condition WHERE
		tableWhere := strings.SplitN(tablePart, "WHERE", 2)
		if len(tableWhere) != 2 {
			return nil, fmt.Errorf("syntaxe WHERE invalide")
		}
		table = strings.TrimSpace(tableWhere[0])
		wherePart := strings.TrimSpace(tableWhere[1])
		
		// Pour le moment, on gère seulement les conditions simples de type "colonne = valeur"
		whereParts := strings.Split(wherePart, "=")
		if len(whereParts) != 2 {
			return nil, fmt.Errorf("syntaxe WHERE invalide")
		}
		
		// Extraire le nom de la colonne et la valeur
		colName := strings.TrimSpace(whereParts[0])
		value := strings.TrimSpace(whereParts[1])
		
		// Supprimer les guillemets si présents
		value = strings.Trim(value, "'")
		
		// Ajouter la condition
		conditions = append(conditions, types.Condition{
			Column:  colName,
			Operator: "=",
			Value:   value,
		})
	} else {
		table = strings.TrimSpace(tablePart)
	}

	// Supprimer les points-virgules
	table = strings.Trim(table, ";")
	table = strings.TrimSpace(table)

	// Afficher les détails pour débogage
	fmt.Printf("DEBUG: SELECT Table: %s\n", table)
	fmt.Printf("DEBUG: SELECT Columns: %v\n", columns)
	fmt.Printf("DEBUG: SELECT Conditions: %v\n", conditions)

	return &types.SelectQuery{
		Columns:    columns,
		Table:      table,
		Conditions: conditions,
	}, nil
}

// parseInsert analyse une requête INSERT
func parseInsert(query string) (types.Query, error) {
	// Supprimer les espaces en début et fin
	query = strings.TrimSpace(query)
	
	// Extraire la partie INSERT
	insertPart := strings.SplitN(query, "VALUES", 2)
	if len(insertPart) != 2 {
		return nil, fmt.Errorf("syntaxe INSERT invalide")
	}
	
	// Extraire la table et les colonnes
	tablePart := strings.SplitN(insertPart[0], "INTO", 2)
	if len(tablePart) != 2 {
		return nil, fmt.Errorf("syntaxe INSERT invalide")
	}
	
	// Extraire le nom de la table
	table := strings.TrimSpace(tablePart[1])
	
	// Extraire les colonnes
	columns := make([]string, 0)
	if strings.Contains(table, "(") {
		// Extraire la partie entre parenthèses
		start := strings.Index(table, "(")
		end := strings.Index(table, ")")
		if start == -1 || end == -1 {
			return nil, fmt.Errorf("syntaxe INSERT invalide")
		}
		columnsStr := table[start+1 : end]
		columns = strings.Split(columnsStr, ",")
		for i, col := range columns {
			columns[i] = strings.TrimSpace(col)
		}
		// Supprimer la partie entre parenthèses du nom de la table
		table = table[:start]
	}
	
	// Supprimer les points-virgules
	table = strings.Trim(table, ";")
	table = strings.TrimSpace(table)
	
	// Extraire les valeurs
	valuesStr := strings.TrimSpace(insertPart[1])
	values := make([]string, 0)
	if strings.Contains(valuesStr, "(") {
		// Extraire la partie entre parenthèses
		start := strings.Index(valuesStr, "(")
		end := strings.Index(valuesStr, ")")
		if start == -1 || end == -1 {
			return nil, fmt.Errorf("syntaxe INSERT invalide")
		}
		valuesStr := valuesStr[start+1 : end]
		valuesStr = strings.Trim(valuesStr, ";")
		values = strings.Split(valuesStr, ",")
		for i, val := range values {
			values[i] = strings.TrimSpace(val)
		}
	}

	// Afficher les colonnes et valeurs pour débogage
	fmt.Printf("DEBUG: Table: %s\n", table)
	fmt.Printf("DEBUG: Columns: %v\n", columns)
	fmt.Printf("DEBUG: Values: %v\n", values)

	return &types.InsertQuery{
		Table:   table,
		Columns: columns,
		Values:  values,
	}, nil
}

// parseCreate analyse une requête CREATE
func parseCreate(query string) (types.Query, error) {
	// Supprimer les espaces en début et fin
	query = strings.TrimSpace(query)
	
	// Extraire la table et les colonnes
	createPart := strings.SplitN(query, "(", 2)
	if len(createPart) != 2 {
		return nil, fmt.Errorf("syntaxe CREATE invalide")
	}
	
	// Extraire le nom de la table
	table := strings.TrimSpace(createPart[0][12:])
	
	// Extraire les colonnes
	columns := make([]string, 0)
	if strings.Contains(createPart[1], ")") {
		// Extraire la partie entre parenthèses
		end := strings.Index(createPart[1], ")")
		columnsStr := createPart[1][:end]
		columns = strings.Split(columnsStr, ",")
		for i, col := range columns {
			columns[i] = strings.TrimSpace(col)
		}
	}

	// Supprimer les points-virgules
	table = strings.Trim(table, ";")
	table = strings.TrimSpace(table)

	// Afficher les détails pour débogage
	fmt.Printf("DEBUG: CREATE Table: %s\n", table)
	fmt.Printf("DEBUG: CREATE Columns: %v\n", columns)

	return &types.CreateQuery{
		Table:   table,
		Columns: columns,
	}, nil
}

// Fonctions utilitaires
func trimSpace(s string) string {
	// TODO: Implémenter la suppression des espaces
	return s
}

func toUpper(s string) string {
	// TODO: Implémenter la conversion en majuscules
	return s
}
