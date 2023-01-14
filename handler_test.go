package app_test

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"testing/fstest"

	"git.sonicoriginal.software/routes/app"
	lib "git.sonicoriginal.software/server"
)

var (
	certs            []tls.Certificate
	rootPath         = ""
	cssFileContents  = `* { padding: 0; margin: 0; }`
	jsFileContents   = `console.log("Hello, world!")`
	htmlFileContents = `<!DOCTYPE html>
<html lang="en-US">
  <head>
    <meta charset="utf-8" />
    <title>Demo Page</title>
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <link rel="stylesheet" type="text/css" media="screen" href="main.css" />
    <script src="main.js"></script>
  </head>
  <body></body>
</html>`

	filesystem = fstest.MapFS{
		fmt.Sprintf("%vindex.html", rootPath): &fstest.MapFile{
			Data: []byte(htmlFileContents),
		},
		fmt.Sprintf("%vmain.css", rootPath): &fstest.MapFile{
			Data: []byte(cssFileContents),
		},
		fmt.Sprintf("%vmain.js", rootPath): &fstest.MapFile{
			Data: []byte(jsFileContents),
		},
	}
)

func TestHandler(t *testing.T) {
	route := fmt.Sprintf("localhost/%v/", app.Name)
	t.Setenv(fmt.Sprintf("%v_SERVE_ADDRESS", strings.ToUpper(app.Name)), route)

	if _, err := app.New(filesystem); err != nil {
		t.Fatalf("%v\n", err)
	}

	ctx, cancelFunction := context.WithCancel(context.Background())
	address, errChan := lib.Run(ctx, certs)

	url := fmt.Sprintf("http://%v/%v/", address, app.Name)
	response, err := http.DefaultClient.Get(url)
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	cancelFunction()

	if err := <-errChan; err != nil {
		t.Fatalf("Server errored: %v", err)
	}

	// TODO Check the http response text is what we expect?

	if response.Status != http.StatusText(http.StatusOK) && response.StatusCode != http.StatusOK {
		t.Fatalf("Server returned: %v", response.Status)
	}
}
