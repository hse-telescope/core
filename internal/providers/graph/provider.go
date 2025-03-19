package graph

import (
	"context"

	"github.com/hse-telescope/core/internal/repository/models"
	"github.com/olegdayo/omniconv"
)

type Repository interface {
	CreateGraph(ctx context.Context, graph models.Graph) (models.Graph, error)
	DeleteGraph(ctx context.Context, graph_id int) error
	UpdateGraph(ctx context.Context, graph_id int, graph models.Graph) error
	GetProjectGraphs(ctx context.Context, project_id int) ([]models.Graph, error)
}

type Provider struct {
	repository Repository
}

func New(repository Repository) Provider {
	return Provider{
		repository: repository,
	}
}

func (p Provider) CreateGraph(ctx context.Context, graph Graph) (Graph, error) {
	newgraph, err := p.repository.CreateGraph(ctx, ProviderGraph2DBGraph(graph))
	return DBGraph2ProviderGraph(newgraph), err
}

func (p Provider) DeleteGraph(ctx context.Context, graph_id int) error {
	err := p.repository.DeleteGraph(ctx, graph_id)
	return err
}

func (p Provider) UpdateGraph(ctx context.Context, graph_id int, graph Graph) error {
	err := p.repository.UpdateGraph(ctx, graph_id, ProviderGraph2DBGraph(graph))
	return err
}

func (p Provider) GetProjectGraphs(ctx context.Context, project_id int) ([]Graph, error) {
	graphs, err := p.repository.GetProjectGraphs(ctx, project_id)
	if err != nil {
		return nil, err
	}
	return omniconv.ConvertSlice(graphs, DBGraph2ProviderGraph), nil
}
