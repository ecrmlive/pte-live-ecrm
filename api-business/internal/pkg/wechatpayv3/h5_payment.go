package wechatpayv3

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type H5Request struct {
	Description   string
	OutTradeNo    string
	AmountCents   int64
	PayerClientIP string
	Attach        string
	ExpireAt      time.Time
}

type H5Response struct {
	MWEBURL string
}

// H5Prepay creates a WeChat H5 (MWEB) payment. The payment URL is transient
// and is returned only to the authenticated order owner.
func (c *Client) H5Prepay(ctx context.Context, config Config, request H5Request) (H5Response, error) {
	if !config.ValidForH5() || request.AmountCents <= 0 || strings.TrimSpace(request.OutTradeNo) == "" || strings.TrimSpace(request.PayerClientIP) == "" {
		return H5Response{}, ErrInvalidConfig
	}
	privateKey, err := parseRSAPrivateKey(config.MerchantPrivateKey)
	if err != nil {
		return H5Response{}, ErrInvalidConfig
	}
	now := c.now()
	expires := request.ExpireAt.UTC()
	if expires.IsZero() || !expires.After(now) {
		expires = now.Add(15 * time.Minute)
	}
	type h5Info struct {
		Type    string `json:"type"`
		AppName string `json:"app_name"`
		AppURL  string `json:"app_url"`
	}
	type sceneInfo struct {
		PayerClientIP string `json:"payer_client_ip"`
		H5Info        h5Info `json:"h5_info"`
	}
	body, err := json.Marshal(struct {
		AppID       string `json:"appid"`
		MchID       string `json:"mchid"`
		Description string `json:"description"`
		OutTradeNo  string `json:"out_trade_no"`
		Attach      string `json:"attach,omitempty"`
		NotifyURL   string `json:"notify_url"`
		TimeExpire  string `json:"time_expire"`
		Amount      struct {
			Total    int64  `json:"total"`
			Currency string `json:"currency"`
		} `json:"amount"`
		SceneInfo sceneInfo `json:"scene_info"`
	}{
		AppID:       strings.TrimSpace(config.H5AppID),
		MchID:       strings.TrimSpace(config.MchID),
		Description: truncate(strings.TrimSpace(request.Description), 127),
		OutTradeNo:  strings.TrimSpace(request.OutTradeNo),
		Attach:      truncate(strings.TrimSpace(request.Attach), 127),
		NotifyURL:   strings.TrimSpace(config.NotifyURL),
		TimeExpire:  expires.Format(time.RFC3339),
		Amount: struct {
			Total    int64  `json:"total"`
			Currency string `json:"currency"`
		}{Total: request.AmountCents, Currency: "CNY"},
		SceneInfo: sceneInfo{PayerClientIP: strings.TrimSpace(request.PayerClientIP), H5Info: h5Info{Type: "Wap", AppName: "CRM Live", AppURL: strings.TrimSpace(config.H5SiteURL)}},
	})
	if err != nil {
		return H5Response{}, err
	}
	const path = "/v3/pay/transactions/h5"
	nonce, err := c.nonce()
	if err != nil {
		return H5Response{}, err
	}
	timestamp := fmt.Sprintf("%d", now.Unix())
	signature, err := sign(privateKey, strings.Join([]string{http.MethodPost, path, timestamp, nonce, string(body), ""}, "\n"))
	if err != nil {
		return H5Response{}, err
	}
	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL()+path, bytes.NewReader(body))
	if err != nil {
		return H5Response{}, err
	}
	httpRequest.Header.Set("Accept", "application/json")
	httpRequest.Header.Set("Content-Type", "application/json")
	httpRequest.Header.Set("Authorization", fmt.Sprintf(`WECHATPAY2-SHA256-RSA2048 mchid="%s",nonce_str="%s",timestamp="%s",serial_no="%s",signature="%s"`, config.MchID, nonce, timestamp, config.resolvedMerchantSerialNo(), signature))
	resp, err := c.httpClient().Do(httpRequest)
	if err != nil {
		return H5Response{}, fmt.Errorf("微信支付下单请求失败: %w", err)
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return H5Response{}, err
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		var failure struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		}
		_ = json.Unmarshal(data, &failure)
		return H5Response{}, fmt.Errorf("微信支付下单失败: %s", providerMessage(failure.Code, failure.Message, resp.StatusCode))
	}
	var output struct {
		MWEBURL string `json:"h5_url"`
	}
	if err := json.Unmarshal(data, &output); err != nil || !strings.HasPrefix(strings.TrimSpace(output.MWEBURL), "https://") {
		return H5Response{}, ErrInvalidCallback
	}
	return H5Response{MWEBURL: output.MWEBURL}, nil
}
