package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/QuantumNous/new-api/setting"
	"github.com/QuantumNous/new-api/setting/operation_setting"
	"github.com/gin-gonic/gin"
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

func TestGetTopUpInfo_CreemDisabled(t *testing.T) {
	gin.SetMode(gin.TestMode)

	prevEnabled := setting.CreemEnabled
	prevApiKey := setting.CreemApiKey
	prevProducts := setting.CreemProducts
	defer func() {
		setting.CreemEnabled = prevEnabled
		setting.CreemApiKey = prevApiKey
		setting.CreemProducts = prevProducts
	}()

	setting.CreemEnabled = false
	setting.CreemApiKey = "test-api-key"
	setting.CreemProducts = `[{"productId":"prod_1","name":"Test","price":4.99,"currency":"USD","quota":100000}]`

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	GetTopUpInfo(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	data, ok := resp["data"].(map[string]interface{})
	if !ok {
		t.Fatal("expected data to be a map")
	}

	if data["enable_creem_topup"] != false {
		t.Fatalf("expected enable_creem_topup=false when CreemEnabled=false, got %v", data["enable_creem_topup"])
	}
}

func TestGetTopUpInfo_CreemEnabled(t *testing.T) {
	gin.SetMode(gin.TestMode)

	prevEnabled := setting.CreemEnabled
	prevApiKey := setting.CreemApiKey
	prevProducts := setting.CreemProducts
	defer func() {
		setting.CreemEnabled = prevEnabled
		setting.CreemApiKey = prevApiKey
		setting.CreemProducts = prevProducts
	}()

	setting.CreemEnabled = true
	setting.CreemApiKey = "test-api-key"
	setting.CreemProducts = `[{"productId":"prod_1","name":"Test","price":4.99,"currency":"USD","quota":100000}]`

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	GetTopUpInfo(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	data, ok := resp["data"].(map[string]interface{})
	if !ok {
		t.Fatal("expected data to be a map")
	}

	if data["enable_creem_topup"] != true {
		t.Fatalf("expected enable_creem_topup=true when CreemEnabled=true, got %v", data["enable_creem_topup"])
	}
}

func TestGetTopUpInfo_CreemDisabled_FiltersCreemFromPayMethods(t *testing.T) {
	gin.SetMode(gin.TestMode)

	prevEnabled := setting.CreemEnabled
	prevApiKey := setting.CreemApiKey
	prevProducts := setting.CreemProducts
	prevPayMethods := operation_setting.PayMethods
	defer func() {
		setting.CreemEnabled = prevEnabled
		setting.CreemApiKey = prevApiKey
		setting.CreemProducts = prevProducts
		operation_setting.PayMethods = prevPayMethods
	}()

	setting.CreemEnabled = false
	setting.CreemApiKey = "test-api-key"
	setting.CreemProducts = `[{"productId":"prod_1","name":"Test","price":4.99,"currency":"USD","quota":100000}]`
	operation_setting.PayMethods = []map[string]string{
		{"name": "支付宝", "type": "alipay"},
		{"name": "Creem", "type": "creem"},
		{"name": "微信", "type": "wxpay"},
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	GetTopUpInfo(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	data, ok := resp["data"].(map[string]interface{})
	if !ok {
		t.Fatal("expected data to be a map")
	}

	payMethodsRaw, ok := data["pay_methods"].([]interface{})
	if !ok {
		t.Fatal("expected pay_methods to be a slice")
	}

	for _, m := range payMethodsRaw {
		method, ok := m.(map[string]interface{})
		if !ok {
			continue
		}
		if method["type"] == "creem" {
			t.Fatal("expected pay_methods to not contain creem when CreemEnabled=false")
		}
	}

	if len(payMethodsRaw) != 2 {
		t.Fatalf("expected 2 pay_methods after filtering creem, got %d", len(payMethodsRaw))
	}
}

func TestRequestCreemPay_Disabled(t *testing.T) {
	gin.SetMode(gin.TestMode)

	prevEnabled := setting.CreemEnabled
	defer func() {
		setting.CreemEnabled = prevEnabled
	}()
	setting.CreemEnabled = false

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body, _ := json.Marshal(CreemPayRequest{
		ProductId:     "prod_1",
		PaymentMethod: PaymentMethodCreem,
	})
	c.Request, _ = http.NewRequest("POST", "/api/user/creem/pay", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	creemAdaptor.RequestPay(c, &CreemPayRequest{
		ProductId:     "prod_1",
		PaymentMethod: PaymentMethodCreem,
	})

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp["message"] != "error" {
		t.Fatalf("expected message=error, got %v", resp["message"])
	}
	if resp["data"] != "Creem 支付未启用" {
		t.Fatalf("expected data='Creem 支付未启用', got %v", resp["data"])
	}
}
