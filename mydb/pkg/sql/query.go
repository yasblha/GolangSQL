package sql

import (
	"fmt"
)

// Parser représente un parser SQL
type Parser interface {
	Parse(query string) (Query, error)
}

// Parse analyse une requête SQL et retourne une Query
func Parse(query string) (Query, error) {
	// Déterminer le type de requête
	queryType := getQueryType(query)
	
	switch queryType {
	case "SELECT":
		return parseSelect(query)
	case "INSERT":
		return parseInsert(query)
	default:
		return nil, fmt.Errorf("type de requête non supporté: %s", queryType)
	}
}

// getQueryType détermine le type de requête SQL
func getQueryType(query string) string {
	// Normaliser la requête
	query = normalizeQuery(query)
	
	// Extraire le premier mot
	words := splitWords(query)
	if len(words) == 0 {
		return ""
	}
	
	return words[0]
}

// normalizeQuery normalise une requête SQL
func normalizeQuery(query string) string {
	// Supprimer les espaces en début et fin
	query = trimSpace(query)
	
	// Convertir en majuscules
	query = toUpper(query)
	
	return query
}

// splitWords divise une requête en mots
func splitWords(query string) []string {
	// TODO: Implémenter la division correcte des mots
	return []string{}
}

// parseSelect analyse une requête SELECT
func parseSelect(query string) (Query, error) {
	// TODO: Implémenter l'analyse des requêtes SELECT
	return nil, nil
}

// parseInsert analyse une requête INSERT
func parseInsert(query string) (Query, error) {
	// TODO: Implémenter l'analyse des requêtes INSERT
	return nil, nil
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