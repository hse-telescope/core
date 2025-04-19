package relation

import "github.com/hse-telescope/core/internal/repository/models"

type Relation struct {
	ID          int
	GraphID     int
	Name        string
	Description string
	FromService int
	ToService   int
}

func ProviderRelation2DBRelation(relation Relation) models.Relation {
	return models.Relation{
		ID:          relation.ID,
		GraphID:     relation.GraphID,
		Name:        relation.Name,
		Description: relation.Description,
		FromService: relation.FromService,
		ToService:   relation.ToService,
	}
}

func DBRelation2ProviderRelation(relation models.Relation) Relation {
	return Relation{
		ID:          relation.ID,
		GraphID:     relation.GraphID,
		Name:        relation.Name,
		Description: relation.Description,
		FromService: relation.FromService,
		ToService:   relation.ToService,
	}
}
