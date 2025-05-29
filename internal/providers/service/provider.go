package service

import (
	"context"

	"github.com/hse-telescope/core/internal/repository/models"
	"github.com/olegdayo/omniconv"
)

type Repository interface {
	GetServiceGraph(ctx context.Context, service_id int) (int, error)
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

func (p Provider) GetServiceGraph(ctx context.Context, service_id int) (int, error) {
	graph_id, err := p.repository.GetServiceGraph(ctx, service_id)
	if err != nil {
		return -1, err
	}
	return graph_id, nil
}

func (p Provider) GetService(ctx context.Context, service_id int) (Service, error) {
	service, err := p.repository.GetService(ctx, service_id)
	if err != nil {
		return Service{}, err
	}
	return DBService2ProviderService(service), nil
}

func (p Provider) GetGraphServices(ctx context.Context, graph_id int) ([]Service, error) {
	services, err := p.repository.GetGraphServices(ctx, graph_id)
	if err != nil {
		return nil, err
	}
	return omniconv.ConvertSlice(services, DBService2ProviderService), nil
}

func (p Provider) CreateService(ctx context.Context, service Service) (Service, error) {
	newservice, err := p.repository.CreateService(ctx, ProviderService2DBService(service))
	return DBService2ProviderService(newservice), err
}

func (p Provider) CreateServices(ctx context.Context, graph_id int, services []Service) ([]int, error) {
	ids, err := p.repository.CreateServices(ctx, graph_id, omniconv.ConvertSlice(services, ProviderService2DBService))
	return ids, err
}

func (p Provider) UpdateService(ctx context.Context, service_id int, service Service) error {
	err := p.repository.UpdateService(ctx, service_id, ProviderService2DBService(service))
	return err
}

func (p Provider) UpdateGraphServices(ctx context.Context, graph_id int, services []Service) error {
	err := p.repository.UpdateGraphServices(ctx, graph_id, omniconv.ConvertSlice(services, ProviderService2DBService))
	return err
}

func (p Provider) DeleteService(ctx context.Context, service_id int) error {
	err := p.repository.DeleteService(ctx, service_id)
	return err
}
