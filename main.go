// Package main Gost Service.
//
// Entry point for the application
//
// Terms Of Service:
//
// All rights reserved by @GostCompany
//
//     Schemes: https
//     Host: gost.com
//     BasePath: /gost/api
//     Version: 1.0.0
//
//     Consumes:
//     - application/json
//     - multipart/form-data
//
//     Produces:
//     - application/json
//     - binary
//
// swagger:meta
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

// Entry points for application
func main() {
	var e porterr.IError

	switch base.App.GetAppType() {
	case ApplicationTypeWeb:
		base.App.GetWeb().Listen(web.GetRoutes())
	case ApplicationTypeScript:
		arguments := base.App.GetConfig().Arguments
		name, ok := arguments["name"]
		if !ok {
			base.App.FatalError(errors.New("no script name specified"))
		}
		scriptName := name.GetString()
		if callback, ok := base.App.GetScripts()[scriptName]; ok {
			base.App.SuccessMessage(fmt.Sprintf("Start script: %s", scriptName))
			callback(arguments)
		} else {
			base.App.FailMessage(fmt.Sprintf("Script file %s.go not found", scriptName))
		}
	case ApplicationTypeConsumer:
		c := gorabbit.CommandConsumer + " " + gorabbit.CommandStart + " " + gorabbit.CommandKeyWordAll
		command := gocli.ParseCommand([]byte(c))
		app := base.App.GetRabbit()
		app.SuccessMessage("Starting AMQP Application...", command)
		app.ConsumerCommander(command)
		e = app.Start("3333", app.ConsumerCommander)
		if e != nil {
			app.FailMessage(e.Error(), command)
		}
	default:
		e = porterr.New(porterr.PortErrorParam, "app type is undefined")
	}
	if e != nil {
		base.App.FatalError(e)
	}
}
