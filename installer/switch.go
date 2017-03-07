package installer

import (
	"github.com/wx13/genesis"
)

// Switch runs a Doer depending on a value.
type Switch struct {
	Tasks map[genesis.Doer]bool
	Name  string
}

func NewSwitch(name string) *Switch {
	return &Switch{
		Tasks: map[genesis.Doer]bool{},
		Name:  name,
	}
}

func (sw Switch) Files() []string {
	files := []string{}
	for task, _ := range sw.Tasks {
		files = append(files, task.Files()...)
	}
	return files
}

func (sw Switch) ID() string {
	id := ""
	if sw.Name == "" {
		for task, _ := range sw.Tasks {
			id += task.ID()
		}
	} else {
		id = sw.Name
	}
	return id
}

func (sw *Switch) Case(condition bool, doer genesis.Doer) {
	sw.Tasks[doer] = condition
}

func (sw Switch) Status() (genesis.Status, error) {
	status := genesis.StatusPass
	for task, condition := range sw.Tasks {
		if !condition {
			continue
		}
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
	for task, condition := range sw.Tasks {
		if !condition {
			continue
		}
		changed, err := task.Do()
		if err != nil {
			return changed, err
		}
	}
	return true, nil
}

func (sw Switch) Undo() (bool, error) {
	for task, condition := range sw.Tasks {
		if !condition {
			continue
		}
		changed, err := task.Undo()
		if err != nil {
			return changed, err
		}
	}
	return true, nil
}
