package main

import (
	"net/http"
	"testing"
)

func TestNoSurf(t *testing.T) {
	var myH myHandler
	h := NoSurf(&myH)
	switch v := h.(type) {
	case http.Handler:
		// do nothing
	default:
		t.Errorf("Expected type http.Handler, Received %T", v)
	}
}

func TestSessionLoad(t *testing.T) {
	var myH myHandler
	h := SessionLoad(&myH)

	switch v := h.(type) {
	case http.Handler:
		// Do nothing
	default:
		t.Errorf("Expected type http.Handler, Received %T", v)
	}
}
