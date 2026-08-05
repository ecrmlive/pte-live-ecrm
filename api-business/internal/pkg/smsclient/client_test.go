package smsclient

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestSendRejectsNonHTTPSAndMissingSecret(t *testing.T) {
	c := New()
	if err := c.Send(context.Background(), Config{Endpoint: "http://example.test", Authorization: "x", Template: "login"}, "13800000000", "123456", time.Minute); err != ErrNotConfigured {
		t.Fatalf("err=%v", err)
	}
	if err := c.Send(context.Background(), Config{Endpoint: "https://example.test", Template: "login"}, "13800000000", "123456", time.Minute); err != ErrNotConfigured {
		t.Fatalf("err=%v", err)
	}
}

func TestSendUsesAuthorizationAndRejectsFailure(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer local-test" {
			w.WriteHeader(401)
			return
		}
		w.WriteHeader(202)
	}))
	defer server.Close()
	c := New()
	c.httpClient = server.Client()
	if err := c.Send(context.Background(), Config{Endpoint: server.URL, Authorization: "Bearer local-test", Template: "login"}, "13800000000", "123456", time.Minute); err != nil {
		t.Fatal(err)
	}
}
