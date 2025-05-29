package relation

import (
	"context"

	"github.com/hse-telescope/core/internal/repository/models"
	"github.com/olegdayo/omniconv"
)

type Repository interface {
	GetRelationGraph(ctx context.Context, relation_id int) (int, error)
	GetRelation(ctx context.Context, relation_id int) (models.Relation, error)
	GetGraphRelations(ctx context.Context, graph_id int) ([]models.Relation, error)
	CreateRelation(ctx context.Context, relation models.Relation) (models.Relation, error)
	CreateRelations(ctx context.Context, graph_id int, relations []models.Relation) error
	UpdateRelation(ctx context.Context, relation_id int, relation models.Relation) error
	UpdateGraphRelations(ctx context.Context, graph_id int, relations []models.Relation) error
	DeleteRelation(ctx context.Context, relation_id int) error
}

type Provider struct {
	repository Repository
}

func New(repository Repository) Provider {
	return Provider{
		repository: repository,
	}
}

func (p Provider) GetRelationGraph(ctx context.Context, relation_id int) (int, error) {
	graph_id, err := p.repository.GetRelationGraph(ctx, relation_id)
	if err != nil {
		return -1, err
	}
	return graph_id, nil
}

func (p Provider) GetRelation(ctx context.Context, relation_id int) (Relation, error) {
	relation, err := p.repository.GetRelation(ctx, relation_id)
	if err != nil {
		return Relation{}, err
	}
	return DBRelation2ProviderRelation(relation), nil
}

func (p Provider) GetGraphRelations(ctx context.Context, graph_id int) ([]Relation, error) {
	relations, err := p.repository.GetGraphRelations(ctx, graph_id)
	if err != nil {
		return nil, err
	}
	return omniconv.ConvertSlice(relations, DBRelation2ProviderRelation), nil
}

func (p Provider) CreateRelation(ctx context.Context, relation Relation) (Relation, error) {
	newrelation, err := p.repository.CreateRelation(ctx, ProviderRelation2DBRelation(relation))
	return DBRelation2ProviderRelation(newrelation), err
}

func (p Provider) CreateRelations(ctx context.Context, graph_id int, relations []Relation) error {
	err := p.repository.CreateRelations(ctx, graph_id, omniconv.ConvertSlice(relations, ProviderRelation2DBRelation))
	return err
}

func (p Provider) UpdateRelation(ctx context.Context, relation_id int, relation Relation) error {
	err := p.repository.UpdateRelation(ctx, relation_id, ProviderRelation2DBRelation(relation))
	return err
}

func (p Provider) UpdateGraphRelations(ctx context.Context, graph_id int, relations []Relation) error {
	err := p.repository.UpdateGraphRelations(ctx, graph_id, omniconv.ConvertSlice(relations, ProviderRelation2DBRelation))
	return err
}

func (p Provider) DeleteRelation(ctx context.Context, relation_id int) error {
	err := p.repository.DeleteRelation(ctx, relation_id)
	return err
}
