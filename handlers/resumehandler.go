package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/a-h/templ"
	"github.com/masonictemple4/masonictempl/components"
	"github.com/masonictemple4/masonictempl/internal/filestore"
	"github.com/masonictemple4/masonictempl/internal/parser"
)

func NewResumeHandler(fh filestore.Filestore) ResumeHandler {
	lgr := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))
	return ResumeHandler{
		DefaultHandler: DefaultHandler{
			Log: lgr,
		},
		FileHandler: fh,
	}
}

type ResumeHandler struct {
	DefaultHandler
	FileHandler filestore.Filestore
}

func (rh ResumeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// PATH: /assets/resume.pdf
	// PATH: /assets/resume.md
	pdfHref := "/assets/resume.pdf"

	fp := os.Getenv("RESUME_MD_PATH")

	if len(strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")) > 1 && strings.Contains(r.URL.Path, "download") {
		println("SHOULD BE DOWNLOADING")
		if fp == "" {
			// "assets/resume.pdf"
			fp = filepath.Join(filestore.GetRootPath(rh.FileHandler), "resume.pdf")
		}

		resumeData, err := rh.FileHandler.Read(r.Context(), fp)
		if err != nil {
			rh.Log.Error("error loading resume markdown: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Disposition", "attachment; filename=resume.pdf")
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(resumeData)))

		// Write the content of the file
		w.Write(resumeData)
		return
	}

	if fp == "" {
		// "assets/resume.md"
		fp = filepath.Join(filestore.GetRootPath(rh.FileHandler), "resume.md")
	}

	resumeData, err := rh.FileHandler.Read(r.Context(), fp)
	if err != nil {
		rh.Log.Error("error loading resume markdown: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Although we know this is clean html, it's good
	// to be in the habit of parsing the frontmatter out.
	resumeData, err = parser.SkipFrontmatter(resumeData)
	if err != nil {
		rh.Log.Error("error parsing resume markdown: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	templ.Handler(components.Resume(pdfHref, string(resumeData))).ServeHTTP(w, r)

}
