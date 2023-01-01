package main

import (
	"errors"
	"fmt"
	"github.com/dimonrus/gocli"
	"github.com/dimonrus/gorabbit"
	"github.com/dimonrus/porterr"
	"gost/app/base"
	_ "gost/app/io/consumer"
	_ "gost/app/io/script"
	"gost/app/io/web"
)

const (
	ApplicationTypeWeb      = "web"
	ApplicationTypeScript   = "script"
	ApplicationTypeConsumer = "consumer"
)

// RunWeb start web application
func RunWeb() porterr.IError {
	base.App.GetWeb().Listen(web.GetRoutes())
	return nil
}

// RunScript start script application
func RunScript() porterr.IError {
	arguments := base.App.GetConfig().Arguments
	name, ok := arguments["name"]
	if !ok {
		base.App.FatalError(errors.New("no script name specified"))
	}
	var e porterr.IError
	scriptName := name.GetString()
	if callback, ok := base.App.GetScripts()[scriptName]; ok {
		base.App.SuccessMessage(fmt.Sprintf("Start script: %s", scriptName))
		callback(arguments)
	} else {
		e = porterr.NewF(porterr.PortErrorScript, "Script file %s.go not found", scriptName)
		base.App.FailMessage(e.Error())
	}
	return e
}

// RunConsumer start consumer application
func RunConsumer() porterr.IError {
	c := []byte(gorabbit.CommandConsumer + " " + gorabbit.CommandStart + " " + gorabbit.CommandKeyWordAll)
	command := gocli.ParseCommand(c)
	app := base.App.GetRabbit()
	app.SuccessMessage("Starting AMQP Application...", command)
	app.ConsumerCommander(command)
	e := app.Start(":3333", app.ConsumerCommander)
	if e != nil {
		app.FailMessage(e.Error(), command)
	}
	return e
}

// Entry points for application
func main() {
	base.App.SuccessMessage("Application environment ENV=" + base.App.GetENV())
	var e porterr.IError
	switch base.App.GetAppType() {
	case ApplicationTypeWeb:
		e = RunWeb()
	case ApplicationTypeScript:
		e = RunScript()
	case ApplicationTypeConsumer:
		e = RunConsumer()
	default:
		e = porterr.New(porterr.PortErrorParam, "app type is undefined")
	}
	if e != nil {
		base.App.FatalError(e)
	}
}
