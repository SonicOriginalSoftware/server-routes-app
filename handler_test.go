package app_test

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"testing/fstest"

	"git.nathanblair.rocks/routes/app"
	lib "git.nathanblair.rocks/server"
)

var (
	certs            []tls.Certificate
	port             = "4430"
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
	t.Setenv("PORT", port)
	route := fmt.Sprintf("localhost/%v/", app.Name)
	t.Setenv(fmt.Sprintf("%v_SERVE_ADDRESS", strings.ToUpper(app.Name)), route)

	if _, err := app.New(filesystem); err != nil {
		t.Fatalf("%v\n", err)
	}

	ctx, cancelContext := context.WithCancel(context.Background())

	exitCode, address := lib.Run(ctx, certs)
	defer close(exitCode)

	url := fmt.Sprintf("http://%v/%v/", address, app.Name)
	response, err := http.DefaultClient.Get(url)
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	cancelContext()

	if returnCode := <-exitCode; returnCode != 0 {
		t.Fatalf("Server errored: %v", returnCode)
	}

	if response.Status != http.StatusText(http.StatusOK) && response.StatusCode != http.StatusOK {
		t.Fatalf("Server returned: %v", response.Status)
	}
}
