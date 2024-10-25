package main

import "testing"

func TestRun(t *testing.T) {
	err := run()
	if err != nil {
		t.Error("run() returned an error")
	}
}
