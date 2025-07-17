package types

import (
	"fmt"
	"strings"
	"strconv"
)



// Engine représente le moteur de base de données
type Engine struct {
	Tables     map[string][]map[string]interface{}
	NextID     map[string]int
}

// New crée une nouvelle instance de Engine
func New() *Engine {
	return &Engine{
		Tables:  make(map[string][]map[string]interface{}),
		NextID:  make(map[string]int),
	}
}

// Execute implémente l'interface Query pour SelectQuery
func (q *SelectQuery) Execute(engine *Engine) (string, error) {
	// Vérifier si la table existe
	rows, exists := engine.Tables[q.Table]
	if !exists {
		return "", fmt.Errorf("table '%s' n'existe pas", q.Table)
	}

	if len(rows) == 0 {
		return "Aucun résultat trouvé", nil
	}

	// Afficher l'état de la table pour débogage
	fmt.Printf("DEBUG: Table state before select:\n")
	for _, r := range rows {
		fmt.Printf("DEBUG: Row: %v\n", r)
	}

	// Déterminer les colonnes à afficher
	var columns []string
	if len(q.Columns) == 1 && q.Columns[0] == "*" {
		// Si SELECT *, utiliser toutes les colonnes
		for col := range rows[0] {
			columns = append(columns, col)
		}
	} else {
		columns = q.Columns
	}

	// Construire l'en-tête
	var header strings.Builder
	var separator strings.Builder
	for _, col := range columns {
		header.WriteString(fmt.Sprintf("%-15s", col))
		separator.WriteString(strings.Repeat("-", 15))
	}
	header.WriteString("\n")
	separator.WriteString("\n")

	// Construire les lignes
	var result strings.Builder
	result.WriteString(header.String())
	result.WriteString(separator.String())

	// Afficher les données
	for _, row := range rows {
		// Vérifier si la ligne est vide
		isEmpty := true
		for _, col := range columns {
			if row[col] != nil {
				isEmpty = false
				break
			}
		}
		if isEmpty {
			continue
		}

		// Si il y a des conditions WHERE, vérifier si la ligne correspond
		if len(q.Conditions) > 0 {
			matches := true
			for _, cond := range q.Conditions {
				// Vérifier si la colonne existe dans la ligne
				if value, ok := row[cond.Column]; ok {
					// Comparer selon le type de la valeur
					switch v := value.(type) {
					case int:
						if numVal, err := strconv.Atoi(cond.Value.(string)); err != nil || v != numVal {
							matches = false
							break
						}
					case string:
						if v != cond.Value {
							matches = false
							break
						}
					default:
						matches = false
						break
					}
				} else {
					matches = false
					break
				}
				if !matches {
					break
				}
			}
			if !matches {
				continue
			}
		}

		// Afficher les valeurs pour débogage
		fmt.Printf("DEBUG: Processing row:\n")
		for _, col := range columns {
			value := row[col]
			fmt.Printf("DEBUG: Column %s = %v (type %T)\n", col, value, value)
		}

		// Afficher les données dans le bon ordre
		for _, col := range columns {
			value := row[col]
			switch v := value.(type) {
			case int:
				result.WriteString(fmt.Sprintf("%-15d", v)) // Aligner à gauche pour les nombres
			case string:
				result.WriteString(fmt.Sprintf("%-15s", v)) // Aligner à gauche pour les textes
			default:
				result.WriteString(fmt.Sprintf("%-15v", v)) // Formatage par défaut
			}
		}
		result.WriteString("\n")
	}

	return result.String(), nil
}

// Execute implémente l'interface Query pour InsertQuery
func (q *InsertQuery) Execute(engine *Engine) (string, error) {
	// Vérifier si la table existe
	if _, exists := engine.Tables[q.Table]; !exists {
		return "", fmt.Errorf("table '%s' n'existe pas", q.Table)
	}

	// Créer une nouvelle ligne
	row := make(map[string]interface{})

	// Ajouter l'ID auto-incrémenté
	id := engine.NextID[q.Table]
	row["id"] = id + 1
	engine.NextID[q.Table] = id + 1

	// Ajouter les autres colonnes
	for i, col := range q.Columns {
		if i < len(q.Values) {
			// Nettoyer la valeur (enlever les guillemets)
			value := strings.Trim(q.Values[i], "'")
			// Convertir en entier si possible, sinon garder en string
			if num, err := strconv.Atoi(value); err == nil {
				row[col] = num
			} else {
				row[col] = value
			}
			fmt.Printf("DEBUG: Inserting %s=%v (type %T)\n", col, row[col], row[col])
		}
	}

	// Ajouter la ligne à la table
	engine.Tables[q.Table] = append(engine.Tables[q.Table], row)

	// Afficher l'état de la table pour débogage
	fmt.Printf("DEBUG: Table state after insert:\n")
	for _, r := range engine.Tables[q.Table] {
		fmt.Printf("DEBUG: Row: %v\n", r)
	}

	return "Ligne insérée avec succès", nil
}

// Execute implémente l'interface Query pour CreateQuery
func (q *CreateQuery) Execute(engine *Engine) (string, error) {
	// Vérifier si la table existe déjà
	if _, exists := engine.Tables[q.Table]; exists {
		return "", fmt.Errorf("table '%s' existe déjà", q.Table)
	}

	// Initialiser la table avec une colonne id
	engine.Tables[q.Table] = make([]map[string]interface{}, 0)
	engine.NextID[q.Table] = 0

	// Ajouter les autres colonnes
	for _, col := range q.Columns {
		// Pour le moment, on gère toutes les colonnes comme des strings
		// On pourrait ajouter la gestion des types plus tard
		for _, row := range engine.Tables[q.Table] {
			row[col] = ""
		}
	}

	return fmt.Sprintf("Table '%s' créée avec succès", q.Table), nil
}

// Query représente une requête SQL générique
type Query interface {
	GetType() string
	GetTable() string
	Execute(engine *Engine) (string, error)
}

// SelectQuery représente une requête SELECT
type SelectQuery struct {
	Columns    []string
	Table      string
	Conditions []Condition
}

// InsertQuery représente une requête INSERT
type InsertQuery struct {
	Table      string
	Columns    []string
	Values     []string
}

// CreateQuery représente une requête CREATE TABLE
type CreateQuery struct {
	Table   string
	Columns []string
}

// Condition représente une condition WHERE
type Condition struct {
	Column    string
	Operator  string
	Value     interface{}
	IsNull    bool
	IsNotNull bool
}

// Implémentation des méthodes de l'interface Query pour SelectQuery
func (q *SelectQuery) GetType() string { return "SELECT" }
func (q *SelectQuery) GetTable() string { return q.Table }

// Implémentation des méthodes de l'interface Query pour InsertQuery
func (q *InsertQuery) GetType() string { return "INSERT" }
func (q *InsertQuery) GetTable() string { return q.Table }

// Implémentation des méthodes de l'interface Query pour CreateQuery
func (q *CreateQuery) GetType() string { return "CREATE" }
func (q *CreateQuery) GetTable() string { return q.Table }
