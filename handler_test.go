package app_test

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"testing"

	"git.nathanblair.rocks/routes/app"
	lib "git.nathanblair.rocks/server"
	"git.nathanblair.rocks/server/handler"
)

func TestHandler(t *testing.T) {
	subdomains := handler.Handlers{app.Prefix: app.New()}

	var certs []tls.Certificate
	ctx, cancelContext := context.WithCancel(context.Background())

	exitCode, address := lib.Run(ctx, subdomains, certs)
	defer close(exitCode)

	url := fmt.Sprintf("http://%v.%v", app.Prefix, address)
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
