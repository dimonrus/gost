package config

import (
	"github.com/dimonrus/gocli"
	"github.com/dimonrus/godb/v2"
	"github.com/dimonrus/gorabbit"
	"github.com/dimonrus/goweb"
)

type Config struct {
	Project   ProjectConfig
	Logger    gocli.LoggerConfig
	Web       goweb.Config
	Db        godb.PostgresConnectionConfig
	Arguments gocli.Arguments
	Rabbit    gorabbit.Config
}

type ProjectConfig struct {
	Name  string
	Debug bool
}
