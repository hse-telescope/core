package project

import "github.com/hse-telescope/core/internal/repository/models"

type Project struct {
	ID   int
	Name string
}

func ProviderProject2DBProject(project Project) models.Project {
	return models.Project{
		ID:   project.ID,
		Name: project.Name,
	}
}

func DBProject2ProviderProject(project models.Project) Project {
	return Project{
		ID:   project.ID,
		Name: project.Name,
	}
}
