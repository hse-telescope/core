package graph

import "github.com/hse-telescope/core/internal/repository/models"

type Graph struct {
	ID        int
	ProjectID int
	Name      string
}

func ProviderGraph2DBGraph(graph Graph) models.Graph {
	return models.Graph{
		ID:        graph.ID,
		ProjectID: graph.ProjectID,
		Name:      graph.Name,
	}
}

func DBGraph2ProviderGraph(graph models.Graph) Graph {
	return Graph{
		ID:        graph.ID,
		ProjectID: graph.ProjectID,
		Name:      graph.Name,
	}
}
