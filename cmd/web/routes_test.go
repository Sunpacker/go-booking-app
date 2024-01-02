package main

import (
	"fmt"
	"github.com/go-chi/chi"
	"testing"
)

func TestRoutes(t *testing.T) {
	mux := routes()

	switch typeReceived := mux.(type) {
	case *chi.Mux:
	default:
		t.Error(fmt.Sprintf("type is not *chi.Mux, but is %T", typeReceived))
	}
}
