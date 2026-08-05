package wechatmini

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCode2Session(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("js_code"); got != "one-time-code" {
			t.Fatalf("js_code = %q", got)
		}
		_, _ = w.Write([]byte(`{"openid":"openid-1","unionid":"unionid-1","session_key":"ignored"}`))
	}))
	defer server.Close()

	session, err := (&Client{Endpoint: server.URL}).Code2Session(context.Background(), Config{AppID: "wx-test", AppSecret: "secret"}, "one-time-code")
	if err != nil {
		t.Fatal(err)
	}
	if session.OpenID != "openid-1" || session.UnionID != "unionid-1" {
		t.Fatalf("session = %#v", session)
	}
}

func TestCode2SessionRejectsProviderError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"errcode":40029,"errmsg":"invalid code"}`))
	}))
	defer server.Close()
	_, err := (&Client{Endpoint: server.URL}).Code2Session(context.Background(), Config{AppID: "wx-test", AppSecret: "secret"}, "bad-code")
	if err != ErrInvalidCode {
		t.Fatalf("error = %v", err)
	}
}
