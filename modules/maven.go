package modules

import (
	"strings"

	"github.com/axelrindle/proxyguy/config"
)

const tpl_maven = `
export MAVEN_OPTS="-Dhttp.proxyHost={{.Host}} -Dhttp.proxyPort={{.Port}} -Dhttps.proxyHost={{.Host}} -Dhttps.proxyPort={{.Port}} -Dhttp.nonProxyHosts={{.NoProxy}}"
`

var TemplateMaven = &Module{
	Name:     "Maven",
	Template: tpl_maven,

	IsEnabled: func(cfg config.StructureModules) bool {
		return cfg.Maven
	},
	Preprocess: func(data Exports) Exports {
		// follow the specification at https://docs.oracle.com/javase/8/docs/api/java/net/doc-files/net-properties.html
		// transform separator "," to "|"
		// use "*" as wildcard

		split := strings.Split(data.NoProxy, ",")

		for i, s := range split {
			if strings.HasPrefix(s, ".") {
				split[i] = "*" + s
			}
		}

		data.NoProxy = strings.Join(split, "|")

		return data
	},
}
