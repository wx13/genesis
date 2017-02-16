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
	Dest string
	Src  string
	Vars interface{}
}

func (tmpl Template) src() string {
	match, _ := regexp.MatchString("^[.]?/", tmpl.Src)
	if match {
		return tmpl.Src
	}
	return filepath.Join(genesis.Tmpdir, tmpl.Src)
}

func (tmpl Template) ID() string {
	return fmt.Sprintf("Template: %s => %s", tmpl.Src, tmpl.Dest)
}

func (tmpl Template) Files() []string {
	return []string{tmpl.src()}
}

func (tmpl Template) Remove() (string, error) {
	err := genesis.Store.RestoreFile(tmpl.Dest, "")
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
	err = genesis.Store.SaveFile(tmpl.Dest, "")
	if err != nil {
		return "Could not save snapshot to file store.", err
	}
	file, err := os.Create(tmpl.Dest)
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

	b, _ := ioutil.ReadFile(tmpl.Dest)
	fStr := string(b)
	if fStr != tmplStr {
		return genesis.StatusFail, "Template and destination differ", errors.New("Template and destination differ.")
	}
	return genesis.StatusPass, "Template file installed.", nil
}
