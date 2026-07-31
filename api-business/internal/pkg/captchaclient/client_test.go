package captchaclient

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/qixi-live/qixi-live-mergers/api-business/internal/pkg/config"
)

func TestClientSignsPTEChallengeRequest(t *testing.T) {
	const secret = "0123456789abcdef0123456789abcdef"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/challenges" || r.Method != http.MethodPost {
			t.Fatalf("unexpected request %s %s", r.Method, r.URL.Path)
		}
		if got := r.Header.Get("X-Captcha-App-ID"); got != "pte-captcha" {
			t.Fatalf("app id = %q", got)
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}
		digest := sha256.Sum256(body)
		canonical := "POST\n/api/v1/challenges\n" + hex.EncodeToString(digest[:]) + "\n" + r.Header.Get("X-Captcha-Timestamp") + "\n" + r.Header.Get("X-Captcha-Nonce")
		mac := hmac.New(sha256.New, []byte(secret))
		_, _ = mac.Write([]byte(canonical))
		if got, want := r.Header.Get("X-Captcha-Signature"), hex.EncodeToString(mac.Sum(nil)); got != want {
			t.Fatalf("signature mismatch: got %q want %q", got, want)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"challenge","challenge_id":"challenge-1","mode":"puzzle"}`))
	}))
	defer server.Close()

	client, err := New(config.CaptchaConfig{
		Enabled: true, BaseURL: server.URL, ApplicationID: "pte-captcha", SecretValue: secret,
	})
	if err != nil {
		t.Fatal(err)
	}
	response, err := client.Create(t.Context(), map[string]any{"action": "login_password", "preferred_mode": "puzzle"})
	if err != nil {
		t.Fatal(err)
	}
	if string(response) != `{"status":"challenge","challenge_id":"challenge-1","mode":"puzzle"}` {
		t.Fatalf("unexpected response: %s", response)
	}
}

func TestNewRejectsMissingServiceSecret(t *testing.T) {
	_, err := New(config.CaptchaConfig{Enabled: true, BaseURL: "http://captcha", ApplicationID: "qixi"})
	if err != ErrUnavailable {
		t.Fatalf("error = %v, want %v", err, ErrUnavailable)
	}
}
