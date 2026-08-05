package wechatpayv3

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// RefundRequest is constructed exclusively from locked server-side payment
// and refund rows. Amounts are integer cents; the client never supplies them.
type RefundRequest struct {
	OutTradeNo  string
	OutRefundNo string
	Reason      string
	TotalCents  int64
	RefundCents int64
}

type RefundResponse struct {
	OutRefundNo      string
	ProviderRefundNo string
	Status           string
}

// Refund starts a WeChat Pay v3 domestic refund. A successful HTTP response
// means only that WeChat accepted the request; business code must still wait
// for a verified REFUND.SUCCESS callback before crediting the user or marking
// the after-sale request refunded.
func (c *Client) Refund(ctx context.Context, config Config, request RefundRequest) (RefundResponse, error) {
	if !config.ValidForRefund() || strings.TrimSpace(request.OutTradeNo) == "" || strings.TrimSpace(request.OutRefundNo) == "" || request.TotalCents <= 0 || request.RefundCents <= 0 || request.RefundCents > request.TotalCents {
		return RefundResponse{}, ErrInvalidConfig
	}
	privateKey, err := parseRSAPrivateKey(config.MerchantPrivateKey)
	if err != nil {
		return RefundResponse{}, ErrInvalidConfig
	}
	body, err := json.Marshal(struct {
		OutTradeNo  string `json:"out_trade_no"`
		OutRefundNo string `json:"out_refund_no"`
		Reason      string `json:"reason,omitempty"`
		NotifyURL   string `json:"notify_url"`
		Amount      struct {
			Refund   int64  `json:"refund"`
			Total    int64  `json:"total"`
			Currency string `json:"currency"`
		} `json:"amount"`
	}{
		OutTradeNo:  strings.TrimSpace(request.OutTradeNo),
		OutRefundNo: strings.TrimSpace(request.OutRefundNo),
		Reason:      truncate(strings.TrimSpace(request.Reason), 80),
		NotifyURL:   strings.TrimSpace(config.RefundNotifyURL),
		Amount: struct {
			Refund   int64  `json:"refund"`
			Total    int64  `json:"total"`
			Currency string `json:"currency"`
		}{Refund: request.RefundCents, Total: request.TotalCents, Currency: "CNY"},
	})
	if err != nil {
		return RefundResponse{}, err
	}
	const path = "/v3/refund/domestic/refunds"
	nonce, err := c.nonce()
	if err != nil {
		return RefundResponse{}, err
	}
	timestamp := fmt.Sprintf("%d", c.now().Unix())
	signature, err := sign(privateKey, strings.Join([]string{http.MethodPost, path, timestamp, nonce, string(body), ""}, "\n"))
	if err != nil {
		return RefundResponse{}, err
	}
	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL()+path, bytes.NewReader(body))
	if err != nil {
		return RefundResponse{}, err
	}
	httpRequest.Header.Set("Accept", "application/json")
	httpRequest.Header.Set("Content-Type", "application/json")
	httpRequest.Header.Set("Authorization", fmt.Sprintf(`WECHATPAY2-SHA256-RSA2048 mchid="%s",nonce_str="%s",timestamp="%s",serial_no="%s",signature="%s"`, config.MchID, nonce, timestamp, config.resolvedMerchantSerialNo(), signature))
	resp, err := c.httpClient().Do(httpRequest)
	if err != nil {
		return RefundResponse{}, fmt.Errorf("微信退款请求失败: %w", err)
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return RefundResponse{}, err
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		var failure struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		}
		_ = json.Unmarshal(data, &failure)
		return RefundResponse{}, fmt.Errorf("微信退款受理失败: %s", providerMessage(failure.Code, failure.Message, resp.StatusCode))
	}
	var output struct {
		OutRefundNo string `json:"out_refund_no"`
		RefundID    string `json:"refund_id"`
		Status      string `json:"status"`
	}
	if err := json.Unmarshal(data, &output); err != nil || strings.TrimSpace(output.OutRefundNo) != strings.TrimSpace(request.OutRefundNo) || strings.TrimSpace(output.Status) == "" {
		return RefundResponse{}, ErrInvalidCallback
	}
	return RefundResponse{OutRefundNo: output.OutRefundNo, ProviderRefundNo: output.RefundID, Status: output.Status}, nil
}
