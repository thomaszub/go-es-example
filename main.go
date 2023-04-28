package main

import (
	"log"

	"github.com/gocql/gocql"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/thomaszub/go-es-example/api"
	"github.com/thomaszub/go-es-example/database"
	"github.com/thomaszub/go-es-example/domain"
)

func main() {
	cfg, err := LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	err = database.InitSchema(cfg.CassandraCluster, cfg.CassandraKeyspace)
	if err != nil {
		log.Fatal(err)
	}

	cluster := gocql.NewCluster(cfg.CassandraCluster...)
	cluster.Keyspace = cfg.CassandraKeyspace
	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	repo := database.InitRepository(session)
	service := domain.NewAccountService(&repo)
	controller := api.NewAccountController(&service)

	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	g := e.Group("/api/accounts")
	controller.RegisterOn(g)
	e.Logger.Fatal(e.Start(":8000"))
}
