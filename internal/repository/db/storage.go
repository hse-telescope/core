package db

import (
	"context"
	"database/sql"

	"github.com/hse-telescope/core/internal/repository/models"
	"github.com/hse-telescope/utils/db/psql"
	"github.com/jmoiron/sqlx"
)

type DB struct {
	db *sql.DB
}

func (s DB) New(param1 string, path string) (any, error) {
	panic("unimplemented")
}

func New(dbURL string, migrationsPath string) (DB, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return DB{}, err
	}
	err = db.Ping()
	if err != nil {
		return DB{}, err
	}
	psql.MigrateDB(db, migrationsPath, psql.PGDriver)
	return DB{
		db: db,
	}, nil

}
func (s DB) GetProjects(ctx context.Context) ([]models.Project, error) {
	q := `
		SELECT
			id,
			name
		FROM projects
	`
	rows, err := s.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	projects := make([]models.Project, 0)
	err = sqlx.StructScan(rows, &projects)
	if err != nil {
		return nil, err
	}
	return projects, nil
}

func (s DB) CreateProject(ctx context.Context, project models.Project) (models.Project, error) {
	q := `
		INSERT INTO projects (name) VALUES ($1) RETURNING id
	`
	var newID int
	err := s.db.QueryRowContext(ctx, q, project.Name).Scan(&newID)
	project.ID = newID
	return project, err
}

func (s DB) UpdateProject(ctx context.Context, project_id int, project models.Project) error {
	q := `
        UPDATE projects
        SET name = $1
        WHERE id = $2
    `

	_, err := s.db.ExecContext(ctx, q, project.Name, project_id)
	return err
}

func (s DB) DeleteProject(ctx context.Context, project_id int) error {
	q := `
        DELETE FROM projects
        WHERE id = $1
    `

	_, err := s.db.ExecContext(ctx, q, project_id)
	return err
}

func (s DB) CreateGraph(ctx context.Context, graph models.Graph) (models.Graph, error) {
	q := `
		INSERT INTO graphs (project_id, name) VALUES ($1, $2) RETURNING id
	`
	var newID int
	err := s.db.QueryRowContext(ctx, q, graph.ProjectID, graph.Name).Scan(&newID)
	graph.ID = newID
	return graph, err
}

func (s DB) DeleteGraph(ctx context.Context, graph_id int) error {
	q := `
		DELETE FROM graphs WHERE id = $1
	`
	_, err := s.db.ExecContext(ctx, q, graph_id)
	return err
}

func (s DB) UpdateGraph(ctx context.Context, graph_id int, graph models.Graph) error {
	q := `
		UPDATE projects
        SET project_id = $1, name = $2
        WHERE id = $3
	`
	_, err := s.db.ExecContext(ctx, q, graph.ProjectID, graph.Name, graph_id)
	return err
}

func (s DB) GetProjectGraphs(ctx context.Context, project_id int) ([]models.Graph, error) {
	q := `
		SELECT
			id,
			project_id,
			name
		FROM graphs WHERE project_id = $1
	`
	rows, err := s.db.QueryContext(ctx, q, project_id)
	if err != nil {
		return nil, err
	}
	graphs := make([]models.Graph, 0)
	err = sqlx.StructScan(rows, &graphs)
	if err != nil {
		return nil, err
	}
	return graphs, nil
}

func (s DB) GetService(ctx context.Context, service_id int) (models.Service, error) {
	q := `
		SELECT id, graph_id, name, description, X, Y FROM services WHERE id = $1
	`
	var service models.Service
	if err := s.db.QueryRowContext(ctx, q, service_id).Scan(&service); err != nil {
		if err == sql.ErrNoRows {
			return models.Service{}, err
		}
		return models.Service{}, err
	}
	return service, nil
}
func (s DB) GetGraphServices(ctx context.Context, graph_id int) ([]models.Service, error) {
	q := `
		SELECT
			id,
			graph_id,
			name,
			description,
			x,
			y
		FROM services WHERE graph_id = $1
	`
	rows, err := s.db.QueryContext(ctx, q, graph_id)
	if err != nil {
		return nil, err
	}
	services := make([]models.Service, 0)
	err = sqlx.StructScan(rows, &services)
	if err != nil {
		return nil, err
	}
	return services, nil
}

func (s DB) UpdateService(ctx context.Context, service_id int, service models.Service) error {
	q := `
		UPDATE services
        SET graph_id = $1, name = $2, description = $3, x = $4, y = $5
        WHERE id = $6
	`
	_, err := s.db.ExecContext(ctx, q, service.GraphID, service.Name, service.Description, service.X, service.Y, service_id)
	return err
}

func (s DB) DeleteService(ctx context.Context, service_id int) error {
	q := `
		DELETE FROM services WHERE id = $1
	`
	_, err := s.db.ExecContext(ctx, q, service_id)
	return err
}

func (s DB) CreateService(ctx context.Context, service models.Service) (models.Service, error) {
	q := `
		INSERT INTO services (graph_id, name, description, x, y) VALUES ($1, $2, $3, $4, $5) RETURNING id
	`
	var newID int
	err := s.db.QueryRowContext(ctx, q, service.GraphID, service.Name, service.Description, service.X, service.Y).Scan(&newID)
	service.ID = newID
	return service, err
}

func (s DB) GetRelation(ctx context.Context, relation_id int) (models.Relation, error) {
	q := `
		SELECT id, graph_id, name, description, from_service, to_service FROM relations WHERE id = $1
	`
	var relation models.Relation
	if err := s.db.QueryRowContext(ctx, q, relation_id).Scan(&relation); err != nil {
		if err == sql.ErrNoRows {
			return models.Relation{}, err
		}
		return models.Relation{}, err
	}
	return relation, nil
}

func (s DB) GetGraphRelations(ctx context.Context, graph_id int) ([]models.Relation, error) {
	q := `
		SELECT
			id,
			graph_id,
			name,
			description,
			from_service,
			to_service
		FROM relations WHERE graph_id = $1
	`
	rows, err := s.db.QueryContext(ctx, q, graph_id)
	if err != nil {
		return nil, err
	}
	relations := make([]models.Relation, 0)
	err = sqlx.StructScan(rows, &relations)
	if err != nil {
		return nil, err
	}
	return relations, nil
}

func (s DB) CreateRelation(ctx context.Context, relation models.Relation) (models.Relation, error) {
	q := `
		INSERT INTO relations (graph_id, name, description, from_service, to_service) VALUES ($1, $2, $3, $4, $5) RETURNING id
	`
	var newID int
	err := s.db.QueryRowContext(ctx, q, relation.GraphID, relation.Name, relation.Description, relation.FromService, relation.ToService).Scan(&newID)
	relation.ID = newID
	return relation, err
}

func (s DB) UpdateRelation(ctx context.Context, relation_id int, relation models.Relation) error {
	q := `
		UPDATE relations
		SET graph_id = $1, name = $2, description = $3, from_service = $4, to_service = $5
		WHERE id = $6
	`
	_, err := s.db.ExecContext(ctx, q, relation.GraphID, relation.Name, relation.Description, relation.FromService, relation.ToService, relation_id)
	return err
}

func (s DB) DeleteRelation(ctx context.Context, relation_id int) error {
	q := `
		DELETE FROM relations WHERE id = $1
	`
	_, err := s.db.ExecContext(ctx, q, relation_id)
	return err
}
