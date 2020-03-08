package base

import (
	"errors"
	"github.com/dimonrus/gocli"
	"github.com/dimonrus/godb"
	"github.com/dimonrus/gohelp"
	"github.com/dimonrus/gorabbit"
	"github.com/dimonrus/goweb"
	"github.com/dimonrus/porterr"
	_ "github.com/lib/pq"
	"gost/app/config"
	"net"
	"net/http"
	"os"
	"path/filepath"
)

// The Application
var (
	App *Application
)

// Application
type Application struct {
	gocli.Application
	baseDb         *godb.DBO
	tcpConnections *goweb.Connections
	rabbit         *gorabbit.Application
	web            *goweb.Application
	migration      *godb.Migration
	scripts        map[string]func(app gocli.Arguments)
}

// Get logger
func (app *Application) GetLogger() gocli.Logger {
	return app.Application.GetLogger(gocli.LogLevelDebug)
}

// Get tcp connections
func (app *Application) GetTcpConnections() *goweb.Connections {
	if app.tcpConnections == nil {
		app.tcpConnections = goweb.NewConnections()
	}
	return app.tcpConnections
}

// Get config
func (app *Application) GetConfig() config.Config {
	cfg := app.Application.GetConfig().(*config.Config)
	return *cfg
}

// Get application type
func (app *Application) GetAppType() string {
	cfg := app.GetConfig()
	appType, ok := cfg.Arguments["app"]
	if ok != true {
		app.FatalError(errors.New("app type is not presents"))
	}
	return appType.GetString()
}

// Get DB Connection
func (app *Application) GetDB() *godb.DBO {
	if app.baseDb != nil {
		return app.baseDb
	}
	cfg := app.GetConfig()
	// init base db connection
	var baseErr error
	app.baseDb, baseErr = godb.DBO{
		Options: godb.Options{
			Debug:  true,
			Logger: App.GetLogger(),
		},
		Connection: &cfg.Db,
	}.Init()
	// if connection failed
	if baseErr != nil {
		App.GetLogger().Errorf("Connection (%s) failed", app.baseDb.Connection)
		App.FatalError(baseErr)
	}
	return app.baseDb
}

// Get rabbit app
func (app *Application) GetRabbit() *gorabbit.Application {
	if app.rabbit == nil {
		app.rabbit = gorabbit.NewApplication(App.GetConfig().Rabbit, app.Application)
		app.rabbit.SetRegistry(make(gorabbit.Registry))
	}
	return app.rabbit
}

// Get web app
func (app *Application) GetWeb() *goweb.Application {
	if app.web == nil {
		app.web = goweb.NewApplication(App.GetConfig().Web, app.Application, nil)
	}
	return app.web
}

// Get migration
func (app *Application) GetMigration() *godb.Migration {
	if app.migration == nil {
		registry := make(godb.MigrationRegistry, 0)
		app.migration = &godb.Migration{
			RegistryPath:  "gost/app/base",
			MigrationPath: "app/io/db/migrations",
			RegistryXPath: "base.App.GetMigration().Registry",
			DBO:           app.GetDB(),
			Registry:      registry,
			Config:        app.GetConfig().Db.ConnectionConfig,
		}
	}
	return app.migration
}

// Connection State Listener
func (app *Application) ConnStateEvent(conn net.Conn, event http.ConnState) {
	id := goweb.ConnectionIdentifier(conn.RemoteAddr().String())
	if event == http.StateActive {
		app.GetTcpConnections().Set(id, conn)
	} else if event == http.StateHijacked || event == http.StateClosed {
		app.GetTcpConnections().Unset(id)
	}
}

// Get Script callback
func (app *Application) GetScripts() map[string]func(app gocli.Arguments) {
	if app.scripts == nil {
		app.scripts = make(map[string]func(app gocli.Arguments))
	}
	return app.scripts
}

// Begin transaction
func (app *Application) StartTransaction() *godb.SqlTx {
	tx, err := app.GetDB().Begin()
	if err != nil {
		app.FatalError(err)
		return nil
	}
	return tx
}

// End transaction
func (app *Application) EndTransaction(q *godb.SqlTx, e porterr.IError) {
	var err error
	if e != nil {
		err = q.Rollback()
	} else {
		err = q.Commit()
	}
	if err != nil {
		app.FatalError(err)
	}
	return
}

// Get absolute path to application
func (app *Application) GetAbsolutePath(path string) string {
	rootPath, err := filepath.Abs("")
	if err != nil {
		panic(err)
	}
	if rootPath[len(rootPath)-1:] != "/" {
		rootPath = rootPath + "/"
	}
	// compatibility with test and docker
	return gohelp.BeforeString(rootPath, "gost") + "gost/" + path
}

// Init application
func init() {
	var cfg config.Config

	// Get ENV value
	environment := os.Getenv("ENV")
	if environment == "" {
		panic("ENV is not defined")
	}

	App = &Application{
		Application: gocli.NewApplication(environment, App.GetAbsolutePath("app/config/yaml"), &cfg),
	}
	App.ParseFlags(&cfg.Arguments)
}
