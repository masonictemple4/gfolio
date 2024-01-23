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
	"github.com/joho/godotenv"
	"github.com/masonictemple4/masonictempl/components"
	"github.com/masonictemple4/masonictempl/db"
	"github.com/masonictemple4/masonictempl/handlers"
	"github.com/masonictemple4/masonictempl/internal/filestore"
	"github.com/masonictemple4/masonictempl/internal/utils"
	"github.com/masonictemple4/masonictempl/services"
	"github.com/spf13/cobra"
)

var (
	hostPtr   *string
	portPtr   *string
	staticPtr *string
)

var defaultConfig string = "/etc/env/.masonictempl.env"

func init() {

	hostPtr = rootCmd.Flags().String("host", "", "host to serve on")
	portPtr = rootCmd.Flags().String("port", os.Getenv("PORT"), "port to serve on. defaults to the PORT env variable.")
	staticPtr = rootCmd.Flags().String("static", "assets", "path to static files directory.")
	rootCmd.PersistentFlags().String("workdir", os.Getenv("WORKDIR"), "REQUIRED. This will ensure you don't have problems with path generation and saving later on. Ideally this is the root of your project but as long as it's the same throughout it doesn't really matter what it is.")
	rootCmd.PersistentFlags().String("config", defaultConfig, "path to your .env file.")
	rootCmd.PersistentFlags().String("pub", os.Getenv("ASSET_DIR"), "name of your public/static file directory.")
	rootCmd.PersistentFlags().String("storage", os.Getenv("STORAGE_METHOD"), "Storage interface. Options are: internal, gcp.")
}

func loadDefaults(cmd *cobra.Command) {
	cmd.PersistentFlags().Lookup("config").Value.Set(defaultConfig)
	cmd.PersistentFlags().Lookup("workdir").Value.Set(os.Getenv("WORKDIR"))
	cmd.Flags().Lookup("port").Value.Set(os.Getenv("PORT"))
	cmd.PersistentFlags().Lookup("pub").Value.Set(os.Getenv("ASSET_DIR"))
	cmd.PersistentFlags().Lookup("storage").Value.Set(os.Getenv("STORAGE_METHOD"))
}

var rootCmd = &cobra.Command{
	Use:   "masonictempl",
	Short: "CLI Interface to interact with masonictempl backend",
	Long:  "CLI Interface to interact with masonictempl backend",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		configFile, err := cmd.PersistentFlags().GetString("config")
		if err != nil {
			log.Fatalf("[config flag]: %v\n", err)
		}

		if utils.FileExists(configFile) {
			if err := godotenv.Load(configFile); err != nil {
				log.Fatalf("[loading env]: %v\n", err)
			}
			loadDefaults(cmd)
		}

		storage, err := cmd.PersistentFlags().GetString("storage")
		if err != nil {
			log.Fatalf("[storage flag]: %v\n", err)
		}

		if storage == "" {
			log.Fatal("storage flag is required or set STORAGE_METHOD env variable")
		}

		wrkDir, err := cmd.PersistentFlags().GetString("workdir")
		if err != nil {
			log.Fatal(err)
		}

		if wrkDir == "" && os.Getenv("WORKDIR") == "" {
			log.Fatal("BOTH Flag and Env setting were empty WORKDIR environment setting is required")
		}

		if wrkDir != "" {
			os.Setenv("WORKDIR", wrkDir)
		}

		if os.Getenv("WORKDIR") == "" {
			log.Fatal("WORKDIR environment setting is required")
		}

	},
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) > 0 {
			cmd.Help()
			return
		}

		startServer()

	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func startServer() {
	if *staticPtr == "" {
		pwd := os.Getenv("WORKDIR")
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

	blogDb := db.NewPostgresGCPProxy()

	blogService := services.NewBlogService(blogDb)

	hndlr := handlers.NewDefaultHandler()
	hndlr.AssetPath = *staticPtr

	fh, err := filestore.NewInternalStore(*staticPtr)
	if err != nil {
		log.Fatal(err)
	}

	blogHandler := handlers.NewBlogsHandler(blogService, fh)

	resumeHandler := handlers.NewResumeHandler(fh)

	hndlr.Routes = map[string]http.Handler{
		"assets": http.StripPrefix(fsStr, fs),
		"/":      templ.Handler(components.Index()),
		"":       templ.Handler(components.Index()),
		"blog":   blogHandler,
		"resume": resumeHandler,
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
