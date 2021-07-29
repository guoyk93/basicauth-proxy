package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

var (
	optPort     = strings.TrimSpace(os.Getenv("PORT"))
	optTarget   = strings.TrimSpace(os.Getenv("PROXY_TARGET"))
	optRealm    = strings.TrimSpace(os.Getenv("BASICAUTH_REALM"))
	optUsername = strings.TrimSpace(os.Getenv("BASICAUTH_USERNAME"))
	optPassword = strings.TrimSpace(os.Getenv("BASICAUTH_PASSWORD"))
)

func main() {
	var err error
	defer func(err *error) {
		if *err != nil {
			log.Println("exited with error:", (*err).Error())
			os.Exit(1)
		}
	}(&err)

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

	var target *url.URL
	if target, err = url.Parse(optTarget); err != nil {
		return
	}

	rp := httputil.NewSingleHostReverseProxy(target)

	s := http.Server{
		Addr: "0.0.0.0:80",
		Handler: http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			username, password, ok := req.BasicAuth()
			log.Println("Auth", username, password, ok)
			log.Println("Required", optUsername, optPassword)
			if !ok || username != optUsername || password != optPassword {
				rw.Header().Set("WWW-Authenticate", "Basic realm="+strconv.Quote(optRealm))
				http.Error(rw, "Unauthorized", http.StatusUnauthorized)
				return
			}
			req.Header.Del("Authorization")
			log.Println("Passed")
			rp.ServeHTTP(rw, req)
		}),
	}
	defer s.Shutdown(context.Background())

	chErr := make(chan error, 1)
	chSig := make(chan os.Signal, 1)
	signal.Notify(chSig, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		chErr <- s.ListenAndServe()
	}()

	select {
	case err = <-chErr:
	case sig := <-chSig:
		log.Println("signal caught:", sig.String())
	}
}
