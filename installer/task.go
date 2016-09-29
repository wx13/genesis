package installer

import (
	"github.com/wx13/genesis"
)

// Task is the most fundamental Doer. It consists of just a single module.
// All other Doers contain tasks at their deepest levels (i.e. only a Task
// can contain a module directly.
type Task struct {
	genesis.Module
}

func (task Task) Status() (genesis.Status, error) {
	id := task.Module.ID()
	if SkipID(id) != "do" {
		return genesis.StatusUnknown, nil
	}
	desc := task.Describe()
	PrintHeader(id, desc)
	status, msg, err := task.Module.Status()
	if err != nil || status == genesis.StatusFail {
		ReportFail(msg, err)
		return status, err
	}
	if status == genesis.StatusUnknown {
		ReportUnknown(msg, err)
		return status, nil
	}
	if status == genesis.StatusPass {
		ReportPass(msg, err)
		return status, nil
	}
	return genesis.StatusUnknown, nil
}

func (task Task) Do() (bool, error) {

	id := task.ID()
	if SkipID(id) != "do" {
		return false, nil
	}

	desc := task.Describe()
	PrintHeader(id, desc)

	// If status is passing, then we don't have
	// to do anything.
	status, msg, err := task.Module.Status()
	if status == genesis.StatusPass {
		ReportPass(msg, err)
		return false, nil
	}

	// Otherwise, run the installer.
	msg, err = task.Install()
	if err != nil {
		ReportFail(msg, err)
		return false, err
	}

	// Check results.
	status, msg2, err := task.Module.Status()
	if status == genesis.StatusFail {
		ReportFail(msg2, err)
		return false, err
	}
	ReportDone(msg, err)
	return true, err
}

func (task Task) Undo() (bool, error) {
	id := task.ID()
	if SkipID(id) != "do" {
		return false, nil
	}
	desc := task.Describe()
	PrintHeader(id, desc)
	status, msg, err := task.Module.Status()
	if err != nil {
		ReportFail(msg, err)
		return false, err
	}
	if status == genesis.StatusFail {
		ReportPass(msg, err)
		return false, nil
	}
	msg, err = task.Remove()
	if err != nil {
		ReportFail(msg, err)
		return false, err
	}
	ReportPass(msg, err)
	return true, nil
}
