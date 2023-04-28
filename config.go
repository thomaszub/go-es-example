package main

import (
	"errors"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	CassandraCluster  []string
	CassandraKeyspace string
}

func LoadConfig() (Config, error) {
	cfg := Config{}
	err := godotenv.Load()
	if err != nil {
		return cfg, err
	}
	cassandraCluster := strings.TrimSpace(os.Getenv("CASSANDRA_CLUSTER"))
	if cassandraCluster == "" {
		return cfg, errors.New("CASSANDRA_CLUSTER is not set")
	}
	cfg.CassandraCluster = strings.Split(cassandraCluster, ",")

	cassandraKeyspace := strings.TrimSpace(os.Getenv("CASSANDRA_KEYSPACE"))
	if cassandraKeyspace == "" {
		return cfg, errors.New("CASSANDRA_KEYSPACE is not set")
	}
	cfg.CassandraKeyspace = cassandraKeyspace
	return cfg, nil
}
