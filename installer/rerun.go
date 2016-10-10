package installer

import (
	"github.com/peterh/liner"
)

func Rerun(dir string) (string, error) {

	lnr := liner.NewLiner()
	defer lnr.Close()
	lnr.SetCtrlCAborts(true)

	cmds := GetHistory(dir)
	for i := len(cmds) - 1; i >= 0; i-- {
		lnr.AppendHistory(cmds[i])
	}

	return lnr.Prompt(">>> ")
}
