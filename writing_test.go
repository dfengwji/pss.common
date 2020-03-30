package main

import (
	"fmt"
	"github.com/dfengwji/pss.common/writing"
	"testing"
)

func TestDot(t *testing.T) {
	// <setup code>
	t.Run("A=1", func(t *testing.T) {
		dot := writing.DotInfo{}
		fmt.Println(dot.String())
	})
	// <tear-down code>
}