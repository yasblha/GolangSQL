package main

import (
	"fmt"
	"mydb/internal/db"
	"mydb/pkg/sql"
	"os"
)

func main() {
	// Test du parsing SQL
	testQueries := []string{
		"SELECT id, name FROM users WHERE id = 1",
		"INSERT INTO users (id, name) VALUES (1, 'John')",
	}

	for _, query := range testQueries {
		fmt.Printf("\nTest de la requête: %s\n", query)

		result, err := sql.Parse(query)
		if err != nil {
			fmt.Printf("Erreur de parsing: %v\n", err)
			continue
		}

		fmt.Printf("Table: %s\n", result.GetTable())
		switch q := result.(type) {
		case *sql.SelectQuery:
			fmt.Printf("Colonnes: %v\n", q.Columns)
			fmt.Printf("Conditions: %v\n", q.Conditions)
		case *sql.InsertQuery:
			fmt.Printf("Colonnes: %v\n", q.Columns)
			fmt.Printf("Valeurs: %v\n", q.Values)
		}
	}

	// Lecture du fichier sample.db
	file, err := os.Open("sample.db")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	info := db.ParseHeader(file)
	fmt.Printf("\nInformations sur la base de données:\n")
	fmt.Printf("Page size: %d bytes\n", info.PageSize)

	tables := db.ReadMasterTable(file, info)
	fmt.Printf("\nTables trouvées:\n")
	for _, t := range tables {
		fmt.Printf("Table: %s\n", t.Name)
		for _, col := range t.Columns {
			fmt.Printf(" - %s (%s)\n", col.Name, col.Type)
		}
	}
}
