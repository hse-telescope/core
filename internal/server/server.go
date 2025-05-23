package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"

	"github.com/hse-telescope/core/internal/config"
	"github.com/hse-telescope/core/internal/providers/graph"

	"github.com/hse-telescope/core/internal/providers/project"
	"github.com/hse-telescope/core/internal/providers/relation"
	"github.com/hse-telescope/core/internal/providers/service"
)

type ProviderProject interface {
	GetProjects(ctx context.Context) ([]project.Project, error)
	CreateProject(ctx context.Context, project project.Project) (project.Project, error)
	UpdateProject(ctx context.Context, project_id int, project project.Project) error
	DeleteProject(ctx context.Context, project_id int) error
}

type ProviderGraph interface {
	GetGraphProject(ctx context.Context, graph_id int) (int, error)
	CreateGraph(ctx context.Context, graph graph.Graph) (graph.Graph, error)
	DeleteGraph(ctx context.Context, graph_id int) error
	UpdateGraph(ctx context.Context, graph_id int, graph graph.Graph) error
	GetProjectGraphs(ctx context.Context, project_id int) ([]graph.Graph, error)
}

type ProviderService interface {
	GetServiceGraph(ctx context.Context, service_id int) (int, error)
	GetService(ctx context.Context, service_id int) (service.Service, error)
	GetGraphServices(ctx context.Context, graph_id int) ([]service.Service, error)
	CreateService(ctx context.Context, service service.Service) (service.Service, error)
	CreateServices(ctx context.Context, graph_id int, service []service.Service) ([]int, error)
	UpdateGraphServices(ctx context.Context, graph_id int, service []service.Service) error
	UpdateService(ctx context.Context, service_id int, service service.Service) error
	DeleteService(ctx context.Context, service_id int) error
}

type ProviderRelation interface {
	GetRelationGraph(ctx context.Context, relation_id int) (int, error)
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

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	s.server.Handler = c.Handler(s.setRouter())
	s.providerProject = provideProject
	s.providerGraph = provideGraph
	s.providerService = provideService
	s.providerRelation = providerRelation
	return s
}

func (s *Server) setRouter() *mux.Router {
	mux := mux.NewRouter()

	mux.Handle("/api/core/projects", s.AuthMiddleware(http.HandlerFunc(s.createProjectHandler))).Methods(http.MethodPost)
	mux.Handle("/api/core/projects", s.AuthMiddleware(http.HandlerFunc(s.getProjectsHanlder))).Methods(http.MethodGet)
	mux.Handle("/api/core/projects/{id}", s.AuthMiddleware(http.HandlerFunc(s.deleteProjectHandler))).Methods(http.MethodDelete)
	mux.Handle("/api/core/projects/{id}", s.AuthMiddleware(http.HandlerFunc(s.updateProjectHandler))).Methods(http.MethodPut)
	mux.Handle("/api/core/projects/{id}/graphs", s.AuthMiddleware(http.HandlerFunc(s.GetProjectGraphsHandler))).Methods(http.MethodGet)

	mux.Handle("/api/core/graphs", s.AuthMiddleware(http.HandlerFunc(s.createGraphHandler))).Methods(http.MethodPost)
	mux.Handle("/api/core/graphs/{id}", s.AuthMiddleware(http.HandlerFunc(s.updateGraphHandler))).Methods(http.MethodPut)
	mux.Handle("/api/core/graphs/{id}", s.AuthMiddleware(http.HandlerFunc(s.deleteGraphHandler))).Methods(http.MethodDelete)
	mux.Handle("/api/core/graphs/{id}/services", s.AuthMiddleware(http.HandlerFunc(s.updateGraphServicesHandler))).Methods(http.MethodPut)
	mux.Handle("/api/core/graphs/{id}/relations", s.AuthMiddleware(http.HandlerFunc(s.updateGraphRelationsHandler))).Methods(http.MethodPut)
	mux.Handle("/api/core/graphs/{id}/services", s.AuthMiddleware(http.HandlerFunc(s.getGraphServicesHandler))).Methods(http.MethodGet)
	mux.Handle("/api/core/graphs/{id}/relations", s.AuthMiddleware(http.HandlerFunc(s.getGraphRelationsHandler))).Methods(http.MethodGet)
	mux.Handle("/api/core/graphs/{id}/services", s.AuthMiddleware(http.HandlerFunc(s.createGraphServicesHandler))).Methods(http.MethodPost)
	mux.Handle("/api/core/graphs/{id}/relations", s.AuthMiddleware(http.HandlerFunc(s.createGraphRelationsHandler))).Methods(http.MethodPost)

	mux.Handle("/api/core/services", s.AuthMiddleware(http.HandlerFunc(s.createServiceHandler))).Methods(http.MethodPost)
	mux.Handle("/api/core/services/{id}", s.AuthMiddleware(http.HandlerFunc(s.updateServiceHandler))).Methods(http.MethodPut)
	mux.Handle("/api/core/services/{id}", s.AuthMiddleware(http.HandlerFunc(s.deleteServiceHandler))).Methods(http.MethodDelete)
	mux.Handle("/api/core/services/{id}", s.AuthMiddleware(http.HandlerFunc(s.getServiceHandler))).Methods(http.MethodGet)

	mux.Handle("/api/core/relations", s.AuthMiddleware(http.HandlerFunc(s.createRelationHandler))).Methods(http.MethodPost)
	mux.Handle("/api/core/relations/{id}", s.AuthMiddleware(http.HandlerFunc(s.updateRelationHandler))).Methods(http.MethodPut)
	mux.Handle("/api/core/relations/{id}", s.AuthMiddleware(http.HandlerFunc(s.deleteRelationHandler))).Methods(http.MethodDelete)
	mux.Handle("/api/core/relations/{id}", s.AuthMiddleware(http.HandlerFunc(s.getRelationHandler))).Methods(http.MethodGet)
	return mux
}

func (s *Server) Start() error {
	return s.server.ListenAndServe()
}
