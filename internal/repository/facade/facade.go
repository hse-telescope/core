package facade

import (
	"context"

	"github.com/hse-telescope/core/internal/repository/models"
)

type Storage interface {
	GetProjects(ctx context.Context) ([]models.Project, error)
	CreateProject(ctx context.Context, project models.Project) (models.Project, error)
	UpdateProject(ctx context.Context, project_id int, project models.Project) error
	DeleteProject(ctx context.Context, project_id int) error

	CreateGraph(ctx context.Context, graph models.Graph) (models.Graph, error)
	DeleteGraph(ctx context.Context, graph_id int) error
	UpdateGraph(ctx context.Context, graph_id int, graph models.Graph) error
	GetProjectGraphs(ctx context.Context, project_id int) ([]models.Graph, error)

	GetService(ctx context.Context, service_id int) (models.Service, error)
	GetGraphServices(ctx context.Context, graph_id int) ([]models.Service, error)
	CreateService(ctx context.Context, service models.Service) (models.Service, error)
	CreateServices(ctx context.Context, graph_id int, services []models.Service) ([]int, error)
	UpdateService(ctx context.Context, service_id int, service models.Service) error
	DeleteService(ctx context.Context, service_id int) error

	GetRelation(ctx context.Context, relation_id int) (models.Relation, error)
	GetGraphRelations(ctx context.Context, graph_id int) ([]models.Relation, error)
	CreateRelation(ctx context.Context, relation models.Relation) (models.Relation, error)
	CreateRelations(ctx context.Context, graph_id int, relations []models.Relation) error
	UpdateRelation(ctx context.Context, relation_id int, relation models.Relation) error
	DeleteRelation(ctx context.Context, relation_id int) error
}

type Facade struct {
	storage Storage
}

func New(storage Storage) Facade {
	return Facade{
		storage: storage,
	}
}

func (f Facade) GetProjects(ctx context.Context) ([]models.Project, error) {
	return f.storage.GetProjects(ctx)
}

func (f Facade) CreateProject(ctx context.Context, project models.Project) (models.Project, error) {
	return f.storage.CreateProject(ctx, project)
}

func (f Facade) UpdateProject(ctx context.Context, project_id int, project models.Project) error {
	return f.storage.UpdateProject(ctx, project_id, project)
}

func (f Facade) DeleteProject(ctx context.Context, project_id int) error {
	return f.storage.DeleteProject(ctx, project_id)
}

func (f Facade) CreateGraph(ctx context.Context, graph models.Graph) (models.Graph, error) {
	return f.storage.CreateGraph(ctx, graph)
}

func (f Facade) DeleteGraph(ctx context.Context, graph_id int) error {
	return f.storage.DeleteGraph(ctx, graph_id)
}

func (f Facade) UpdateGraph(ctx context.Context, graph_id int, graph models.Graph) error {
	return f.storage.UpdateGraph(ctx, graph_id, graph)
}

func (f Facade) GetProjectGraphs(ctx context.Context, project_id int) ([]models.Graph, error) {
	return f.storage.GetProjectGraphs(ctx, project_id)
}

func (f Facade) GetService(ctx context.Context, service_id int) (models.Service, error) {
	return f.storage.GetService(ctx, service_id)
}

func (f Facade) GetGraphServices(ctx context.Context, graph_id int) ([]models.Service, error) {
	return f.storage.GetGraphServices(ctx, graph_id)
}

func (f Facade) CreateService(ctx context.Context, service models.Service) (models.Service, error) {
	return f.storage.CreateService(ctx, service)
}

func (f Facade) CreateServices(ctx context.Context, graph_id int, services []models.Service) ([]int, error) {
	return f.storage.CreateServices(ctx, graph_id, services)
}

func (f Facade) UpdateService(ctx context.Context, service_id int, service models.Service) error {
	return f.storage.UpdateService(ctx, service_id, service)
}

func (f Facade) DeleteService(ctx context.Context, service_id int) error {
	return f.storage.DeleteService(ctx, service_id)
}

func (f Facade) GetRelation(ctx context.Context, relation_id int) (models.Relation, error) {
	return f.storage.GetRelation(ctx, relation_id)
}

func (f Facade) GetGraphRelations(ctx context.Context, graph_id int) ([]models.Relation, error) {
	return f.storage.GetGraphRelations(ctx, graph_id)
}

func (f Facade) CreateRelation(ctx context.Context, relation models.Relation) (models.Relation, error) {
	return f.storage.CreateRelation(ctx, relation)
}

func (f Facade) UpdateRelation(ctx context.Context, relation_id int, relation models.Relation) error {
	return f.storage.UpdateRelation(ctx, relation_id, relation)
}

func (f Facade) CreateRelations(ctx context.Context, graph_id int, relations []models.Relation) error {
	return f.storage.CreateRelations(ctx, graph_id, relations)
}

func (f Facade) DeleteRelation(ctx context.Context, relation_id int) error {
	return f.storage.DeleteRelation(ctx, relation_id)
}
