package database

import (
	"bytes"
	_ "embed"
	"log"
	"text/template"

	"github.com/gocql/gocql"
)

//go:embed keyspace.cql
var keyspaceCql string

//go:embed schema.cql
var schemaCql string

func Initialize(cluster []string, keyspace string) error {
	err := createKeyspace(cluster, keyspace)
	if err != nil {
		return err
	}
	cl := gocql.NewCluster(cluster...)
	cl.Keyspace = keyspace
	session, err := cl.CreateSession()
	if err != nil {
		return err
	}
	if err = session.Query(schemaCql).Exec(); err != nil {
		return err
	}
	defer session.Close()
	log.Println("Schema successfully initialized")
	return nil
}

func createKeyspace(cluster []string, keyspace string) error {
	cl := gocql.NewCluster(cluster...)
	session, err := cl.CreateSession()
	if err != nil {
		return err
	}
	defer session.Close()
	temp, err := parseTemplate(keyspace)
	if err != nil {
		return err
	}
	if err = session.Query(temp).Exec(); err != nil {
		return err
	}
	return nil
}

func parseTemplate(keyspace string) (string, error) {
	bytes := &bytes.Buffer{}
	temp, err := template.New("create-keyspace").Parse(keyspaceCql)
	if err != nil {
		return "", err
	}
	err = temp.Execute(bytes, map[string]string{"Keyspace": keyspace})
	if err != nil {
		return "", err
	}
	return bytes.String(), nil
}
