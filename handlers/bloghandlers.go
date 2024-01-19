package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/a-h/templ"
	"github.com/masonictemple4/masonictempl/components"
	"github.com/masonictemple4/masonictempl/services"
)

func NewBlogsHandler() BlogsHandler {
	// Include source path to the error or calling function in the log output.
	lgr := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))
	return BlogsHandler{
		Log:         lgr,
		BlogService: services.NewBlogService(),
	}
}

type BlogsHandler struct {
	Log         *slog.Logger
	BlogService *services.BlogService
}

func (ph BlogsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ps := ph.BlogService.List(r.Context())

	fmt.Printf("BlogsHandler: %+v\n", ps)

	templ.Handler(components.Index()).ServeHTTP(w, r)
}
