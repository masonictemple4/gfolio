package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/a-h/templ"
	"github.com/masonictemple4/masonictempl/components"
)

var (
	hostPtr   = flag.String("host", "", "host to serve on")
	portPtr   = flag.String("port", "8080", "port to serve on")
	staticPtr = flag.String("static", "assets", "static files directory to serve. I.e css, js, images, icons etc.. Note: Directory must be inside of the project directory, paths to outside of the root project directory (i.e, /home/user/...) are not supported.")
)

func main() {

	// Serve the static files.
	fs := http.FileServer(http.Dir("assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	indexComponent := components.Index()
	blogComponent := components.BlogList()

	// Render via server.

	http.Handle("/", templ.Handler(indexComponent))
	http.Handle("/blog", templ.Handler(blogComponent))

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
