package handlers

import (
	"github.com/Sunpacker/go-booking-app/pkg/render"
	"net/http"
)

// Home is the home page handler
func Home(w http.ResponseWriter, r *http.Request) {
	_ = r
	render.RenderTemplate(w, "home.page")
}

func About(w http.ResponseWriter, r *http.Request) {
	_ = r
	render.RenderTemplate(w, "about.page")
}
