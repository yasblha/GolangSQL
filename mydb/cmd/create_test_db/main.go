package main

import (
	"fmt"
	"mydb/tests"
)

func main() {
	err := tests.CreateTestDB("test.db")
	if err != nil {
		fmt.Printf("Erreur lors de la création du fichier de test: %v\n", err)
		return
	}
	fmt.Println("Fichier de test créé avec succès: test.db")
} 