package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/QuantumNous/new-api/common"
	"github.com/bytedance/gopkg/util/gopool"
)

type AdminTelegramEvent string

const (
	AdminEventUserRegistered AdminTelegramEvent = "admin_user_registered"
	AdminEventTopUpCreated   AdminTelegramEvent = "admin_topup_created"
	AdminEventTopUpPaid      AdminTelegramEvent = "admin_topup_paid"
	AdminEventTokenCreated   AdminTelegramEvent = "admin_token_created"
)

type AdminTelegramNotifyPayload struct {
	Event   AdminTelegramEvent `json:"event"`
	Details map[string]any     `json:"details"`
}

// NotifyAdminTelegramAsync sends an admin notification to Telegram asynchronously.
// It never returns an error to the caller and never blocks the business flow.
func NotifyAdminTelegramAsync(event AdminTelegramEvent, details map[string]any) {
	if !common.AdminTelegramNotifyEnabled {
		return
	}
	if common.AdminTelegramBotToken == "" || common.AdminTelegramChatID == "" {
		return
	}

	payload := AdminTelegramNotifyPayload{
		Event:   event,
		Details: details,
	}

	gopool.Go(func() {
		sendAdminTelegramNotify(payload)
	})
}

func sendAdminTelegramNotify(payload AdminTelegramNotifyPayload) {
	messageText := formatAdminTelegramMessage(payload)

	type tgSendMessageReq struct {
		ChatID string `json:"chat_id"`
		Text   string `json:"text"`
	}

	reqBody := tgSendMessageReq{
		ChatID: common.AdminTelegramChatID,
		Text:   messageText,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		common.SysLog(fmt.Sprintf("admin telegram notify marshal error: %s", err.Error()))
		return
	}

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", common.AdminTelegramBotToken)

	client := http.Client{Timeout: 3 * time.Second}
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		common.SysLog(fmt.Sprintf("admin telegram notify request error: %s", err.Error()))
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		common.SysLog(fmt.Sprintf("admin telegram notify non-2xx status: %d", resp.StatusCode))
		return
	}

	common.SysLog(fmt.Sprintf("admin telegram notify sent: %s", payload.Event))
}

func formatAdminTelegramMessage(payload AdminTelegramNotifyPayload) string {
	var title string
	switch payload.Event {
	case AdminEventUserRegistered:
		title = "🆕 新用户注册"
	case AdminEventTopUpCreated:
		title = "💰 充值订单创建"
	case AdminEventTopUpPaid:
		title = "✅ 充值付款成功"
	case AdminEventTokenCreated:
		title = "🔑 新 API Key 创建"
	default:
		title = string(payload.Event)
	}

	var detailsText string
	for k, v := range payload.Details {
		detailsText += fmt.Sprintf("\n%s: %v", k, v)
	}

	return title + detailsText
}
