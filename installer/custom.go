package installer

import (
	"github.com/wx13/genesis"
)

type Custom struct {
	Task genesis.Doer
	S    func() (genesis.Status, error)
	D    func() (bool, error)
	U    func() (bool, error)
	I    func() string
}

func NewCustom(task genesis.Doer) *Custom {
	custom := Custom{
		Task: task,
		S:    task.Status,
		D:    task.Do,
		U:    task.Undo,
		I:    task.ID,
	}
	return &custom
}

func (custom Custom) Status() (genesis.Status, error) {
	return custom.S()
}

func (custom Custom) Do() (bool, error) {
	return custom.D()
}

func (custom Custom) Undo() (bool, error) {
	return custom.U()
}

func (custom Custom) ID() string {
	return custom.I()
}
