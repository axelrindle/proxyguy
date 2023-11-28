package modules

import (
	"os"
	"path"

	"github.com/axelrindle/proxyguy/config"
	"github.com/magiconair/properties"
	"github.com/sirupsen/logrus"
)

const gradleProperties = ".gradle/gradle.properties"

type TemplateGradle struct {
	Logger *logrus.Logger
}

func (t TemplateGradle) GetName() string {
	return "Gradle"
}

func (t TemplateGradle) GetTemplate() string {
	return ""
}

func (t TemplateGradle) IsEnabled(cfg config.StructureModules) bool {
	return cfg.Gradle
}

func (t TemplateGradle) Preprocess(data *Exports) {
	// add props to global gradle.properties
	properties.LogPrintf = func(fmt string, args ...interface{}) {
		t.Logger.Debugf(fmt, args)
	}

	t.doWithProperties(func(p *properties.Properties) {
		p.Set("systemProp.http.proxyHost", data.Host)
		p.Set("systemProp.http.proxyPort", data.Port)
		p.Set("systemProp.https.proxyHost", data.Host)
		p.Set("systemProp.https.proxyPort", data.Port)
		p.Set("systemProp.http.nonProxyHosts", data.NoProxy)
	})
}

func (t TemplateGradle) OnNoProxy() {
	// remove props from global gradle.properties
	properties.LogPrintf = func(fmt string, args ...interface{}) {
		t.Logger.Debugf(fmt, args)
	}

	t.doWithProperties(func(p *properties.Properties) {
		p.Delete("systemProp.http.proxyHost")
		p.Delete("systemProp.http.proxyPort")
		p.Delete("systemProp.https.proxyHost")
		p.Delete("systemProp.https.proxyPort")
		p.Delete("systemProp.http.nonProxyHosts")
	})
}

func (t TemplateGradle) doWithProperties(handler func(p *properties.Properties)) {
	properties.LogPrintf = func(fmt string, args ...interface{}) {
		t.Logger.Debugf(fmt, args)
	}

	p, err := loadFile()
	if err != nil {
		t.Logger.Error(err)
		return
	}

	handler(p)

	err = saveFile(p)
	if err != nil {
		t.Logger.Error(err)
		return
	}

	t.Logger.Debug("Written global gradle.properties")
}

func loadFile() (*properties.Properties, error) {
	file, err := absoluteFile()
	if err != nil {
		return nil, err
	}

	p, err := properties.LoadFiles([]string{file}, properties.UTF8, true)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func saveFile(p *properties.Properties) error {
	file, err := absoluteFile()
	if err != nil {
		return err
	}
	tmpFile := file + ".tmp"

	handle, err := os.OpenFile(tmpFile, os.O_SYNC|os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer handle.Close()

	_, err = p.Write(handle, properties.UTF8)
	if err != nil {
		return err
	}

	if err = os.Rename(tmpFile, file); err != nil {
		return err
	}

	return nil
}

func absoluteFile() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return path.Join(home, gradleProperties), nil
}
