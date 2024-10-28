package main

import "testing"

func TestRun(t *testing.T) {
	db, err := run()
	defer db.SQL.Close()
	if err != nil {
		t.Error("run() returned an error")
	}
}
