package engine

import (
	"mydb/pkg/types"
)

// Engine est une implémentation de types.Engine
// Il contient les méthodes d'exécution des requêtes
type Engine struct {
	*types.Engine
}

// New crée une nouvelle instance du moteur
func New() *Engine {
	return &Engine{
		Engine: types.New(),
	}
}

// Execute exécute une requête
func (e *Engine) Execute(query types.Query) (string, error) {
	return query.Execute(e.Engine)
}
