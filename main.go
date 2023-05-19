package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	opsLabels = []string{"request_method", "request_path", "authenticated"}

	opsRequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "ezauth_http_requests_total",
		Help: "The total number of handled http request",
	}, opsLabels)

	opsRequestsDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "ezauth_http_requests_duration",
		Help: "The duration of handled http request",
	}, opsLabels)
)

func main() {
	var err error
	defer func(err *error) {
		if *err != nil {
			log.Println("exited with error:", (*err).Error())
			os.Exit(1)
		}
	}(&err)

	var (
		optPort     = strings.TrimSpace(os.Getenv("PORT"))
		optTarget   = strings.TrimSpace(os.Getenv("PROXY_TARGET"))
		optRealm    = strings.TrimSpace(os.Getenv("BASICAUTH_REALM"))
		optUsername = strings.TrimSpace(os.Getenv("BASICAUTH_USERNAME"))
		optPassword = strings.TrimSpace(os.Getenv("BASICAUTH_PASSWORD"))

		optTargetInsecure, _ = strconv.ParseBool(strings.TrimSpace(os.Getenv("PROXY_TARGET_INSECURE")))
	)

	if optPort == "" {
		optPort = "80"
	}
	if optRealm == "" {
		optRealm = "BasicAuth Proxy"
	}

	if optTarget == "" {
		err = errors.New("missing environment variable PROXY_TARGET")
		return
	}

	if optUsername == "" {
		err = errors.New("missing environment variable BASICAUTH_USERNAME")
		return
	}

	if optPassword == "" {
		err = errors.New("missing environment variable BASICAUTH_PASSWORD")
		return
	}

	var s *http.Server
	if s, err = newServer(serverOptions{
		port:     optPort,
		target:   optTarget,
		realm:    optRealm,
		username: optUsername,
		password: optPassword,

		targetInsecure: optTargetInsecure,
	}); err != nil {
		return
	}

	chErr := make(chan error, 1)
	chSig := make(chan os.Signal, 1)
	signal.Notify(chSig, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		chErr <- s.ListenAndServe()
	}()

	select {
	case err = <-chErr:
		return
	case sig := <-chSig:
		log.Println("signal caught:", sig.String())
	}

	err = s.Shutdown(context.Background())
}
