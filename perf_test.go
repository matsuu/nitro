package nitro

import (
	"testing"
)

func TestStart(t *testing.T) {
	b := Start("test start")
	b.Step("test step 1")
	b.Step("test step 2")
	b.Step("test step 3")
	b.Stop("test stop")
}
