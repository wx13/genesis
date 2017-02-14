package modules

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"text/template"

	"github.com/wx13/genesis"
)

type Template struct {
	DestFile     string
	TemplateFile string
	Vars         interface{}
}

func (tmpl Template) src() string {
	match, _ := regexp.MatchString("^[.]?/", tmpl.TemplateFile)
	if match {
		return tmpl.TemplateFile
	}
	return filepath.Join(genesis.Tmpdir, tmpl.TemplateFile)
}

func (tmpl Template) ID() string {
	return fmt.Sprintf("Template: %s => %s", tmpl.TemplateFile, tmpl.DestFile)
}

func (tmpl Template) Files() []string {
	return []string{tmpl.src()}
}

func (tmpl Template) Remove() (string, error) {
	err := genesis.Store.RestoreFile(tmpl.DestFile, "")
	if err == nil {
		return "Successfully restored template file.", nil
	}
	return "Failed to restore template file.", err
}

func (tmpl Template) Install() (string, error) {

	t, err := template.ParseFiles(tmpl.src())
	if err != nil {
		return "Could not read template file.", err
	}
	err = genesis.Store.SaveFile(tmpl.DestFile, "")
	if err != nil {
		return "Could not save snapshot to file store.", err
	}
	file, err := os.Create(tmpl.DestFile)
	if err != nil {
		return "Could not create destination file.", err
	}
	err = t.Execute(file, tmpl.Vars)
	if err != nil {
		return "Failed to execute template.", err
	}
	return "Successfully ran template file.", nil

}

func (tmpl Template) Status() (genesis.Status, string, error) {

	t, err := template.ParseFiles(tmpl.src())
	if err != nil {
		return genesis.StatusFail, "Could not read template file.", err
	}

	buf := new(bytes.Buffer)
	t.Execute(buf, tmpl.Vars)
	tmplStr := buf.String()

	b, _ := ioutil.ReadFile(tmpl.DestFile)
	fStr := string(b)
	if fStr != tmplStr {
		return genesis.StatusFail, "Template and destination differ", errors.New("Template and destination differ.")
	}
	return genesis.StatusPass, "Template file installed.", nil
}
