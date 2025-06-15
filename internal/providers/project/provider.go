package project

import (
	"context"

	"github.com/hse-telescope/core/internal/repository/models"
	"github.com/hse-telescope/tracer"
	"github.com/olegdayo/omniconv"
)

type Repository interface {
	GetProjects(ctx context.Context) ([]models.Project, error)
	CreateProject(ctx context.Context, project models.Project) (models.Project, error)
	UpdateProject(ctx context.Context, project_id int, project models.Project) error
	DeleteProject(ctx context.Context, project_id int) error
}

type Provider struct {
	repository Repository
}

func New(repository Repository) Provider {
	return Provider{
		repository: repository,
	}
}

func (p Provider) GetProjects(ctx context.Context) ([]Project, error) {
	ctx, span := tracer.Start(ctx, "provider/GetProjects")
	defer span.End()

	projects, err := p.repository.GetProjects(ctx)
	if err != nil {
		return nil, err
	}
	return omniconv.ConvertSlice(projects, DBProject2ProviderProject), nil
}

func (p Provider) CreateProject(ctx context.Context, project Project) (Project, error) {
	ctx, span := tracer.Start(ctx, "provider/CreateProject")
	defer span.End()

	newproject, err := p.repository.CreateProject(ctx, ProviderProject2DBProject(project))
	return DBProject2ProviderProject(newproject), err
}

func (p Provider) UpdateProject(ctx context.Context, project_id int, project Project) error {
	ctx, span := tracer.Start(ctx, "provider/UpdateProject")
	defer span.End()

	err := p.repository.UpdateProject(ctx, project_id, ProviderProject2DBProject(project))
	return err
}

func (p Provider) DeleteProject(ctx context.Context, project_id int) error {
	ctx, span := tracer.Start(ctx, "provider/DeleteProject")
	defer span.End()

	err := p.repository.DeleteProject(ctx, project_id)
	return err
}
