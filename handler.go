//revive:disable:package-comments

package app

import (
	"fmt"
	"io/fs"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"git.sonicoriginal.software/logger.git"
	"git.sonicoriginal.software/server.git/v2"
)

const name = "app"

// Handler handles App requests
type handler struct {
	logger logger.Log
	fsys   fs.FS
}

func (handler *handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	handler.logger.Info("%v %v\n", request.Method, request.URL.Path)
	requestPath := strings.TrimPrefix(request.URL.Path, fmt.Sprintf("/%v/", name))
	if filepath.Ext(requestPath) == "" {
		requestPath = "index.html"
	}

	file, err := handler.fsys.Open(requestPath)
	if err != nil {
		handler.logger.Error("Could not open resource at: %v\n", requestPath)
		http.Error(writer, err.Error(), http.StatusNotFound)
		return
	}

	stats, err := file.Stat()
	if err != nil {
		handler.logger.Error("Could not stat resource at: %v\n", requestPath)
		http.Error(writer, err.Error(), http.StatusNoContent)
		return
	}

	contents := make([]byte, stats.Size())
	_, err = file.Read(contents)
	if err != nil {
		handler.logger.Error("Could not read resource at: %v\n", requestPath)
		http.Error(writer, err.Error(), http.StatusNoContent)
		return
	}

	fileExtension := filepath.Ext(requestPath)
	contentType := mime.TypeByExtension(fileExtension)
	writer.Header().Set("Content-Type", contentType)

	if _, err = writer.Write(contents); err != nil {
		handler.logger.Error("Could not write response: %v", err)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
}

// New returns a new Handler
func New(fsys fs.FS, mux *http.ServeMux) (route string) {
	logger := logger.New(
		name,
		logger.DefaultSeverity,
		os.Stdout,
		os.Stderr,
	)

	if mux == nil {
		mux = http.DefaultServeMux
	}

	return server.RegisterHandler(name, &handler{logger, fsys}, mux)
}
