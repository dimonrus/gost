// command file
package script

import (
	"fmt"
	"github.com/dimonrus/gocli"
	"gost/app/base"
	"os"
	"text/template"
)

func init() {
	base.App.GetScripts()["create"] = func(args gocli.ArgumentMap) {
		fileName, ok := args["file"]
		if !ok {
			base.App.GetLogger().Errorln("File name is not specified")
			return
		}
		name := fileName.GetString()
		if name == "" {
			base.App.GetLogger().Errorln("File name is not specified")
			return
		}
		folderPath := "app/io/script"
		err := os.MkdirAll(folderPath, os.ModePerm)
		if err != nil {
			panic(err)
		}
		filePath := fmt.Sprintf("%s/%s.go", folderPath, name)

		f, err := os.Create(filePath)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		scriptTemplate := getScriptTemplate()
		err = scriptTemplate.Execute(f, struct {
			Name string
		}{
			Name: name,
		})
		if err != nil {
			base.App.FatalError(err)
		}
		base.App.GetLogger().Infof("script file created: %s", name)
	}
}

func getScriptTemplate() *template.Template {
	var scriptTemplate = template.Must(template.New("").
		Parse(`// script file
package script

import (
	"github.com/dimonrus/gocli"
	"gost/app/base"
)

func init() {
	base.App.GetScripts()["{{ .Name }}"] = func(args gocli.Arguments) {
		//write yours code here
	}
}`))
	return scriptTemplate
}
