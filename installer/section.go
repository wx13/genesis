package installer

import (
	"github.com/wx13/genesis"
)

type Section struct {
	Tasks []genesis.Doer
	Name  string
}

func NewGroup() *Section {
	return &Section{
		Tasks: []genesis.Doer{},
	}
}

func NewSection(name string) *Section {
	return &Section{
		Tasks: []genesis.Doer{},
		Name:  name,
	}
}

func (section Section) ID() string {
	id := ""
	if section.Name == "" {
		for _, task := range section.Tasks {
			id += task.ID()
		}
	} else {
		id = section.Name
	}
	return id
}

func (section *Section) AddTask(module genesis.Module) {
	section.Tasks = append(section.Tasks, Task{module})
}

func (section *Section) Add(doer genesis.Doer) {
	section.Tasks = append(section.Tasks, doer)
}

func (section Section) Status() (genesis.Status, error) {
	skip := SkipID(section.ID())
	if skip == "skip" {
		return genesis.StatusUnknown, nil
	}
	if skip == "do" {
		TempEmptyDoTags()
		defer RestoreDoTags()
	}
	PrintSectionHeader(section.Name)
	defer PrintSectionFooter(section.Name)
	status := genesis.StatusPass
	for _, task := range section.Tasks {
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

func (section Section) Do() (bool, error) {
	skip := SkipID(section.ID())
	if skip == "skip" {
		return false, nil
	}
	if skip == "do" {
		TempEmptyDoTags()
		defer RestoreDoTags()
	}
	PrintSectionHeader(section.Name)
	defer PrintSectionFooter(section.Name)
	for _, task := range section.Tasks {
		changed, err := task.Do()
		if err != nil {
			return changed, err
		}
	}
	return true, nil
}

func (section Section) Undo() (bool, error) {
	skip := SkipID(section.ID())
	if skip == "skip" {
		return false, nil
	}
	if skip == "do" {
		TempEmptyDoTags()
		defer RestoreDoTags()
	}
	PrintSectionHeader(section.Name)
	defer PrintSectionFooter(section.Name)
	for k := len(section.Tasks) - 1; k >= 0; k-- {
		task := section.Tasks[k]
		changed, err := task.Do()
		if err != nil {
			return changed, err
		}
	}
	return true, nil
}
