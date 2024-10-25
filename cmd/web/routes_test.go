package main

import "testing"

func TestRoutes(t *testing.T) {
	routes := routes()
	if routes == nil {
		t.Error("routes() returned nil")
	}
}
