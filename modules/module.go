package modules

import (
	"html/template"
	"os"
	"reflect"

	"github.com/axelrindle/proxyguy/config"
	"github.com/axelrindle/proxyguy/logger"
	"github.com/sirupsen/logrus"
)

type Exports struct {
	Host    string
	Port    string
	NoProxy string
}

type Module interface {
	GetName() string
	GetTemplate() string
	GetLogger() *logrus.Entry

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

type DefaultModule struct{}

func (t DefaultModule) GetName() string {
	return reflect.TypeOf(t).Name()
}

func (t DefaultModule) GetTemplate() string {
	return ""
}

func (t DefaultModule) GetLogger() *logrus.Entry {
	return logger.ForModule(t.GetName())
}

func (t DefaultModule) IsEnabled(cfg config.StructureModules) bool {
	return true
}

func (t DefaultModule) Preprocess(data *Exports) {
}

func (t DefaultModule) OnNoProxy() {
}
