package config

import (
	"github.com/dimonrus/gocli"
	"github.com/dimonrus/godb/v2"
	"github.com/dimonrus/gorabbit"
	"github.com/dimonrus/goweb"
)

// Config application config
type Config struct {
	// Business main configuration. Extendable
	Project ProjectConfig
	// Logger config
	Logger gocli.LoggerConfig
	// Web api server configuration
	Web goweb.Config
	// Database credentials and config
	Db godb.PostgresConnectionConfig
	// Cli application allowed arguments
	Arguments gocli.Arguments
	// Rabbit MQ configuration
	Rabbit gorabbit.Config
}

// ProjectConfig Business config
type ProjectConfig struct {
	// Project name
	Name string
	// Debug mode
	Debug bool
}
