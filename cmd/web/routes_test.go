package main

import (
	"booking/internal/config"
	"fmt"
	"testing"

	"github.com/go-chi/chi"
)

func TestRoutes(t *testing.T) {
	var app config.AppConfig

	mux := routes(&app)

	switch v := mux.(type) {
	case *chi.Mux:
		//do nothin; test passed
	default:
		t.Error(fmt.Sprintf("type is not *chip.Mux, type is %T", v))
	}
}
