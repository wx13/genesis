package installer

import (
	"github.com/wx13/genesis"
)

// Custom is a type of genesis.Doer.  It is a wrapper around another
// Doer which allows for custom Status/Do/Undo functions.
type Custom struct {
	Task genesis.Doer
	S    func() (genesis.Status, error)
	D    func() (bool, error)
	U    func() (bool, error)
	I    func() string
	F    func() []string
}

func NewCustom(task genesis.Doer) *Custom {
	custom := Custom{
		Task: task,
		S:    task.Status,
		D:    task.Do,
		U:    task.Undo,
		I:    task.ID,
		F:    task.Files,
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

func (custom Custom) Files() []string {
	return custom.F()
}
