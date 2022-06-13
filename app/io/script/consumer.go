// script file
package script

import (
	"fmt"
	"github.com/dimonrus/gocli"
	"gost/app/base"
	"os"
	"text/template"
)

func init() {
	base.App.GetScripts()["consumer"] = func(args gocli.Arguments) {
		fileName, ok := args["file"]
		if !ok || fileName.GetString() == "" {
			base.App.GetLogger().Errorln("File name is not specified")
			return
		}
		queue, ok := args["queue"]
		if !ok || queue.GetString() == "" {
			base.App.GetLogger().Errorln("Queue name is not specified")
			return
		}
		server, ok := args["server"]
		if !ok || server.GetString() == "" {
			base.App.GetLogger().Errorln("Server name is not specified")
			return
		}
		folderPath := "app/io/consumer"
		err := os.MkdirAll(folderPath, os.ModePerm)
		if err != nil {
			base.App.FatalError(err)
		}
		filePath := fmt.Sprintf("%s/%s.go", folderPath, fileName.GetString())

		f, err := os.Create(filePath)
		if err != nil {
			base.App.FatalError(err)
		}
		defer f.Close()

		scriptTemplate := getConsumerScriptTemplate()

		err = scriptTemplate.Execute(f, struct {
			Name   string
			Queue  string
			Server string
		}{
			Name:   fileName.GetString(),
			Queue:  queue.GetString(),
			Server: server.GetString(),
		})

		if err != nil {
			base.App.FatalError(err)
		}

		base.App.GetLogger().Infof("consumer file created: %s", fileName.GetString())
	}
}

// Consumer file template
func getConsumerScriptTemplate() *template.Template {
	var consumerTemplate = template.Must(template.New("").
		Parse(`// {{ .Name }} consumer file
package consumer

import (
	"github.com/dimonrus/gorabbit"
	"github.com/rabbitmq/amqp091-go"
	"gost/app/base"
)

func init() {
	base.App.GetRabbit().GetRegistry()["{{ .Name }}"] = &gorabbit.Consumer{Queue: "{{ .Queue }}", Server: "{{ .Server }}", Count: 1, Callback: func(d amqp.Delivery) {

	}}
}`))
	return consumerTemplate
}
