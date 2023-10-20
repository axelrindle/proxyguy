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

var TemplateMain = &Module{
	Name:     "Main",
	Template: tpl_main,

	IsEnabled: func(cfg config.StructureModules) bool {
		return true
	},
}
