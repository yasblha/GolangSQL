package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"mydb/internal/engine"
	"mydb/pkg/sql"
)

func main() {
	fmt.Println("GoSQL CLI - Tapez 'exit' pour quitter")
	fmt.Println("-----------------------------------")

	// Créer une instance du moteur
	engine := engine.New()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("gosql> ")
		if !scanner.Scan() {
			break
		}

		query := strings.TrimSpace(scanner.Text())
		if query == "exit" {
			break
		}

		if query == "" {
			continue
		}

		// Parser la requête
		result, err := sql.Parse(query)
		if err != nil {
			fmt.Printf("Erreur de parsing: %v\n", err)
			continue
		}

		// Exécuter la requête
		output, err := engine.Execute(result)
		if err != nil {
			fmt.Printf("Erreur d'exécution: %v\n", err)
			continue
		}

		// Afficher les résultats
		fmt.Println(output)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Erreur de lecture:", err)
	}
}
