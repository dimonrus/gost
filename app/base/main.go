package base

import (
	"errors"
	"github.com/dimonrus/gocli"
	"github.com/dimonrus/godb/v2"
	"github.com/dimonrus/gohelp"
	"github.com/dimonrus/gorabbit"
	"github.com/dimonrus/goweb"
	"github.com/dimonrus/migrate"
	"github.com/dimonrus/porterr"
	_ "github.com/lib/pq"
	"gost/app/config"
	"net"
	"net/http"
	"os"
	"path/filepath"
)

var (
	// App Application object
	App *Application
	// StaticConfigPath relative config path
	StaticConfigPath = "app/config/yaml"
)

// Application main application struct
type Application struct {
	gocli.Application
	baseDb         *godb.DBO
	tcpConnections *goweb.Connections
	rabbit         *gorabbit.Application
	web            *goweb.Application
	migration      *migrate.Migration
	scripts        map[string]func(app gocli.ArgumentMap)
}

// GetTcpConnections Get tcp connections
func (app *Application) GetTcpConnections() *goweb.Connections {
	if app.tcpConnections == nil {
		app.tcpConnections = goweb.NewConnections()
	}
	return app.tcpConnections
}

// GetConfig Get app config
func (app *Application) GetConfig() config.Config {
	cfg := app.Application.GetConfig().(*config.Config)
	return *cfg
}

// GetAppType Get application type
func (app *Application) GetAppType() string {
	cfg := app.GetConfig()
	appType, ok := cfg.Arguments["app"]
	if !ok {
		app.FatalError(errors.New("app type is not presents"))
	}
	return appType.GetString()
}

// GetDB Get DB Connection
func (app *Application) GetDB() *godb.DBO {
	if app.baseDb != nil {
		return app.baseDb
	}
	cfg := app.GetConfig()
	// init base db connection
	var baseErr error
	app.baseDb, baseErr = godb.DBO{
		Options: godb.Options{
			Debug:          cfg.Project.Debug,
			Logger:         App.GetLogger(),
			QueryProcessor: godb.PreparePositionalArgsQuery,
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

// GetRabbit Get rabbit app
func (app *Application) GetRabbit() *gorabbit.Application {
	if app.rabbit == nil {
		app.rabbit = gorabbit.NewApplication(App.GetConfig().Rabbit, app.Application)
	}
	return app.rabbit
}

// GetWeb Get web app
func (app *Application) GetWeb() *goweb.Application {
	if app.web == nil {
		cfg := App.GetConfig().Web
		cfg.Debug = app.GetConfig().Project.Debug
		app.web = goweb.NewApplication(cfg, app.Application, nil)
	}
	return app.web
}

// GetMigration Get migration
func (app *Application) GetMigration() *migrate.Migration {
	if app.migration == nil {
		app.migration = &migrate.Migration{
			RegistryPath:  "gost/app/base",
			MigrationPath: "app/io/db/migrations",
			RegistryXPath: "base.App.GetMigration().Registry",
			DBO:           app.baseDb,
			Registry:      make(migrate.MigrationRegistry),
		}
	}
	return app.migration
}

// ConnStateEvent Connection State Listener
func (app *Application) ConnStateEvent(conn net.Conn, event http.ConnState) {
	id := goweb.ConnectionIdentifier(conn.RemoteAddr().String())
	if event == http.StateActive {
		app.GetTcpConnections().Set(id, conn)
	} else if event == http.StateHijacked || event == http.StateClosed {
		app.GetTcpConnections().Unset(id)
	}
}

// GetScripts Get Script callback
func (app *Application) GetScripts() map[string]func(app gocli.ArgumentMap) {
	if app.scripts == nil {
		app.scripts = make(map[string]func(app gocli.ArgumentMap))
	}
	return app.scripts
}

// StartTransaction Begin transaction
func (app *Application) StartTransaction() *godb.SqlTx {
	tx, err := app.GetDB().Begin()
	if err != nil {
		app.FatalError(err)
		return nil
	}
	return tx
}

// EndTransaction End transaction
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
}

// GetAbsolutePath Get absolute path to application
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

// GetENV return main env value
func (app *Application) GetENV() string {
	environment := os.Getenv("ENV")
	if environment == "" {
		environment = "local"
	}
	return environment
}

func init() {
	// Init app
	var cfg config.Config
	App = &Application{
		Application: gocli.NewApplication(App.GetENV(), App.GetAbsolutePath(StaticConfigPath), &cfg),
	}
	App.SetLogger(gocli.NewLogger(cfg.Logger))
	App.ParseFlags(cfg.Arguments)
}
