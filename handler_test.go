package app_test

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"testing"
	"testing/fstest"

	app "git.sonicoriginal.software/server-routes-app.git"
	"git.sonicoriginal.software/server.git/v2"
)

const (
	portEnvKey       = "TEST_PORT"
	htmlFileContents = `<html><head></head><body>hello</body></html>`
)

var (
	certs      []tls.Certificate
	filesystem                = fstest.MapFS{"index.html": &fstest.MapFile{Data: []byte(htmlFileContents)}}
	mux        *http.ServeMux = nil
)

func TestHandler(t *testing.T) {
	route := app.New(filesystem, mux)

	t.Logf("Handler registered for route [%v]\n", route)

	ctx, cancelFunction := context.WithCancel(context.Background())
	address, serverErrorChannel := server.Run(ctx, &certs, mux, portEnvKey)

	t.Logf("Serving on [%v]\n", address)

	url := fmt.Sprintf("http://%v%v", address, route)

	t.Logf("Requesting [%v]\n", url)

	response, err := http.DefaultClient.Get(url)
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	cancelFunction()

	serverError := <-serverErrorChannel
	if serverError.Close != nil {
		t.Fatalf("Error closing server: %v", serverError.Close.Error())
	}
	contextError := serverError.Context.Error()

	t.Logf("%v\n", contextError)

	if contextError != server.ErrContextCancelled.Error() {
		t.Fatalf("Server failed unexpectedly: %v", contextError)
	}

	t.Log("Response:")
	t.Logf("  Status code: %v", response.StatusCode)
	t.Logf("  Status text: %v", response.Status)

	if response.Status != http.StatusText(http.StatusOK) && response.StatusCode != http.StatusOK {
		t.Fatalf("Server returned: %v", response.Status)
	}

	responseBody, err := io.ReadAll(response.Body)
	responseFormatted := string(responseBody)
	if err != nil {
		t.Fatalf("Could not read response: %v", err)
	} else if responseFormatted != htmlFileContents {
		t.Fatalf("%v != %v", responseFormatted, htmlFileContents)
	}

	t.Logf("  Body:\n%v", responseFormatted)
}
