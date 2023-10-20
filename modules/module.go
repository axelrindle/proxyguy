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

type Module struct {
	Name     string
	Template string

	IsEnabled  func(cfg config.StructureModules) bool
	Preprocess func(data Exports) Exports
}

func Process(mdl Module, data Exports) bool {
	tmpl, err := template.New("env").Parse(mdl.Template)
	if err != nil {
		return false
	}

	tmpl.Execute(os.Stdout, data)
	return true
}
