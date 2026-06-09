package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// ============================================================================
// Epay Server — HTTP Handlers
// ============================================================================

type EpayServer struct {
	cfg    *Config
	store  *Store
	alipay *AlipayClient
	wechat *WechatClient // nil if not configured
}

func NewEpayServer(cfg *Config, store *Store, alipay *AlipayClient, wechat *WechatClient) *EpayServer {
	return &EpayServer{cfg: cfg, store: store, alipay: alipay, wechat: wechat}
}

func (s *EpayServer) RegisterRoutes(r *gin.Engine) {
	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	// Epay purchase endpoint — receives form POST from fast-api frontend
	r.POST("/submit.php", s.HandleSubmit)
	r.GET("/submit.php", s.HandleSubmit)

	// Alipay return — user browser gets redirected here after payment
	r.GET("/alipay/return", s.HandleAlipayReturn)
	r.POST("/alipay/return", s.HandleAlipayReturn)

	// Alipay async notify — Alipay server POSTs payment result here
	r.POST("/alipay/notify", s.HandleAlipayNotify)

	// WeChat Pay async notify — WeChat server POSTs payment result here
	r.POST("/wechat/notify", s.HandleWechatNotify)

	// WeChat payment status polling (for QR page auto-redirect)
	r.GET("/wechat/check", s.HandleWechatCheck)
}

// ============================================================================
// POST /submit.php — Epay Purchase (routes to Alipay or WeChat)
// ============================================================================

func (s *EpayServer) HandleSubmit(c *gin.Context) {
	if c.Request.Method == "POST" {
		if err := c.Request.ParseForm(); err != nil {
			log.Printf("[submit] parse form error: %v", err)
			c.String(http.StatusBadRequest, "参数错误")
			return
		}
	}

	params := make(map[string]string)
	if c.Request.Method == "POST" {
		for k, v := range c.Request.PostForm {
			if len(v) > 0 {
				params[k] = v[0]
			}
		}
	} else {
		for k, v := range c.Request.URL.Query() {
			if len(v) > 0 {
				params[k] = v[0]
			}
		}
	}

	log.Printf("[submit] received params: pid=%s type=%s out_trade_no=%s money=%s name=%s",
		params["pid"], params["type"], params["out_trade_no"], params["money"], params["name"])

	// 1. Verify epay signature
	if !verifySign(params, s.cfg.Epay.Key) {
		log.Printf("[submit] sign verification FAILED for out_trade_no=%s", params["out_trade_no"])
		c.String(http.StatusBadRequest, "签名验证失败")
		return
	}

	// 2. Check required fields
	outTradeNo := params["out_trade_no"]
	paymentType := params["type"]
	amount := params["money"]
	productName := params["name"]
	notifyURL := params["notify_url"]
	returnURL := params["return_url"]
	merchantID := params["pid"]

	if outTradeNo == "" || amount == "" || notifyURL == "" {
		log.Printf("[submit] missing required fields: out_trade_no=%s amount=%s notify_url=%s",
			outTradeNo, amount, notifyURL)
		c.String(http.StatusBadRequest, "缺少必要参数")
		return
	}

	// 3. Check duplicate
	if existing, err := s.store.GetOrderByEpayNo(outTradeNo); err == nil && existing != nil {
		log.Printf("[submit] duplicate order out_trade_no=%s existing_status=%s", outTradeNo, existing.Status)
		if existing.Status == OrderStatusPaid {
			if returnURL != "" {
				c.Redirect(http.StatusFound, returnURL)
				return
			}
		}
		c.String(http.StatusBadRequest, "订单已存在")
		return
	}

	// 4. Create local order
	order := &Order{
		EpayOrderNo: outTradeNo,
		PaymentType: paymentType,
		ProductName: productName,
		Amount:      amount,
		Status:      OrderStatusPending,
		NotifyURL:   notifyURL,
		ReturnURL:   returnURL,
		MerchantID:  merchantID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if err := s.store.CreateOrder(order); err != nil {
		log.Printf("[submit] create order error: %v out_trade_no=%s", err, outTradeNo)
		c.String(http.StatusInternalServerError, "创建订单失败")
		return
	}

	log.Printf("[submit] order created out_trade_no=%s amount=%s type=%s", outTradeNo, amount, paymentType)

	// 5. Route to payment provider
	subject := productName
	if subject == "" {
		subject = fmt.Sprintf("充值订单 %s", outTradeNo)
	}

	switch paymentType {
	case "wxpay":
		s.handleWechatPay(c, subject, outTradeNo, amount, returnURL)
	default:
		s.handleAlipayPay(c, subject, outTradeNo, amount)
	}
}

func (s *EpayServer) handleAlipayPay(c *gin.Context, subject, outTradeNo, amount string) {
	payURL, err := s.alipay.CreatePagePay(subject, outTradeNo, amount)
	if err != nil {
		log.Printf("[submit] create alipay payment error: %v out_trade_no=%s", err, outTradeNo)
		c.String(http.StatusInternalServerError, "创建支付宝支付失败")
		return
	}
	log.Printf("[submit] redirecting to alipay out_trade_no=%s", outTradeNo)
	c.Redirect(http.StatusFound, payURL)
}

func (s *EpayServer) handleWechatPay(c *gin.Context, subject, outTradeNo, amount, returnURL string) {
	if s.wechat == nil || !s.cfg.Wechat.Configured {
		log.Printf("[wechat] wechat pay not configured, rejecting out_trade_no=%s", outTradeNo)
		c.String(http.StatusBadRequest, "微信支付未配置")
		return
	}

	codeURL, err := s.wechat.CreateNativeOrder(outTradeNo, subject, amount)
	if err != nil {
		log.Printf("[wechat] create native order error: %v out_trade_no=%s", err, outTradeNo)
		c.String(http.StatusInternalServerError, "创建微信支付订单失败")
		return
	}

	log.Printf("[wechat] native order created out_trade_no=%s code_url=%s", outTradeNo, codeURL)
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, RenderQRPage(codeURL, outTradeNo, returnURL, amount))
}

// ============================================================================
// GET/POST /alipay/return — Alipay Sync Return
// ============================================================================

func (s *EpayServer) HandleAlipayReturn(c *gin.Context) {
	if c.Request.Method == "POST" {
		if err := c.Request.ParseForm(); err != nil {
			c.String(http.StatusBadRequest, "参数错误")
			return
		}
	}

	var queryValues url.Values
	if c.Request.Method == "POST" {
		queryValues = c.Request.PostForm
	} else {
		queryValues = c.Request.URL.Query()
	}

	outTradeNo := queryValues.Get("out_trade_no")
	log.Printf("[return] received out_trade_no=%s", outTradeNo)

	if err := s.alipay.VerifySign(c.Request.Context(), queryValues); err != nil {
		log.Printf("[return] verify sign failed out_trade_no=%s err=%v", outTradeNo, err)
		c.String(http.StatusBadRequest, "签名验证失败")
		return
	}

	rsp, err := s.alipay.QueryTrade(outTradeNo)
	if err != nil {
		log.Printf("[return] query trade error: %v out_trade_no=%s", err, outTradeNo)
	} else if rsp != nil && rsp.TradeStatus == "TRADE_SUCCESS" {
		s.forwardEpayNotify(outTradeNo, "")
	}

	order, err := s.store.GetOrderByEpayNo(outTradeNo)
	if err != nil {
		c.String(http.StatusBadRequest, "订单不存在")
		return
	}

	if order.ReturnURL != "" {
		c.Redirect(http.StatusFound, order.ReturnURL)
		return
	}
	c.String(http.StatusOK, "支付成功，请返回充值页面查看余额。")
}

// ============================================================================
// POST /alipay/notify — Alipay Async Notify
// ============================================================================

func (s *EpayServer) HandleAlipayNotify(c *gin.Context) {
	if err := c.Request.ParseForm(); err != nil {
		log.Printf("[notify] parse form error: %v", err)
		c.String(http.StatusBadRequest, "fail")
		return
	}

	outTradeNo := c.Request.PostForm.Get("out_trade_no")
	tradeNo := c.Request.PostForm.Get("trade_no")
	tradeStatus := c.Request.PostForm.Get("trade_status")
	log.Printf("[notify] received out_trade_no=%s trade_no=%s trade_status=%s",
		outTradeNo, tradeNo, tradeStatus)

	if err := s.alipay.VerifySign(c.Request.Context(), c.Request.PostForm); err != nil {
		log.Printf("[notify] verify sign FAILED out_trade_no=%s err=%v", outTradeNo, err)
		c.String(http.StatusBadRequest, "fail")
		return
	}

	if tradeStatus != "TRADE_SUCCESS" {
		log.Printf("[notify] non-success status out_trade_no=%s status=%s", outTradeNo, tradeStatus)
		c.String(http.StatusOK, "success")
		return
	}

	s.forwardEpayNotify(outTradeNo, tradeNo)
	c.String(http.StatusOK, "success")
}

// ============================================================================
// POST /wechat/notify — WeChat Pay Async Notify
// ============================================================================

func (s *EpayServer) HandleWechatNotify(c *gin.Context) {
	if s.wechat == nil {
		log.Printf("[wechat-notify] wechat not configured")
		c.JSON(http.StatusInternalServerError, gin.H{"code": "FAIL", "message": "not configured"})
		return
	}

	txn, err := s.wechat.ParseNotify(c.Request)
	if err != nil {
		log.Printf("[wechat-notify] parse error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"code": "FAIL", "message": err.Error()})
		return
	}

	log.Printf("[wechat-notify] received out_trade_no=%s trade_state=%s txn_id=%s",
		txn.OutTradeNo, txn.TradeState, txn.TransactionID)

	// Only process successful trades
	if txn.TradeState != "SUCCESS" {
		log.Printf("[wechat-notify] non-success state out_trade_no=%s state=%s", txn.OutTradeNo, txn.TradeState)
		c.JSON(http.StatusOK, gin.H{"code": "SUCCESS", "message": "OK"})
		return
	}

	s.forwardEpayNotify(txn.OutTradeNo, txn.TransactionID)
	c.JSON(http.StatusOK, gin.H{"code": "SUCCESS", "message": "OK"})
}

// ============================================================================
// GET /wechat/check — Payment status polling for QR page
// ============================================================================

func (s *EpayServer) HandleWechatCheck(c *gin.Context) {
	outTradeNo := c.Query("out_trade_no")
	if outTradeNo == "" {
		c.JSON(http.StatusBadRequest, gin.H{"paid": false})
		return
	}

	order, err := s.store.GetOrderByEpayNo(outTradeNo)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"paid": false})
		return
	}

	c.JSON(http.StatusOK, gin.H{"paid": order.Status == OrderStatusPaid})
}

// ============================================================================
// Forward Epay Notify to fast-api
// ============================================================================

func (s *EpayServer) forwardEpayNotify(outTradeNo, providerTradeNo string) {
	if err := s.store.UpdateOrderPaid(outTradeNo, providerTradeNo); err != nil {
		order, err2 := s.store.GetOrderByEpayNo(outTradeNo)
		if err2 != nil {
			log.Printf("[notify-forward] order not found out_trade_no=%s", outTradeNo)
			return
		}
		if order.Status == OrderStatusPaid {
			log.Printf("[notify-forward] order already paid out_trade_no=%s", outTradeNo)
			return
		}
		log.Printf("[notify-forward] update order error: %v out_trade_no=%s", err, outTradeNo)
		return
	}

	order, err := s.store.GetOrderByEpayNo(outTradeNo)
	if err != nil {
		log.Printf("[notify-forward] re-read order error: %v out_trade_no=%s", err, outTradeNo)
		return
	}

	notifyParams := map[string]string{
		"pid":          s.cfg.Epay.PartnerID,
		"trade_no":     outTradeNo,
		"out_trade_no": outTradeNo,
		"type":         order.PaymentType,
		"name":         order.ProductName,
		"money":        order.Amount,
		"trade_status": "TRADE_SUCCESS",
	}

	notifyParams = generateParams(notifyParams, s.cfg.Epay.Key)

	log.Printf("[notify-forward] posting to %s type=%s", order.NotifyURL, order.PaymentType)
	if err := postNotify(order.NotifyURL, notifyParams); err != nil {
		log.Printf("[notify-forward] POST FAILED out_trade_no=%s url=%s error=%v",
			outTradeNo, order.NotifyURL, err)
		return
	}
	log.Printf("[notify-forward] SUCCESS out_trade_no=%s", outTradeNo)
}

// postNotify sends an HTTP POST with form-encoded params to the notify URL
func postNotify(notifyURL string, params map[string]string) error {
	form := url.Values{}
	for k, v := range params {
		form.Set(k, v)
	}

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Post(notifyURL, "application/x-www-form-urlencoded",
		strings.NewReader(form.Encode()))
	if err != nil {
		return fmt.Errorf("http post: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
	bodyStr := strings.TrimSpace(string(body))
	if bodyStr != "success" {
		return fmt.Errorf("unexpected response: status=%d body=%s", resp.StatusCode, bodyStr)
	}
	return nil
}
