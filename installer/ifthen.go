package installer

import (
	"github.com/wx13/genesis"
)

type IfThen struct {
	If   genesis.Doer
	Then genesis.Doer
}

func (ifthen IfThen) ID() string {
	return ifthen.If.ID() + ifthen.Then.ID()
}

func (ifthen IfThen) Status() (genesis.Status, error) {
	skip := SkipID(ifthen.ID())
	if skip == "skip" {
		return genesis.StatusUnknown, nil
	}
	if skip == "do" {
		TempEmptyDoTags()
		defer RestoreDoTags()
	}
	status := genesis.StatusPass
	for _, task := range []genesis.Doer{ifthen.If, ifthen.Then} {
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

func (ifthen IfThen) Do() (bool, error) {
	skip := SkipID(ifthen.ID())
	if skip == "skip" {
		return false, nil
	}
	if skip == "do" {
		TempEmptyDoTags()
		defer RestoreDoTags()
	}
	changed, err := ifthen.If.Do()
	if err != nil {
		return changed, err
	}
	if changed {
		return ifthen.Then.Do()
	}
	return false, nil
}

func (ifthen IfThen) Undo() (bool, error) {
	skip := SkipID(ifthen.ID())
	if skip == "skip" {
		return false, nil
	}
	if skip == "do" {
		TempEmptyDoTags()
		defer RestoreDoTags()
	}
	changed, err := ifthen.If.Undo()
	if err != nil {
		return changed, err
	}
	if changed {
		return ifthen.Then.Undo()
	}
	return false, nil
}
