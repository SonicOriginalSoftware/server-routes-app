//revive:disable:package-comments

package app

import (
	"fmt"
	"io/fs"
	"mime"
	"net/http"
	"path/filepath"
	"strings"

	"git.sonicoriginal.software/server/handlers"
	"git.sonicoriginal.software/server/logging"
)

const defaultServePath = "public"

// Name is the name used to identify the service
const Name = "app"

// Handler handles App requests
type Handler struct {
	logger logging.Log
	fsys   fs.FS
}

func (handler *Handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	handler.logger.Info("%v %v\n", request.Method, request.URL.Path)
	requestPath := strings.TrimPrefix(request.URL.Path, fmt.Sprintf("/%v/", Name))
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
func New(fsys fs.FS) (handler *Handler, err error) {
	logger := logging.New(Name)
	handler = &Handler{logger, fsys}
	handlers.Register(Name, "", Name, handler, logger)

	return
}
