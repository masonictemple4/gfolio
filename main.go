package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	// "github.com/a-h/templ"

	"github.com/a-h/templ"
	"github.com/masonictemple4/masonictempl/components"
	"github.com/masonictemple4/masonictempl/handlers"
)

var (
	hostPtr   = flag.String("host", "", "host to serve on")
	portPtr   = flag.String("port", "8080", "port to serve on")
	staticPtr = flag.String("static", "assets", "path to static files directory.")
)

func main() {

	if *staticPtr == "" {
		pwd, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}

		if !confirmationLoop(pwd) {
			println("Exiting...")
			os.Exit(0)
		}

	}

	// TODO: Make sure this will be supported in prod.
	fsStr := fmt.Sprintf("/%s/", *staticPtr)
	fs := http.FileServer(http.Dir(*staticPtr))

	var hostStr string
	if *hostPtr == "" {
		hostStr = fmt.Sprintf(":%s", *portPtr)
	} else {
		hostStr = fmt.Sprintf("%s:%s", *hostPtr, *portPtr)
	}

	// TODO: Might make more sense to define it with some DI
	// Like first specifying the store, then the service, then the handler.
	hndlr := handlers.NewDefaultHandler()
	hndlr.AssetPath = *staticPtr
	blogHandler := handlers.NewBlogsHandler()
	hndlr.Routes = map[string]http.Handler{
		"assets": http.StripPrefix(fsStr, fs),
		"/":      templ.Handler(components.Index()),
		"":       templ.Handler(components.Index()),
		"blog":   blogHandler,
	}

	server := &http.Server{
		Addr:         hostStr,
		Handler:      hndlr,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}

	fmt.Printf("Listening on http://%s\n", hostStr)
	server.ListenAndServe()

}

func confirmationLoop(path string) bool {
	promptMsg := fmt.Sprintf("WARNING: You are hosting the %s directory. Is this what you want? (yes/no): ", path)

	scanner := bufio.NewScanner(os.Stdin)

	print(promptMsg)

	for scanner.Scan() {
		text := scanner.Text()
		text = strings.ToLower(strings.TrimSpace(text))

		switch text {
		case "yes", "y":
			return true
		case "no", "n":
			return false
		default:
			print(promptMsg)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error:", err)
	}

	return false
}
