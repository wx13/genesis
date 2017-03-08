package installer

import (
	"github.com/wx13/genesis"
)

// Switch runs a Doer depending on a value.
type Switch struct {
	Dos   []genesis.Doer
	Donts []genesis.Doer
	Name  string
}

func NewSwitch(name string) *Switch {
	return &Switch{
		Name: name,
	}
}

func (sw Switch) Files() []string {
	files := []string{}
	for _, task := range append(sw.Dos, sw.Donts...) {
		files = append(files, task.Files()...)
	}
	return files
}

func (sw Switch) ID() string {
	id := ""
	if sw.Name == "" {
		for _, task := range append(sw.Dos, sw.Donts...) {
			id += task.ID()
		}
	} else {
		id = sw.Name
	}
	return id
}

func (sw *Switch) Case(condition bool, doer genesis.Doer) {
	if condition {
		sw.Dos = append(sw.Dos, doer)
	} else {
		sw.Donts = append(sw.Donts, doer)
	}
}

// Else -- if all existing tasks are false, then this is true.
// If any existing task is true, this is false.
func (sw *Switch) Else(doer genesis.Doer) {
	if len(sw.Dos) == 0 {
		sw.Dos = append(sw.Dos, doer)
	} else {
		sw.Donts = append(sw.Donts, doer)
	}
}

func (sw Switch) Status() (genesis.Status, error) {
	status := genesis.StatusPass
	for _, task := range sw.Dos {
		s, _ := task.Status()
		if s == genesis.StatusFail {
			status = s
		}
		if s == genesis.StatusUnknown && status == genesis.StatusPass {
			status = s
		}
	}
	return status, nil
}

func (sw Switch) Do() (bool, error) {
	for _, task := range sw.Dos {
		changed, err := task.Do()
		if err != nil {
			return changed, err
		}
	}
	return true, nil
}

func (sw Switch) Undo() (bool, error) {
	for _, task := range sw.Dos {
		changed, err := task.Undo()
		if err != nil {
			return changed, err
		}
	}
	return true, nil
}
