package main

import (
	"errors"
	"fmt"
	"html"
	"log"
	"net/http"
	"time"

	"github.com/getsentry/sentry-go"
)

func main() {
	err := sentry.Init(sentry.ClientOptions{
		Dsn:         env.Dsn,
		Debug:       true,
		Environment: "dev",
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
	defer sentry.Flush(2 * time.Second)

	http.HandleFunc("/bar", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))

		ctx := r.Context()
		hub := sentry.CurrentHub()
		hub.Scope().SetRequest(r)
		ctx = sentry.SetHubOnContext(ctx, hub)

		err := errors.New("errors new")
		if err != nil {
			hub.Recover(err)
		}
	})
	http.ListenAndServe(":8080", nil)

	// hub.CaptureMessage("It works!!")

}
