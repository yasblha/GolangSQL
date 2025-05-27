package main

import (
	"fmt"
	"log"

	"mydb/pkg/sql"
)

func main() {
	// Exemple de requête SELECT
	selectQuery := "SELECT name, age FROM users WHERE age > 25"
	query, err := sql.Parse(selectQuery)
	if err != nil {
		log.Fatalf("Erreur lors de l'analyse de la requête SELECT: %v", err)
	}

	// Vérifier le type de requête
	switch q := query.(type) {
	case *sql.SelectQuery:
		fmt.Printf("Requête SELECT analysée:\n")
		fmt.Printf("  Table: %s\n", q.Table)
		fmt.Printf("  Colonnes: %v\n", q.Columns)
		fmt.Printf("  Conditions: %v\n", q.Conditions)
	default:
		fmt.Printf("Type de requête inattendu: %T\n", query)
	}

	// Exemple de requête INSERT
	insertQuery := "INSERT INTO users (name, age) VALUES ('John', 25)"
	query, err = sql.Parse(insertQuery)
	if err != nil {
		log.Fatalf("Erreur lors de l'analyse de la requête INSERT: %v", err)
	}

	// Vérifier le type de requête
	switch q := query.(type) {
	case *sql.InsertQuery:
		fmt.Printf("\nRequête INSERT analysée:\n")
		fmt.Printf("  Table: %s\n", q.Table)
		fmt.Printf("  Colonnes: %v\n", q.Columns)
		fmt.Printf("  Valeurs: %v\n", q.Values)
	default:
		fmt.Printf("Type de requête inattendu: %T\n", query)
	}
} 