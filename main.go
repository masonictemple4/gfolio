package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/a-h/templ"
	"github.com/masonictemple4/masonictempl/components"
)

func main() {

	// Serve the static files.
	fs := http.FileServer(http.Dir("assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	indexComponent := components.Index()

	// Will print output to the console.
	println("Render via component.Render")
	indexComponent.Render(context.Background(), os.Stdout)
	// Give a newline after the sample render
	println()

	// Render via server.

	http.Handle("/", templ.Handler(indexComponent))

	fmt.Println("Listening on :8080")

	http.ListenAndServe(":8080", nil)

}
