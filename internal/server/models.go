package server

import (
	"github.com/hse-telescope/core/internal/providers/graph"
	"github.com/hse-telescope/core/internal/providers/project"
	"github.com/hse-telescope/core/internal/providers/relation"
	"github.com/hse-telescope/core/internal/providers/service"
)

type Project struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Graph struct {
	ID        int    `json:"id"`
	ProjectID int    `json:"project_id"`
	Name      string `json:"name"`
}

type Service struct {
	ID          int     `json:"id"`
	GraphID     int     `json:"graph_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	X           float32 `json:"x"`
	Y           float32 `json:"y"`
}

type Relation struct {
	ID          int    `json:"id"`
	GraphID     int    `json:"graph_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	FromService int    `json:"from_service"`
	ToService   int    `json:"to_service"`
}

func ServerProject2ProviderProject(pr Project) project.Project {
	return project.Project{
		ID:   pr.ID,
		Name: pr.Name,
	}
}

func ProviderProject2ServerProject(pr project.Project) Project {
	return Project{
		ID:   pr.ID,
		Name: pr.Name,
	}
}

func ServerGraph2ProviderGraph(gr Graph) graph.Graph {
	return graph.Graph{
		ID:        gr.ID,
		ProjectID: gr.ProjectID,
		Name:      gr.Name,
	}
}

func ProviderGraph2ServerGraph(gr graph.Graph) Graph {
	return Graph{
		ID:        gr.ID,
		ProjectID: gr.ProjectID,
		Name:      gr.Name,
	}
}

func ServerService2ProviderService(serv Service) service.Service {
	return service.Service{
		ID:          serv.ID,
		GraphID:     serv.GraphID,
		Name:        serv.Name,
		Description: serv.Description,
		X:           serv.X,
		Y:           serv.Y,
	}
}

func ProviderService2ServerService(serv service.Service) Service {
	return Service{
		ID:          serv.ID,
		GraphID:     serv.GraphID,
		Name:        serv.Name,
		Description: serv.Description,
		X:           serv.X,
		Y:           serv.Y,
	}
}

func ServerRelation2ProviderRelation(rel Relation) relation.Relation {
	return relation.Relation{
		ID:          rel.ID,
		GraphID:     rel.GraphID,
		Name:        rel.Name,
		Description: rel.Description,
		FromService: rel.FromService,
		ToService:   rel.ToService,
	}
}

func ProviderRelation2ServerRelation(rel relation.Relation) Relation {
	return Relation{
		ID:          rel.ID,
		GraphID:     rel.GraphID,
		Name:        rel.Name,
		Description: rel.Description,
		FromService: rel.FromService,
		ToService:   rel.ToService,
	}
}
