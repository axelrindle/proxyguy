package modules

import (
	"strings"

	"github.com/axelrindle/proxyguy/config"
)

var tpl_gradle = strings.Replace(tpl_maven, "MAVEN_OPTS", "GRADLE_OPTS", 1)

var TemplateGradle = &Module{
	Name:     "Gradle",
	Template: tpl_gradle,

	IsEnabled: func(cfg config.StructureModules) bool {
		return cfg.Gradle
	},
	Preprocess: TemplateMaven.Preprocess,
}
