//go:build e2e

package hippo

import (
	"fmt"
	"testing"
	"time"

	"github.com/jhaals/yopass/pkg/yopass"
)

const (
	walrusURL = "http://localhost:31415"
	redisURL  = "redis://localhost:6379/0"
)

func newHippo() (*Hippo, error) {
	hippo, err := NewHippo(walrusURL, redisURL)
	if err != nil {
		return nil, err
	}
	return hippo.(*Hippo), nil
}

func generateTimestampedSecret() yopass.Secret {
	return yopass.Secret{
		Message:    fmt.Sprintf("test hippos data %s", time.Now().Format(time.RFC3339)),
		Expiration: 3600,
		OneTime:    true,
	}
}

func TestHippo_PutAndGet_E2E(t *testing.T) {
	hippo, err := newHippo()
	if err != nil {
		t.Fatalf("Failed to create Hippo instance: %v", err)
	}

	secret := generateTimestampedSecret()
	key := "test-key"

	err = hippo.Put(key, secret)
	if err != nil {
		t.Fatalf("Failed to put secret: %v", err)
	}

	retrievedSecret, err := hippo.Get(key)
	if err != nil {
		t.Fatalf("Failed to get secret: %v", err)
	}

	if retrievedSecret.Message != secret.Message {
		t.Errorf("Expected message %s, but got %s", secret.Message, retrievedSecret.Message)
	}
}

func TestHippo_Delete_E2E(t *testing.T) {
	hippo, err := newHippo()
	if err != nil {
		t.Fatalf("Failed to create Hippo instance: %v", err)
	}

	secret := generateTimestampedSecret()
	key := "test-key"
	err = hippo.Put(key, secret)
	if err != nil {
		t.Fatalf("Failed to put secret: %v", err)
	}

	deleted, err := hippo.Delete(key)
	if err != nil {
		t.Fatalf("Failed to delete secret: %v", err)
	}
	if !deleted {
		t.Errorf("Expected secret to be deleted, but it was not")
	}

	_, err = hippo.Get(key)
	if err == nil {
		t.Errorf("Expected error when getting deleted secret, but got none")
	}
}
