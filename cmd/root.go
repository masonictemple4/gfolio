package cmd

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/a-h/templ"
	"github.com/masonictemple4/masonictempl/components"
	"github.com/masonictemple4/masonictempl/handlers"
	"github.com/masonictemple4/masonictempl/internal/filestore"
	"github.com/spf13/cobra"
)

var (
	hostPtr   *string
	portPtr   *string
	staticPtr *string
)

var rootCmd = &cobra.Command{
	Use:   "masonictempl",
	Short: "CLI Interface to interact with masonictempl backend",
	Long:  "CLI Interface to interact with masonictempl backend",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// TODO: Here we can add config loading
		// And db setup if necessary..
		// This can be done later....

	},
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) > 0 {
			cmd.Help()
			return
		}

		startServer()

	},
}

func init() {
	hostPtr = rootCmd.Flags().String("host", "", "host to serve on")
	portPtr = rootCmd.Flags().String("port", "8080", "port to serve on")
	staticPtr = rootCmd.Flags().String("static", "assets", "path to static files directory.")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func startServer() {
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

	fh, err := filestore.NewInternalStore(*staticPtr)
	if err != nil {
		log.Fatal(err)
	}

	blogHandler := handlers.NewBlogsHandler(fh)
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
