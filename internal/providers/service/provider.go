package service

import (
	"context"

	"github.com/hse-telescope/core/internal/repository/models"
	"github.com/hse-telescope/tracer"
	"github.com/olegdayo/omniconv"
)

type Repository interface {
	GetService(ctx context.Context, service_id int) (models.Service, error)
	GetGraphServices(ctx context.Context, graph_id int) ([]models.Service, error)
	CreateService(ctx context.Context, service models.Service) (models.Service, error)
	CreateServices(ctx context.Context, graph_id int, services []models.Service) ([]int, error)
	UpdateService(ctx context.Context, service_id int, service models.Service) error
	UpdateGraphServices(ctx context.Context, graph_id int, services []models.Service) error
	DeleteService(ctx context.Context, service_id int) error
}

type Provider struct {
	repository Repository
}

func New(repository Repository) Provider {
	return Provider{
		repository: repository,
	}
}

func (p Provider) GetService(ctx context.Context, service_id int) (Service, error) {
	ctx, span := tracer.Start(ctx, "storage/GetService")
	defer span.End()

	service, err := p.repository.GetService(ctx, service_id)
	if err != nil {
		return Service{}, err
	}
	return DBService2ProviderService(service), nil
}

func (p Provider) GetGraphServices(ctx context.Context, graph_id int) ([]Service, error) {
	ctx, span := tracer.Start(ctx, "storage/GetGraphServices")
	defer span.End()

	services, err := p.repository.GetGraphServices(ctx, graph_id)
	if err != nil {
		return nil, err
	}
	return omniconv.ConvertSlice(services, DBService2ProviderService), nil
}

func (p Provider) CreateService(ctx context.Context, service Service) (Service, error) {
	ctx, span := tracer.Start(ctx, "storage/CreateService")
	defer span.End()

	newservice, err := p.repository.CreateService(ctx, ProviderService2DBService(service))
	return DBService2ProviderService(newservice), err
}

func (p Provider) CreateServices(ctx context.Context, graph_id int, services []Service) ([]int, error) {
	ctx, span := tracer.Start(ctx, "storage/CreateServices")
	defer span.End()

	ids, err := p.repository.CreateServices(ctx, graph_id, omniconv.ConvertSlice(services, ProviderService2DBService))
	return ids, err
}

func (p Provider) UpdateService(ctx context.Context, service_id int, service Service) error {
	ctx, span := tracer.Start(ctx, "storage/UpdateService")
	defer span.End()

	err := p.repository.UpdateService(ctx, service_id, ProviderService2DBService(service))
	return err
}

func (p Provider) UpdateGraphServices(ctx context.Context, graph_id int, services []Service) error {
	ctx, span := tracer.Start(ctx, "storage/UpdateGraphServices")
	defer span.End()

	err := p.repository.UpdateGraphServices(ctx, graph_id, omniconv.ConvertSlice(services, ProviderService2DBService))
	return err
}

func (p Provider) DeleteService(ctx context.Context, service_id int) error {
	ctx, span := tracer.Start(ctx, "storage/DeleteService")
	defer span.End()

	err := p.repository.DeleteService(ctx, service_id)
	return err
}
