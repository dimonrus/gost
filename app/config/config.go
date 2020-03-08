package config

import (
	"github.com/dimonrus/gocli"
	"github.com/dimonrus/godb"
	"github.com/dimonrus/gorabbit"
	"github.com/dimonrus/goweb"
)

type Config struct {
	Project   ProjectConfig
	Web       goweb.Config
	Db        godb.PostgresConnectionConfig
	Arguments gocli.Arguments
	Rabbit    gorabbit.Config
}

type ProjectConfig struct {
	Name        string
	Debug       bool
}
