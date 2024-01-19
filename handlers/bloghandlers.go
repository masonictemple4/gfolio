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

func (bh BlogsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	blogs := bh.BlogService.List(r.Context())

	fmt.Printf("BlogsHandler: %+v\n", blogs)

	templ.Handler(components.BlogList(blogs)).ServeHTTP(w, r)
}
