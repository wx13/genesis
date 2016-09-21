package installer

import (
	"github.com/wx13/genesis"
)

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
	status, msg, err := task.Module.Status()
	if status == genesis.StatusPass {
		ReportPass(msg, err)
		return false, nil
	}
	msg, err = task.Install()
	if err != nil {
		ReportFail(msg, err)
		return false, err
	}
	status, msg2, err := task.Module.Status()
	if status == genesis.StatusPass {
		ReportDone(msg, err)
		return true, err
	}
	ReportFail(msg2, err)
	return false, err
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
