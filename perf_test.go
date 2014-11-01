package nitro

import (
	"math/rand"
	"testing"
)

func TestStart(t *testing.T) {
	b := Start("test start")
	b.Step("test step 1")
	b.Step("test step 2")
	b.Step("test step 3")
	b.Stop("test stop")
}

func TestCondition(t *testing.T) {
	b := Start("test start")
	SetCondition(func() bool {
		return rand.Intn(100) == 42
	})
	b.Step("test step 1")
	b.Step("test step 2")
	b.Step("test step 3")
	b.Stop("test stop")
}
