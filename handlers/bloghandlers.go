package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/a-h/templ"
	"github.com/masonictemple4/masonictempl/components"
	"github.com/masonictemple4/masonictempl/internal/filestore"
	"github.com/masonictemple4/masonictempl/internal/parser"
	"github.com/masonictemple4/masonictempl/services"
)

func NewBlogsHandler(service *services.BlogService, fh filestore.Filestore) BlogsHandler {
	// Include source path to the error or calling function in the log output.
	if service == nil {
		slog.Error("blogservice is nil")
	}

	lgr := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))

	return BlogsHandler{
		Log:         lgr,
		BlogService: service,
		Filehandler: fh,
	}
}

type BlogsHandler struct {
	Log         *slog.Logger
	Filehandler filestore.Filestore
	BlogService *services.BlogService
}

func (bh BlogsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	blogs := bh.BlogService.List(r.Context())
	urlParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")

	if len(urlParts) > 1 && !strings.Contains(strings.TrimPrefix(r.URL.Path, "/"), "assets") {
		blog, err := bh.BlogService.GetWithSlug(r.Context(), urlParts[len(urlParts)-1], "Authors", "Tags", "Media")
		if err != nil {
			fmt.Printf("Error getting blog: %v\n", err)
		}

		if bh.Filehandler == nil {
			fmt.Printf("Filehandler is nil\n")
		}

		fp := strings.Replace(blog.Docpath, "./", "", 1)
		blogData, err := bh.Filehandler.Read(r.Context(), fp)
		if err != nil {
			fmt.Printf("Error reading blog: %v\n", err)
		}

		cleanData, _ := parser.SkipFrontmatter(blogData)

		templ.Handler(components.BlogDetail(*blog, string(cleanData))).ServeHTTP(w, r)
		return
	}

	templ.Handler(components.BlogList(blogs)).ServeHTTP(w, r)
}
