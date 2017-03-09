package installer

import (
	"strings"

	"github.com/peterh/liner"
)

func Rerun(dir string) (string, error) {

	lnr := liner.NewLiner()
	defer lnr.Close()
	lnr.SetCtrlCAborts(true)

	hist := GetHistory(dir)
	for i := len(hist) - 1; i >= 0; i-- {
		words := strings.Fields(hist[i])
		if len(words) <= 1 {
			continue
		}
		lnr.AppendHistory(strings.Join(words[1:], " "))
	}

	return lnr.Prompt(">>> ")
}
