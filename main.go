package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/a-h/templ"
	"github.com/masonictemple4/gfolio/components"
)

var (
	hostPtr = flag.String("host", "", "host to serve on")
	portPtr = flag.String("port", "8080", "port to serve on")
)

func main() {

	// Serve the static files.
	fs := http.FileServer(http.Dir("assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	indexComponent := components.Index()

	// Render via server.

	http.Handle("/", templ.Handler(indexComponent))

	var hostStr string
	if *hostPtr == "" {
		hostStr = fmt.Sprintf(":%s", *portPtr)
	} else {
		hostStr = fmt.Sprintf("%s:%s", *hostPtr, *portPtr)
	}

	fmt.Printf("Listening on http://%s\n", hostStr)

	http.ListenAndServe(hostStr, nil)

	// random comment.

}
