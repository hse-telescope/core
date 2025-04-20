package main

import (
	"os"

	"github.com/hse-telescope/core/internal/config"
	"github.com/hse-telescope/core/internal/providers/graph"
	"github.com/hse-telescope/core/internal/providers/project"
	"github.com/hse-telescope/core/internal/providers/relation"
	"github.com/hse-telescope/core/internal/providers/service"
	"github.com/hse-telescope/core/internal/repository/db"
	"github.com/hse-telescope/core/internal/repository/facade"
	"github.com/hse-telescope/core/internal/server"
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

	facade := facade.New(storage)

	ProjectProvide := project.New(facade)
	GraphProvider := graph.New(facade)
	ServiceProvide := service.New(facade)
	RelationProvide := relation.New(facade)

	s := server.New(conf, ProjectProvide, GraphProvider, ServiceProvide, RelationProvide)
	panic(s.Start())
}
