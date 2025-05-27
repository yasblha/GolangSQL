package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"mydb/db"
	"mydb/parser"
)

// Structure pour stocker les données en mémoire
var dataStore = make(map[string][]map[string]interface{})

func executeCreateTable(query string) error {
	// Extraire le nom de la table et les colonnes
	parts := strings.Split(query, "(")
	if len(parts) != 2 {
		return fmt.Errorf("syntaxe CREATE TABLE invalide")
	}

	// Extraire le nom de la table
	tableName := strings.TrimSpace(strings.Replace(parts[0], "CREATE TABLE", "", -1))
	
	// Extraire les définitions des colonnes
	columnsDef := strings.TrimRight(parts[1], ")")
	columns := strings.Split(columnsDef, ",")

	// Créer la table dans le store
	dataStore[strings.ToLower(tableName)] = make([]map[string]interface{}, 0)
	
	fmt.Printf("Table '%s' créée avec succès\n", tableName)
	fmt.Println("Colonnes:")
	for _, col := range columns {
		fmt.Printf("  * %s\n", strings.TrimSpace(col))
	}
	
	return nil
}

func executeInsert(insert *parser.InsertQuery) error {
	// Vérifier si la table existe
	if _, exists := dataStore[insert.Table]; !exists {
		dataStore[insert.Table] = make([]map[string]interface{}, 0)
	}

	// Créer une nouvelle ligne
	row := make(map[string]interface{})
	for i, col := range insert.Columns {
		if i < len(insert.Values) {
			row[col] = insert.Values[i]
		}
	}

	// Ajouter la ligne à la table
	dataStore[insert.Table] = append(dataStore[insert.Table], row)
	return nil
}

func executeSelect(selectQuery *parser.SelectQuery) {
	table := strings.ToLower(selectQuery.Table)
	if rows, exists := dataStore[table]; exists {
		fmt.Println("\nRésultats:")
		fmt.Println("----------")
		
		// Afficher les en-têtes
		if len(rows) > 0 {
			headers := make([]string, 0)
			if selectQuery.Columns[0] == "*" {
				for col := range rows[0] {
					headers = append(headers, col)
				}
			} else {
				headers = selectQuery.Columns
			}
			
			// Afficher les en-têtes
			for _, h := range headers {
				fmt.Printf("%-15s", h)
			}
			fmt.Println("\n----------")

			// Afficher les données
			for _, row := range rows {
				// Vérifier les conditions WHERE si elles existent
				match := true
				for _, condition := range selectQuery.Conditions {
					val := row[condition.Column]
					if val != condition.Value {
						match = false
						break
					}
				}
				if !match {
					continue
				}

				// Afficher les colonnes demandées
				if selectQuery.Columns[0] == "*" {
					for _, h := range headers {
						fmt.Printf("%-15v", row[h])
					}
				} else {
					for _, col := range selectQuery.Columns {
						fmt.Printf("%-15v", row[col])
					}
				}
				fmt.Println()
			}
		} else {
			fmt.Println("Aucun résultat trouvé")
		}
	} else {
		fmt.Printf("Table '%s' non trouvée\n", table)
	}
	fmt.Println()
}

func main() {
	// Ouvrir le fichier de base de données
	file, err := os.OpenFile("test.db", os.O_RDWR, 0644)
	if err != nil {
		fmt.Printf("Erreur d'ouverture du fichier: %v\n", err)
		return
	}
	defer file.Close()

	// Lire l'en-tête
	info := db.ParseHeader(file)
	fmt.Printf("Taille de page: %d bytes\n", info.PageSize)

	// Lire la table master
	tables := db.ReadMasterTable(file, info)
	fmt.Println("\nTables trouvées:")
	for _, table := range tables {
		fmt.Printf("- %s\n", table.Name)
		for _, col := range table.Columns {
			fmt.Printf("  * %s (%s)\n", col.Name, col.Type)
		}
		// Initialiser la table dans le store
		dataStore[strings.ToLower(table.Name)] = make([]map[string]interface{}, 0)
	}

	// Boucle principale pour les commandes
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("\nEntrez vos commandes SQL (ou 'exit' pour quitter):")
	
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}

		query := strings.TrimSpace(scanner.Text())
		if query == "exit" {
			break
		}

		// Parser la requête
		if strings.HasPrefix(strings.ToUpper(query), "CREATE TABLE") {
			if err := executeCreateTable(query); err != nil {
				fmt.Printf("Erreur de création de table: %v\n", err)
			}
		} else if strings.HasPrefix(strings.ToUpper(query), "SELECT") {
			selectQuery, err := parser.ParseSelect(query)
			if err != nil {
				fmt.Printf("Erreur de parsing SELECT: %v\n", err)
				continue
			}
			executeSelect(selectQuery)

		} else if strings.HasPrefix(strings.ToUpper(query), "INSERT") {
			insertQuery, err := parser.ParseInsert(query)
			if err != nil {
				fmt.Printf("Erreur de parsing INSERT: %v\n", err)
				continue
			}
			if err := executeInsert(insertQuery); err != nil {
				fmt.Printf("Erreur d'exécution INSERT: %v\n", err)
				continue
			}
			fmt.Println("Insertion réussie")

		} else {
			fmt.Println("Commande non reconnue. Utilisez CREATE TABLE, SELECT ou INSERT.")
		}
	}
} 