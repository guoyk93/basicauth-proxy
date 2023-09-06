package main

import (
	"crypto/tls"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type serverOptions struct {
	port     string
	target   string
	realm    string
	username string
	password string

	targetInsecure bool
}

func newServer(opts serverOptions) (s *http.Server, err error) {
	var target *url.URL
	if target, err = url.Parse(opts.target); err != nil {
		return
	}

	rh := httputil.NewSingleHostReverseProxy(target)
	rh.Transport = &http.Transport{
		Proxy:             http.ProxyFromEnvironment,
		DialContext:       (&net.Dialer{}).DialContext,
		ForceAttemptHTTP2: true,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: opts.targetInsecure,
		},
	}

	ph := promhttp.Handler()

	s = &http.Server{
		Addr: ":" + opts.port,
		Handler: http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			// metrics
			if req.URL.Path == "/metrics" {
				ph.ServeHTTP(rw, req)
				return
			}
			// ready
			if req.URL.Path == "/ready" {
				http.Error(rw, "OK", http.StatusOK)
				return
			}

			var (
				startedAt  = time.Now()
				authorized = false
			)

			username, password, ok := req.BasicAuth()

			if (!ok) ||
				(username != opts.username) ||
				(password != opts.password) {

				rw.Header().Set("WWW-Authenticate", "Basic realm="+strconv.Quote(opts.realm))
				http.Error(rw, "Unauthorized", http.StatusUnauthorized)
			} else {
				authorized = true
				req.Header.Del("Authorization")
				rh.ServeHTTP(rw, req)
			}

			duration := time.Since(startedAt)
			labels := prometheus.Labels{
				"request_method": req.Method,
				"request_path":   req.URL.Path,
				"authenticated":  strconv.FormatBool(authorized),
			}
			opsRequestsTotal.With(labels).Inc()
			opsRequestsDuration.With(labels).Observe(float64(duration/time.Millisecond) / float64(time.Second/time.Millisecond))
		}),
	}
	return
}
