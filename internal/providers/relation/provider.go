package relation

import (
	"context"

	"github.com/hse-telescope/core/internal/repository/models"
	"github.com/hse-telescope/tracer"
	"github.com/olegdayo/omniconv"
)

type Repository interface {
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

func (p Provider) GetRelation(ctx context.Context, relation_id int) (Relation, error) {
	ctx, span := tracer.Start(ctx, "provider/GetRelation")
	defer span.End()

	relation, err := p.repository.GetRelation(ctx, relation_id)
	if err != nil {
		return Relation{}, err
	}
	return DBRelation2ProviderRelation(relation), nil
}

func (p Provider) GetGraphRelations(ctx context.Context, graph_id int) ([]Relation, error) {
	ctx, span := tracer.Start(ctx, "provider/GetGraphRelations")
	defer span.End()

	relations, err := p.repository.GetGraphRelations(ctx, graph_id)
	if err != nil {
		return nil, err
	}
	return omniconv.ConvertSlice(relations, DBRelation2ProviderRelation), nil
}

func (p Provider) CreateRelation(ctx context.Context, relation Relation) (Relation, error) {
	ctx, span := tracer.Start(ctx, "provider/CreateRelation")
	defer span.End()

	newrelation, err := p.repository.CreateRelation(ctx, ProviderRelation2DBRelation(relation))
	return DBRelation2ProviderRelation(newrelation), err
}

func (p Provider) CreateRelations(ctx context.Context, graph_id int, relations []Relation) error {
	ctx, span := tracer.Start(ctx, "provider/CreateRelations")
	defer span.End()

	err := p.repository.CreateRelations(ctx, graph_id, omniconv.ConvertSlice(relations, ProviderRelation2DBRelation))
	return err
}

func (p Provider) UpdateRelation(ctx context.Context, relation_id int, relation Relation) error {
	ctx, span := tracer.Start(ctx, "provider/UpdateRelation")
	defer span.End()

	err := p.repository.UpdateRelation(ctx, relation_id, ProviderRelation2DBRelation(relation))
	return err
}

func (p Provider) UpdateGraphRelations(ctx context.Context, graph_id int, relations []Relation) error {
	ctx, span := tracer.Start(ctx, "provider/UpdateGraphRelations")
	defer span.End()

	err := p.repository.UpdateGraphRelations(ctx, graph_id, omniconv.ConvertSlice(relations, ProviderRelation2DBRelation))
	return err
}

func (p Provider) DeleteRelation(ctx context.Context, relation_id int) error {
	ctx, span := tracer.Start(ctx, "provider/DeleteRelation")
	defer span.End()

	err := p.repository.DeleteRelation(ctx, relation_id)
	return err
}
