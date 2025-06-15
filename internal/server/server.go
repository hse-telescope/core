package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/hse-telescope/core/internal/config"
	"github.com/hse-telescope/core/internal/providers/graph"
	"github.com/hse-telescope/core/internal/providers/project"
	"github.com/hse-telescope/core/internal/providers/relation"
	"github.com/hse-telescope/core/internal/providers/service"
	"github.com/hse-telescope/logger"
	"github.com/hse-telescope/tracer"
)

type ProviderProject interface {
	GetProjects(ctx context.Context) ([]project.Project, error)
	CreateProject(ctx context.Context, project project.Project) (project.Project, error)
	UpdateProject(ctx context.Context, project_id int, project project.Project) error
	DeleteProject(ctx context.Context, project_id int) error
}

type ProviderGraph interface {
	CreateGraph(ctx context.Context, graph graph.Graph) (graph.Graph, error)
	DeleteGraph(ctx context.Context, graph_id int) error
	UpdateGraph(ctx context.Context, graph_id int, graph graph.Graph) error
	GetProjectGraphs(ctx context.Context, project_id int) ([]graph.Graph, error)
}

type ProviderService interface {
	GetService(ctx context.Context, service_id int) (service.Service, error)
	GetGraphServices(ctx context.Context, graph_id int) ([]service.Service, error)
	CreateService(ctx context.Context, service service.Service) (service.Service, error)
	CreateServices(ctx context.Context, graph_id int, service []service.Service) ([]int, error)
	UpdateGraphServices(ctx context.Context, graph_id int, service []service.Service) error
	UpdateService(ctx context.Context, service_id int, service service.Service) error
	DeleteService(ctx context.Context, service_id int) error
}

type ProviderRelation interface {
	GetRelation(ctx context.Context, relation_id int) (relation.Relation, error)
	GetGraphRelations(ctx context.Context, graph_id int) ([]relation.Relation, error)
	CreateRelation(ctx context.Context, relation relation.Relation) (relation.Relation, error)
	CreateRelations(ctx context.Context, graph_id int, relations []relation.Relation) error
	UpdateGraphRelations(ctx context.Context, graph_id int, relation []relation.Relation) error
	UpdateRelation(ctx context.Context, relation_id int, relation relation.Relation) error
	DeleteRelation(ctx context.Context, relation_id int) error
}

type Server struct {
	server           http.Server
	providerProject  ProviderProject
	providerGraph    ProviderGraph
	providerService  ProviderService
	providerRelation ProviderRelation
}

func New(conf config.Config, provideProject ProviderProject, provideGraph ProviderGraph, provideService ProviderService, providerRelation ProviderRelation) *Server {
	s := new(Server)
	s.server.Addr = fmt.Sprintf(":%d", conf.Port)
	s.server.Handler = s.setRouter()
	s.providerProject = provideProject
	s.providerGraph = provideGraph
	s.providerService = provideService
	s.providerRelation = providerRelation
	return s
}

func (s *Server) setRouter() *mux.Router {
	mux := mux.NewRouter()

	mux.Use(logger.AddLoggingMiddleware, tracer.AddTracingMiddleware)

	mux.HandleFunc("/projects", s.createProjectHandler).Methods(http.MethodPost)
	mux.HandleFunc("/projects", s.getProjectsHanlder).Methods(http.MethodGet)
	mux.HandleFunc("/projects/{id}", s.deleteProjectHandler).Methods(http.MethodDelete)
	mux.HandleFunc("/projects/{id}", s.updateProjectHandler).Methods(http.MethodPut)
	mux.HandleFunc("/projects/{id}/graphs", s.GetProjectGraphsHandler).Methods(http.MethodGet)

	mux.HandleFunc("/graphs", s.createGraphHandler).Methods(http.MethodPost)
	mux.HandleFunc("/graphs/{id}", s.updateGraphHandler).Methods(http.MethodPut)
	mux.HandleFunc("/graphs/{id}", s.deleteGraphHandler).Methods(http.MethodDelete)
	mux.HandleFunc("/graphs/{id}/services", s.updateGraphServicesHandler).Methods(http.MethodPut)
	mux.HandleFunc("/graphs/{id}/relations", s.updateGraphRelationsHandler).Methods(http.MethodPut)
	mux.HandleFunc("/graphs/{id}/services", s.getGraphServicesHandler).Methods(http.MethodGet)
	mux.HandleFunc("/graphs/{id}/relations", s.getGraphRelationsHandler).Methods(http.MethodGet)
	mux.HandleFunc("/graphs/{id}/services", s.createGraphServicesHandler).Methods(http.MethodPost)
	mux.HandleFunc("/graphs/{id}/relations", s.createGraphRelationsHandler).Methods(http.MethodPost)

	mux.HandleFunc("/services", s.createServiceHandler).Methods(http.MethodPost)
	mux.HandleFunc("/services/{id}", s.updateServiceHandler).Methods(http.MethodPut)
	mux.HandleFunc("/services/{id}", s.deleteServiceHandler).Methods(http.MethodDelete)
	mux.HandleFunc("/services/{id}", s.getServiceHandler).Methods(http.MethodGet)

	mux.HandleFunc("/relations", s.createRelationHandler).Methods(http.MethodPost)
	mux.HandleFunc("/relations/{id}", s.updateRelationHandler).Methods(http.MethodPut)
	mux.HandleFunc("/relations/{id}", s.deleteRelationHandler).Methods(http.MethodDelete)
	mux.HandleFunc("/relations/{id}", s.getRelationHandler).Methods(http.MethodGet)

	return mux
}

func (s *Server) Start() error {
	return s.server.ListenAndServe()
}
