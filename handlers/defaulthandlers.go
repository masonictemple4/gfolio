package handlers

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
)

type DefaultHandler struct {
	Log       *slog.Logger
	Routes    map[string]http.Handler
	AssetPath string
}

func NewDefaultHandler() *DefaultHandler {
	Log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))
	Routes := make(map[string]http.Handler, 0)
	return &DefaultHandler{
		Log:    Log,
		Routes: Routes,
	}
}

func (dh *DefaultHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	urlParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")
	ctx := context.WithValue(r.Context(), "path", r.URL.Path)

	r = r.WithContext(ctx)

	part := urlParts[0]
	if part == "" {
		part = "/"
	}

	if part == "favicon.ico" {
		favPath := fmt.Sprintf("/%s/favicon.ico", dh.AssetPath)
		if _, err := os.Stat(favPath); os.IsNotExist(err) {
			return
		}
		http.Redirect(w, r, favPath, http.StatusMovedPermanently)
		return
	}

	handler, ok := dh.Routes[urlParts[0]]
	if !ok {

		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	handler.ServeHTTP(w, r)
}
