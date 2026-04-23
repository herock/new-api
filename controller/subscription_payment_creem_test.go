package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/QuantumNous/new-api/setting"
	"github.com/gin-gonic/gin"
)

func TestSubscriptionRequestCreemPay_Disabled(t *testing.T) {
	gin.SetMode(gin.TestMode)

	prevEnabled := setting.CreemEnabled
	defer func() {
		setting.CreemEnabled = prevEnabled
	}()
	setting.CreemEnabled = false

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body, _ := json.Marshal(SubscriptionCreemPayRequest{
		PlanId: 1,
	})
	c.Request, _ = http.NewRequest("POST", "/api/subscription/creem/pay", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	SubscriptionRequestCreemPay(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp["success"] != false {
		t.Fatalf("expected success=false, got %v", resp["success"])
	}
	if resp["message"] != "Creem 支付未启用" {
		t.Fatalf("expected message='Creem 支付未启用', got %v", resp["message"])
	}
}
