package controller

import (
	"testing"

	"github.com/QuantumNous/new-api/setting"
)

func TestVerifyCreemSignatureRequiresSecretInTestMode(t *testing.T) {
	previousTestMode := setting.CreemTestMode
	setting.CreemTestMode = true
	t.Cleanup(func() {
		setting.CreemTestMode = previousTestMode
	})

	if verifyCreemSignature(`{"eventType":"checkout.completed"}`, "", "") {
		t.Fatal("expected empty Creem webhook secret to reject signatures in test mode")
	}
}

func TestVerifyCreemSignatureRequiresSignature(t *testing.T) {
	if verifyCreemSignature(`{"eventType":"checkout.completed"}`, "", "secret") {
		t.Fatal("expected empty Creem webhook signature to be rejected")
	}
}

func TestVerifyCreemSignatureAcceptsValidSignature(t *testing.T) {
	payload := `{"eventType":"checkout.completed"}`
	secret := "secret"
	signature := generateCreemSignature(payload, secret)

	if !verifyCreemSignature(payload, signature, secret) {
		t.Fatal("expected valid Creem webhook signature to be accepted")
	}
}
