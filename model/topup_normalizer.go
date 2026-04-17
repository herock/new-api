package model

import (
	"github.com/QuantumNous/new-api/common"
	"github.com/shopspring/decimal"
)

// NormalizedTopUp 统一语义的充值记录结构
// 所有读取方应优先使用归一化后的字段，而不是原始的 amount/money
type NormalizedTopUp struct {
	Id   int    `json:"id"`
	UserId int   `json:"user_id"`
	TradeNo string `json:"trade_no"`
	Status string `json:"status"`
	PaymentMethod string `json:"payment_method"`
	PaymentMethodDisplay string `json:"payment_method_display"` // 统一展示名
	OrderType string `json:"order_type"` // topup / subscription

	RawAmount int64 `json:"raw_amount"` // 原始 amount，仅用于调试
	RawMoney float64 `json:"raw_money"` // 原始 money，仅用于调试

	PaidAmountUSD float64 `json:"paid_amount_usd"` // 用户实际支付金额（USD）
	CreditedQuota int64 `json:"credited_quota"` // 实际入账 quota
	CreditedAmountUSD float64 `json:"credited_amount_usd"` // 入账美元值

	CreateTime int64 `json:"create_time"`
	CompleteTime int64 `json:"complete_time"`
}

// NormalizeTopUp 将原始 TopUp 转换为统一语义结构
func NormalizeTopUp(topUp *TopUp) *NormalizedTopUp {
	if topUp == nil {
		return nil
	}

	result := &NormalizedTopUp{
		Id:         topUp.Id,
		UserId:     topUp.UserId,
		TradeNo:    topUp.TradeNo,
		Status:     topUp.Status,
		PaymentMethod: topUp.PaymentMethod,
		PaymentMethodDisplay: NormalizePaymentMethodDisplay(topUp.PaymentMethod),
		OrderType:  DetectTopUpOrderType(topUp),
		RawAmount:  topUp.Amount,
		RawMoney:   topUp.Money,
		CreateTime: topUp.CreateTime,
		CompleteTime: topUp.CompleteTime,
	}

	// 计算支付金额
	// 本阶段统一把 Money 解释为"支付金额字段"
	result.PaidAmountUSD = topUp.Money

	// 计算入账额度
	// 根据支付方式的不同历史语义进行归一化
	result.CreditedQuota = ComputeCreditedQuota(topUp)

	// 计算入账美元值
	if result.CreditedQuota > 0 {
		dQuota := decimal.NewFromInt(result.CreditedQuota)
		dQuotaPerUnit := decimal.NewFromFloat(common.QuotaPerUnit)
		result.CreditedAmountUSD = dQuota.Div(dQuotaPerUnit).InexactFloat64()
	} else {
		result.CreditedAmountUSD = 0
	}

	return result
}

// NormalizeTopUps 批量归一化
func NormalizeTopUps(topUps []*TopUp) []*NormalizedTopUp {
	if topUps == nil {
		return nil
	}
	result := make([]*NormalizedTopUp, 0, len(topUps))
	for _, tu := range topUps {
		if tu != nil {
			result = append(result, NormalizeTopUp(tu))
		}
	}
	return result
}

// ComputeCreditedQuota 根据支付方式计算实际入账额度
//
// 各支付方式的归一规则：
// - Stripe: Money * QuotaPerUnit (Money 是经分组倍率换算后的美元数量)
// - Creem: Amount (Amount 直接存 quota 整数)
// - Waffo: Amount * QuotaPerUnit
// - Epay/USDC/其他在线充值: Amount * QuotaPerUnit
// - Subscription bridge: Amount=0 时返回 0
func ComputeCreditedQuota(topUp *TopUp) int64 {
	if topUp == nil {
		return 0
	}

	// 订阅桥接记录特殊处理
	if topUp.Amount == 0 && DetectTopUpOrderType(topUp) == "subscription" {
		return 0
	}

	var creditedQuota int64

	switch topUp.PaymentMethod {
	case "stripe":
		// Stripe: Money 代表经分组倍率换算后的美元数量
		dMoney := decimal.NewFromFloat(topUp.Money)
		dQuotaPerUnit := decimal.NewFromFloat(common.QuotaPerUnit)
		creditedQuota = dMoney.Mul(dQuotaPerUnit).IntPart()

	case "creem":
		// Creem: Amount 直接是 quota 整数
		creditedQuota = topUp.Amount

	default:
		// 其他支付方式（Waffo, Epay, USDC 等）: Amount * QuotaPerUnit
		dAmount := decimal.NewFromInt(topUp.Amount)
		dQuotaPerUnit := decimal.NewFromFloat(common.QuotaPerUnit)
		creditedQuota = dAmount.Mul(dQuotaPerUnit).IntPart()
	}

	return creditedQuota
}

// NormalizePaymentMethodDisplay 统一支付方式展示名称
func NormalizePaymentMethodDisplay(method string) string {
	switch method {
	case "stripe":
		return "Stripe"
	case "creem":
		return "Creem"
	case "waffo":
		return "Waffo"
	case "sol_usdc":
		return "SOL USDC"
	case "trx_usdt":
		return "TRX USDT"
	case "alipay":
		return "支付宝"
	case "wxpay":
		return "微信"
	case "":
		return "-"
	default:
		// 未知支付方式，返回原始值的首字母大写形式
		if len(method) > 0 {
			return method
		}
		return "-"
	}
}

// DetectTopUpOrderType 检测订单类型
// 返回 "topup" 或 "subscription"
func DetectTopUpOrderType(topUp *TopUp) string {
	if topUp == nil {
		return "topup"
	}

	tradeNo := topUp.TradeNo

	// 订阅订单特征：Amount=0 且 TradeNo 以 sub 开头（不区分大小写）
	if topUp.Amount == 0 && len(tradeNo) >= 3 {
		prefix := tradeNo[:3]
		if prefix == "sub" || prefix == "Sub" || prefix == "SUB" {
			return "subscription"
		}
	}

	return "topup"
}
