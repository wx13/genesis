package installer

import (
	"fmt"

	"github.com/wx13/genesis"
)

var StatusCount struct {
	Pass, Fail, Unknown, Done int
}

func ReportSummary() {
	fmt.Println("")
	fmt.Println("    Summary:")
	fmt.Println("      Pass:   ", StatusCount.Pass)
	fmt.Println("      Done:   ", StatusCount.Done)
	fmt.Println("      Unknown:", StatusCount.Unknown)
	if StatusCount.Fail > 0 {
		fmt.Println("      \033[31mFail:   ", StatusCount.Fail, "\033[0m")
	} else {
		fmt.Println("      Fail:   ", StatusCount.Fail)
	}
}

func ReportDone(msg string, err error) {
	StatusCount.Done++
	fmt.Println("    \033[1;32m[DONE]\033[0m", msg)
	if err != nil {
		fmt.Println("    ", err)
	}
}

func ReportPass(msg string, err error) {
	StatusCount.Pass++
	fmt.Println("    \033[32m[PASS]\033[0m", msg)
	if err != nil {
		fmt.Println("   ", err)
	}
}

func ReportFail(msg string, err error) {
	StatusCount.Fail++
	fmt.Println("    \033[31m[FAIL]\033[0m", msg)
	if err != nil {
		fmt.Println("   ", err)
	}
}

func ReportUnknown(msg string, err error) {
	StatusCount.Unknown++
	fmt.Println("    \033[33m[UNKNOWN]\033[0m", msg)
	if err != nil {
		fmt.Println("   ", err)
	}
}

func PrintHeader(tag, desc string) {
	fmt.Println("")
	id := "\033[36m" + genesis.StringHash(tag) + "\033[0m"
	fmt.Println("   ", id, desc)
}

func PrintSectionHeader(name string) {
	if name == "" {
		return
	}
	fmt.Println("")
	id := "\033[36m" + genesis.StringHash(name) + "\033[0m"
	fmt.Println("    ======== ", id, name, "========")
}

func PrintSectionFooter(name string) {
	if name == "" {
		return
	}
	fmt.Println("")
	id := "\033[36m" + genesis.StringHash(name) + "\033[0m"
	fmt.Println("    -------- ", id, name, "--------")
}
