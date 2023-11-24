package modules

import (
	"strings"

	"github.com/axelrindle/proxyguy/config"
)

const tplMaven = `
export MAVEN_OPTS="-Dhttp.proxyHost={{.Host}} -Dhttp.proxyPort={{.Port}} -Dhttps.proxyHost={{.Host}} -Dhttps.proxyPort={{.Port}} -Dhttp.nonProxyHosts={{.NoProxy}}"
`

type TemplateMaven struct {
}

func (t TemplateMaven) GetName() string {
	return "Maven"
}

func (t TemplateMaven) GetTemplate() string {
	return tplMaven
}

func (t TemplateMaven) IsEnabled(cfg config.StructureModules) bool {
	return cfg.Maven
}

func (t TemplateMaven) Preprocess(data *Exports) {
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
}

func (t TemplateMaven) OnNoProxy() {
}
