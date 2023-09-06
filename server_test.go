package main

import (
	"context"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServer(t *testing.T) {
	ts := &http.Server{
		Addr: ":12123",
		Handler: http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.Write([]byte(req.URL.Path))
		}),
	}
	defer ts.Shutdown(context.Background())

	ps, err := newServer(serverOptions{
		port:     "12124",
		target:   "http://127.0.0.1:12123",
		realm:    "rko",
		username: "hello",
		password: "world",
	})
	require.NoError(t, err)
	defer ps.Shutdown(context.Background())

	fetch := func(path string, auth bool) (code int, headers http.Header, body string, err error) {
		var req *http.Request
		if req, err = http.NewRequest(http.MethodGet, "http://127.0.0.1:12124"+path, nil); err != nil {
			return
		}
		if auth {
			req.SetBasicAuth("hello", "world")
		}
		var res *http.Response
		if res, err = http.DefaultClient.Do(req); err != nil {
			return
		}
		defer res.Body.Close()
		code = res.StatusCode
		headers = res.Header
		var buf []byte
		if buf, err = io.ReadAll(res.Body); err != nil {
			return
		}
		body = string(buf)
		return
	}

	go func() {
		ts.ListenAndServe()
	}()

	go func() {
		ps.ListenAndServe()
	}()

	time.Sleep(time.Second)

	code, _, _, err := fetch("/metrics", false)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, code)

	code, _, body, err := fetch("/ready", false)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, code)
	assert.Equal(t, "OK\n", body)

	code, headers, body, err := fetch("/test", false)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, code)
	assert.Equal(t, `Basic realm="rko"`, headers.Get("WWW-Authenticate"))
	assert.Equal(t, "Unauthorized\n", body)

	code, _, body, err = fetch("/test", true)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, code)
	assert.Equal(t, "/test", body)
}
