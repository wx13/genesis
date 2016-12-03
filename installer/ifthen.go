package installer

import (
	"github.com/wx13/genesis"
)

// IfThen is a type of genesis.Doer. It runs doer B if and only if doer A runs.
type IfThen struct {
	If   genesis.Doer
	Then genesis.Doer
}

// ID for IfThen is just the combination of the two Doer IDs.
func (ifthen IfThen) ID() string {
	return ifthen.If.ID() + ifthen.Then.ID()
}

func (ifthen IfThen) Files() []string {
	return append(ifthen.If.Files(), ifthen.Then.Files()...)
}

func (ifthen IfThen) Status() (genesis.Status, error) {
	skip := SkipID(ifthen.ID())
	if skip == "skip" {
		return genesis.StatusUnknown, nil
	}
	if skip == "do" {
		doTags := EmptyDoTags()
		defer RestoreDoTags(doTags)
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
		doTags := EmptyDoTags()
		defer RestoreDoTags(doTags)
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
		doTags := EmptyDoTags()
		defer RestoreDoTags(doTags)
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
