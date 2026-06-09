package main

import (
	"context"
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/skip2/go-qrcode"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/core/downloader"
	"github.com/wechatpay-apiv3/wechatpay-go/core/notify"
	"github.com/wechatpay-apiv3/wechatpay-go/core/option"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/native"
	"github.com/wechatpay-apiv3/wechatpay-go/utils"
)

// WechatClient wraps WeChat Pay APIv3 Native payment
type WechatClient struct {
	client  *core.Client
	cfg     *WechatConfig
	notifyH *notify.Handler
}

func NewWechatClient(cfg *WechatConfig) (*WechatClient, error) {
	if cfg.MchID == "" || cfg.APIv3Key == "" || cfg.SerialNo == "" || cfg.PrivateKey == "" {
		return nil, fmt.Errorf("wechat pay not configured")
	}

	// Parse private key
	pk, err := utils.LoadPrivateKey(cfg.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("parse private key: %w", err)
	}

	// Create client with auto cert download + auth
	client, err := core.NewClient(
		context.Background(),
		option.WithWechatPayAutoAuthCipher(cfg.MchID, cfg.SerialNo, pk, cfg.APIv3Key),
	)
	if err != nil {
		return nil, fmt.Errorf("create wechat pay client: %w", err)
	}

	cfg.Configured = true

	// Setup notify handler (separate cert downloader for callback verification)
	mgr := downloader.MgrInstance()
	if err := mgr.RegisterDownloaderWithClient(context.Background(), client, cfg.MchID, cfg.APIv3Key); err != nil {
		return nil, fmt.Errorf("register cert downloader: %w", err)
	}

	// Create adapter so CertificateDownloaderMgr satisfies CertificateGetter
	certGetter := &certGetterAdapter{mgr: mgr, mchID: cfg.MchID}

	// Create AES-GCM for decrypting callback data
	c, err := aes.NewCipher([]byte(cfg.APIv3Key))
	if err != nil {
		return nil, fmt.Errorf("create aes cipher: %w", err)
	}
	aesgcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, fmt.Errorf("create gcm: %w", err)
	}

	notifyH := notify.NewEmptyHandler().AddRSAWithAESGCM(
		NewWechatPayVerifier(certGetter),
		aesgcm,
	)

	return &WechatClient{client: client, cfg: cfg, notifyH: notifyH}, nil
}

// CreateNativeOrder creates a Native payment order, returns code_url (QR code content)
func (w *WechatClient) CreateNativeOrder(outTradeNo, description, amount string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	svc := native.NativeApiService{Client: w.client}
	req := native.PrepayRequest{
		Appid:       core.String(""),
		Mchid:       core.String(w.cfg.MchID),
		Description: core.String(description),
		OutTradeNo:  core.String(outTradeNo),
		NotifyUrl:   core.String(w.cfg.NotifyURL),
		Amount: &native.Amount{
			Total:    core.Int64(yuanToFen(amount)),
			Currency: core.String("CNY"),
		},
	}

	resp, _, err := svc.Prepay(ctx, req)
	if err != nil {
		return "", fmt.Errorf("prepay: %w", err)
	}
	if resp.CodeUrl == nil || *resp.CodeUrl == "" {
		return "", fmt.Errorf("empty code_url in response")
	}
	return *resp.CodeUrl, nil
}

// WechatTradeResult is the decoded trade info from a WeChat Pay callback
type WechatTradeResult struct {
	OutTradeNo    string `json:"out_trade_no"`
	TradeState    string `json:"trade_state"`
	TransactionID string `json:"transaction_id"`
}

// ParseNotify parses and verifies a WeChat Pay callback
func (w *WechatClient) ParseNotify(req *http.Request) (*WechatTradeResult, error) {
	content := new(WechatTradeResult)
	_, err := w.notifyH.ParseNotifyRequest(context.Background(), req, content)
	if err != nil {
		return nil, fmt.Errorf("parse notify: %w", err)
	}
	if content.OutTradeNo == "" {
		return nil, fmt.Errorf("empty out_trade_no in callback")
	}
	return content, nil
}

// RenderQRPage generates an HTML page showing the WeChat Pay QR code
func RenderQRPage(codeURL, outTradeNo, returnURL string, amount string) string {
	png, err := qrcode.Encode(codeURL, qrcode.Medium, 280)
	if err != nil {
		log.Printf("[wechat] QR encode error: %v", err)
		return `<html><body><p>生成二维码失败，请重试</p></body></html>`
	}
	qrBase64 := base64.StdEncoding.EncodeToString(png)

	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="zh-CN">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<title>微信支付</title>
<style>
*{margin:0;padding:0;box-sizing:border-box}
body{font-family:-apple-system,BlinkMacSystemFont,"Segoe UI",Roboto,sans-serif;
background:#f5f5f5;display:flex;justify-content:center;align-items:center;min-height:100vh}
.card{background:#fff;border-radius:12px;padding:40px 32px;text-align:center;
box-shadow:0 2px 16px rgba(0,0,0,.08);max-width:400px;width:90%%}
.title{font-size:20px;font-weight:600;color:#333;margin-bottom:8px}
.amount{font-size:28px;color:#07C160;font-weight:700;margin-bottom:8px}
.subtitle{font-size:14px;color:#999;margin-bottom:24px}
.qrcode{border:1px solid #eee;border-radius:8px;padding:8px;margin-bottom:16px;display:inline-block}
.qrcode img{display:block;width:240px;height:240px}
.tip{font-size:14px;color:#666;margin-top:16px}
.loading{display:inline-block;width:16px;height:16px;border:2px solid #07C160;
border-top-color:transparent;border-radius:50%%;animation:spin .8s linear infinite;
vertical-align:middle;margin-right:6px}
@keyframes spin{to{transform:rotate(360deg)}}
.footer{margin-top:24px;font-size:12px;color:#bbb}
</style>
</head>
<body>
<div class="card">
<div class="title">微信扫码支付</div>
<div class="amount">¥%s</div>
<div class="subtitle">订单号：%s</div>
<div class="qrcode"><img src="data:image/png;base64,%s" alt="微信支付二维码"></div>
<div class="tip"><span class="loading"></span>请使用微信扫一扫完成支付</div>
<div class="footer">支付完成后将自动跳转</div>
</div>
<script>
(function(){var c=false;setInterval(function(){if(c)return;
fetch('/wechat/check?out_trade_no=%s').then(function(r){return r.json()})
.then(function(d){if(d.paid){c=true;window.location.href='%s'}})},3000)})();
</script>
</body>
</html>`, amount, outTradeNo, qrBase64, outTradeNo, returnURL)
}

// yuanToFen converts yuan string to fen int64
func yuanToFen(yuan string) int64 {
	var f float64
	if _, err := fmt.Sscanf(yuan, "%f", &f); err != nil {
		return 0
	}
	return int64(f*100 + 0.5)
}

// ============================================================================
// Adapter: CertificateDownloaderMgr → CertificateGetter
// ============================================================================

type certGetterAdapter struct {
	mgr   *downloader.CertificateDownloaderMgr
	mchID string
}

func (c *certGetterAdapter) Get(ctx context.Context, serialNumber string) (*x509.Certificate, bool) {
	return c.mgr.GetCertificate(ctx, c.mchID, serialNumber)
}

func (c *certGetterAdapter) GetAll(ctx context.Context) map[string]*x509.Certificate {
	return c.mgr.GetCertificateMap(ctx, c.mchID)
}

func (c *certGetterAdapter) GetNewestSerial(ctx context.Context) string {
	return c.mgr.GetNewestCertificateSerial(ctx, c.mchID)
}

// ============================================================================
// Custom RSA Verifier (minimal implementation for wechat pay callbacks)
// ============================================================================

type wechatPayVerifier struct {
	getter core.CertificateGetter
}

func NewWechatPayVerifier(getter core.CertificateGetter) *wechatPayVerifier {
	return &wechatPayVerifier{getter: getter}
}

func (v *wechatPayVerifier) Verify(ctx context.Context, serialNumber, message, signature string) error {
	cert, ok := v.getter.Get(ctx, serialNumber)
	if !ok {
		return fmt.Errorf("certificate not found for serial: %s", serialNumber)
	}
	pubKey, ok := cert.PublicKey.(*rsa.PublicKey)
	if !ok {
		return fmt.Errorf("certificate public key is not RSA")
	}
	sig, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return fmt.Errorf("decode signature: %w", err)
	}
	hashed := sha256.Sum256([]byte(message))
	return rsa.VerifyPKCS1v15(pubKey, crypto.SHA256, hashed[:], sig)
}

func (v *wechatPayVerifier) GetSerial(ctx context.Context) (string, error) {
	serial := v.getter.GetNewestSerial(ctx)
	if serial == "" {
		return "", fmt.Errorf("no certificate available")
	}
	return serial, nil
}
