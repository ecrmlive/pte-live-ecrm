package wechatpayv3

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

// CallbackRefund is the verified, decrypted subset of a WeChat v3 refund
// notification. It is intentionally separate from CallbackTransaction: a
// REFUND.SUCCESS notification must never be mistaken for payment success.
type CallbackRefund struct {
	EventID       string
	OutTradeNo    string
	OutRefundNo   string
	TransactionID string
	RefundID      string
	MchID         string
	Status        string
	RefundCents   int64
	TotalCents    int64
	SuccessTime   time.Time
	RawPayload    json.RawMessage
}

func VerifyAndDecryptRefundCallback(config Config, header http.Header, body []byte, now time.Time) (CallbackRefund, error) {
	if !config.ValidForCallback() || len(body) == 0 || len(body) > 1<<20 {
		return CallbackRefund{}, ErrInvalidConfig
	}
	if err := VerifyCallbackSignature(config, header, body, now); err != nil {
		return CallbackRefund{}, err
	}
	var envelope struct {
		ID        string `json:"id"`
		EventType string `json:"event_type"`
		Resource  struct {
			Algorithm      string `json:"algorithm"`
			Ciphertext     string `json:"ciphertext"`
			AssociatedData string `json:"associated_data"`
			Nonce          string `json:"nonce"`
		} `json:"resource"`
	}
	if err := json.Unmarshal(body, &envelope); err != nil || envelope.ID == "" || envelope.EventType != "REFUND.SUCCESS" || envelope.Resource.Algorithm != "AEAD_AES_256_GCM" {
		return CallbackRefund{}, ErrInvalidCallback
	}
	plain, err := decryptResource(config.APIv3Key, envelope.Resource.AssociatedData, envelope.Resource.Nonce, envelope.Resource.Ciphertext)
	if err != nil {
		return CallbackRefund{}, ErrInvalidCallback
	}
	var refund struct {
		OutTradeNo    string `json:"out_trade_no"`
		OutRefundNo   string `json:"out_refund_no"`
		TransactionID string `json:"transaction_id"`
		RefundID      string `json:"refund_id"`
		MchID         string `json:"mchid"`
		Status        string `json:"refund_status"`
		SuccessTime   string `json:"success_time"`
		Amount        struct {
			Refund   int64  `json:"refund"`
			Total    int64  `json:"total"`
			Currency string `json:"currency"`
		} `json:"amount"`
	}
	if err := json.Unmarshal(plain, &refund); err != nil || strings.TrimSpace(refund.OutTradeNo) == "" || strings.TrimSpace(refund.OutRefundNo) == "" || strings.TrimSpace(refund.TransactionID) == "" || strings.TrimSpace(refund.RefundID) == "" || strings.TrimSpace(refund.MchID) == "" || refund.Status != "SUCCESS" || refund.Amount.Refund <= 0 || refund.Amount.Total < refund.Amount.Refund || refund.Amount.Currency != "CNY" {
		return CallbackRefund{}, ErrInvalidCallback
	}
	successTime, err := time.Parse(time.RFC3339, refund.SuccessTime)
	if err != nil {
		return CallbackRefund{}, ErrInvalidCallback
	}
	return CallbackRefund{EventID: envelope.ID, OutTradeNo: refund.OutTradeNo, OutRefundNo: refund.OutRefundNo, TransactionID: refund.TransactionID, RefundID: refund.RefundID, MchID: refund.MchID, Status: refund.Status, RefundCents: refund.Amount.Refund, TotalCents: refund.Amount.Total, SuccessTime: successTime.UTC(), RawPayload: append(json.RawMessage(nil), plain...)}, nil
}
