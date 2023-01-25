package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/maxatome/go-testdeep/td"
)

func TestEnableCORS(t *testing.T) {
	rr := httptest.NewRecorder()

	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// a mock http handler to pass to enableCORS middleware
	// which write a 200 OK status and an "OK OVH" response body
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK OVH"))
	})

	// write response into httptest recorder
	enableCORS(next).ServeHTTP(rr, r)

	rs := rr.Result()

	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	bytes.TrimSpace(body)

	td.Cmp(t, rs.StatusCode, http.StatusOK)
	td.Cmp(t, rs.Header.Get("Access-Control-Allow-Origin"), "*")
	td.Cmp(t, string(body), "OK OVH")
}
