package modules

import (
	"html/template"
	"os"

	"github.com/axelrindle/proxyguy/config"
)

type Exports struct {
	Host    string
	Port    string
	NoProxy string
}

type Module interface {
	GetName() string
	GetTemplate() string

	IsEnabled(cfg config.StructureModules) bool

	Preprocess(data *Exports)
	OnNoProxy()
}

func Process(mdl Module, data Exports) bool {
	tmplStr := mdl.GetTemplate()
	if tmplStr == "" {
		return true
	}

	tmpl, err := template.New(mdl.GetName()).Parse(tmplStr)
	if err != nil {
		return false
	}

	tmpl.Execute(os.Stdout, data)
	return true
}
