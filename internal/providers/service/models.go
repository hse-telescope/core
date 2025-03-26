package service

import (
	"github.com/hse-telescope/core/internal/repository/models"
)

type Service struct {
	ID          int
	GraphID     int
	Name        string
	Description string
	X           float32
	Y           float32
}

func ProviderService2DBService(service Service) models.Service {
	return models.Service{
		ID:          service.ID,
		GraphID:     service.GraphID,
		Name:        service.Name,
		Description: service.Description,
		X:           service.X,
		Y:           service.Y,
	}
}

func DBService2ProviderService(service models.Service) Service {
	return Service{
		ID:          service.ID,
		GraphID:     service.GraphID,
		Name:        service.Name,
		Description: service.Description,
		X:           service.X,
		Y:           service.Y,
	}
}
