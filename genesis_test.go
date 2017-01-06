package genesis_test

import (
	"testing"

	"github.com/wx13/genesis"
)

func TestIsRunning(t *testing.T) {

	running, err := genesis.IsRunning("asjdofsj983sjdf98ee8jfsnviidf")
	if running {
		t.Error("Crazy process should not be running:", err, running)
	}

}

func TestExpandHome(t *testing.T) {

	dir := "hello/there"
	expDir := genesis.ExpandHome(dir)
	if dir != expDir {
		t.Error("Path should not have changed:", dir, "=>", expDir)
	}

}
