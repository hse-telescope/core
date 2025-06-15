package main

import (
	"context"
	"os"

	"github.com/hse-telescope/core/internal/config"
	"github.com/hse-telescope/core/internal/providers/graph"
	"github.com/hse-telescope/core/internal/providers/project"
	"github.com/hse-telescope/core/internal/providers/relation"
	"github.com/hse-telescope/core/internal/providers/service"
	"github.com/hse-telescope/core/internal/repository/db"
	"github.com/hse-telescope/core/internal/repository/facade"
	"github.com/hse-telescope/core/internal/server"
	"github.com/hse-telescope/logger"
	"github.com/hse-telescope/tracer"
)

func main() {
	configPath := os.Args[1]
	conf, err := config.Parse(configPath)
	if err != nil {
		panic(err)
	}
	storage, err := db.New(conf.DB.GetDBURL(), conf.DB.MigrationsPath)
	if err != nil {
		panic(err)
	}

	logger.SetupLogger(context.Background(), "core", conf.OTELCollectorURL, conf.Logger)
	tracer.SetupTracer(context.Background(), "core", conf.OTELCollectorURL)

	facade := facade.New(storage)

	ProjectProvide := project.New(facade)
	GraphProvider := graph.New(facade)
	ServiceProvide := service.New(facade)
	RelationProvide := relation.New(facade)

	s := server.New(conf, ProjectProvide, GraphProvider, ServiceProvide, RelationProvide)
	panic(s.Start())
}
