package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/smartwalle/alipay/v3"
)

// ============================================================================
// Alipay Integration
// ============================================================================

type AlipayClient struct {
	client *alipay.Client
	cfg    *AlipayConfig
}

func NewAlipayClient(cfg *AlipayConfig) (*AlipayClient, error) {
	client, err := alipay.New(cfg.AppID, cfg.PrivateKey, cfg.IsProduction)
	if err != nil {
		return nil, fmt.Errorf("create alipay client: %w", err)
	}
	if err := client.LoadAliPayPublicKey(cfg.AliPayPublicKey); err != nil {
		return nil, fmt.Errorf("load alipay public key: %w", err)
	}
	return &AlipayClient{client: client, cfg: cfg}, nil
}

// CreatePagePay creates a PC website payment and returns the redirect URL
func (a *AlipayClient) CreatePagePay(subject, outTradeNo, amount string) (string, error) {
	p := alipay.TradePagePay{
		Trade: alipay.Trade{
			Subject:     subject,
			OutTradeNo:  outTradeNo,
			TotalAmount: amount,
			ProductCode: "FAST_INSTANT_TRADE_PAY",
			NotifyURL:   a.cfg.NotifyURL,
			ReturnURL:   a.cfg.ReturnURL,
		},
	}

	payURL, err := a.client.TradePagePay(p)
	if err != nil {
		return "", fmt.Errorf("create page pay: %w", err)
	}
	return payURL.String(), nil
}

// VerifySign verifies the Alipay notification signature
func (a *AlipayClient) VerifySign(ctx context.Context, form url.Values) error {
	return a.client.VerifySign(ctx, form)
}

// DecodeNotification decodes the Alipay async notification from form data
func (a *AlipayClient) DecodeNotification(req *http.Request) (*alipay.Notification, error) {
	return a.client.GetTradeNotification(req)
}

// QueryTrade queries the payment status from Alipay
func (a *AlipayClient) QueryTrade(outTradeNo string) (*alipay.TradeQueryRsp, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	p := alipay.TradeQuery{OutTradeNo: outTradeNo}
	rsp, err := a.client.TradeQuery(ctx, p)
	if err != nil {
		return nil, fmt.Errorf("query trade: %w", err)
	}
	return rsp, nil
}
