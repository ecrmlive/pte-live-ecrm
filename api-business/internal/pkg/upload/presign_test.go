package upload

import (
	"context"
	"strings"
	"testing"
)

func TestCOSPresignPutReturnsOnlyDirectPUTContract(t *testing.T) {
	storage := COS{Bucket: "demo-1250000000", Region: "ap-shanghai", SecretID: "test-id", SecretKey: "test-key", KeyPrefix: "pte-live-ecrm"}
	intent, err := storage.PresignPut(context.Background(), PresignInput{Scope: "app/merchant-applications/7", Filename: "license.png", ContentType: "image/png", Size: 1024})
	if err != nil {
		t.Fatal(err)
	}
	if intent.Method != "PUT" || !strings.HasPrefix(intent.ObjectKey, "pte-live-ecrm/app/merchant-applications/7/") || intent.Headers["Content-Type"] != "image/png" || !strings.Contains(intent.UploadURL, "q-sign-algorithm") {
		t.Fatalf("unexpected intent: %+v", intent)
	}
}

func TestCOSPresignPutRejectsInvalidFiles(t *testing.T) {
	storage := COS{Bucket: "demo-1250000000", Region: "ap-shanghai", SecretID: "test-id", SecretKey: "test-key"}
	if _, err := storage.PresignPut(context.Background(), PresignInput{Scope: "app/merchant-applications/7", Filename: "bad.exe", Size: 1024}); err == nil {
		t.Fatal("expected extension rejection")
	}
	if _, err := storage.PresignPut(context.Background(), PresignInput{Scope: "../escape", Filename: "license.png", Size: 1024}); err == nil {
		t.Fatal("expected scope rejection")
	}
}
