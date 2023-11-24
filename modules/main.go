package modules

import "github.com/axelrindle/proxyguy/config"

const tpl_main = `
export http_proxy="http://{{.Host}}:{{.Port}}"
export https_proxy="http://{{.Host}}:{{.Port}}"
export no_proxy="{{.NoProxy}}"
export HTTP_PROXY="http://{{.Host}}:{{.Port}}"
export HTTPS_PROXY="http://{{.Host}}:{{.Port}}"
export NO_PROXY="{{.NoProxy}}"
`

type TemplateMain struct {
}

func (t TemplateMain) GetName() string {
	return "Main"
}

func (t TemplateMain) GetTemplate() string {
	return tpl_main
}

func (t TemplateMain) IsEnabled(cfg config.StructureModules) bool {
	return true
}

func (t TemplateMain) Preprocess(data *Exports) {
}

func (t TemplateMain) OnNoProxy() {
}
