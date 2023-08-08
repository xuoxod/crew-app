package main

import (
	"fmt"
	"testing"
)

func TestRun(t *testing.T) {
	drvr, err := run()

	if err != nil {
		t.Error("failed run")
	}

	fmt.Println("driver: ", drvr)
}
