package main

import (
	"fmt"
	"github.com/Sunpacker/go-booking-app/pkg/handlers"
	"net/http"
)

const PORT = ":8080"

func main() {
	http.HandleFunc("/", handlers.Home)
	http.HandleFunc("/about", handlers.About)

	fmt.Printf("Starting application on port %s", PORT)
	_ = http.ListenAndServe(PORT, nil)
}
