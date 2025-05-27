package main

import (
	"fmt"
	"mydb/db"
	"mydb/parser"
	"os"
	"strings"
)

func main() {
	// Test du parsing SQL
	testQueries := []string{
		"SELECT id, name FROM users WHERE id = 1",
		"INSERT INTO users (id, name) VALUES (1, 'John')",
	}

	for _, query := range testQueries {
		fmt.Printf("\nTest de la requête: %s\n", query)
		
		if strings.HasPrefix(strings.ToUpper(query), "SELECT") {
			result, err := parser.ParseSelect(query)
			if err != nil {
				fmt.Printf("Erreur de parsing SELECT: %v\n", err)
				continue
			}
			fmt.Printf("Table: %s\n", result.Table)
			fmt.Printf("Colonnes: %v\n", result.Columns)
			fmt.Printf("Conditions: %v\n", result.Conditions)
		} else if strings.HasPrefix(strings.ToUpper(query), "INSERT") {
			result, err := parser.ParseInsert(query)
			if err != nil {
				fmt.Printf("Erreur de parsing INSERT: %v\n", err)
				continue
			}
			fmt.Printf("Table: %s\n", result.Table)
			fmt.Printf("Colonnes: %v\n", result.Columns)
			fmt.Printf("Valeurs: %v\n", result.Values)
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
