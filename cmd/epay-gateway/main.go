package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/gin-gonic/gin"
)

func main() {
	configPath := flag.String("config", "config.yaml", "path to config file")
	dbPath := flag.String("db", "epay_gateway.db", "path to SQLite database file")
	flag.Parse()

	// Load config
	cfg, err := LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	log.Printf("Config loaded: server=%s, epay_pid=%s, alipay_appid=%s, production=%v, wechat=%v",
		cfg.ListenAddr(), cfg.Epay.PartnerID, maskString(cfg.Alipay.AppID),
		cfg.Alipay.IsProduction, isWechatConfigured(&cfg.Wechat))

	// Init SQLite store
	store, err := NewStore(*dbPath)
	if err != nil {
		log.Fatalf("Failed to init database: %v", err)
	}
	log.Printf("Database initialized: %s", *dbPath)

	// Init Alipay client
	alipayClient, err := NewAlipayClient(&cfg.Alipay)
	if err != nil {
		log.Fatalf("Failed to init Alipay client: %v", err)
	}
	log.Printf("Alipay client initialized (production=%v)", cfg.Alipay.IsProduction)

	// Init WeChat client (optional — skips gracefully if not configured)
	var wechatClient *WechatClient
	if isWechatConfigured(&cfg.Wechat) {
		wc, err := NewWechatClient(&cfg.Wechat)
		if err != nil {
			log.Printf("WARNING: WeChat Pay init failed (will be disabled): %v", err)
		} else {
			wechatClient = wc
			log.Printf("WeChat Pay client initialized (mch_id=%s)", maskString(cfg.Wechat.MchID))
		}
	} else {
		log.Printf("WeChat Pay not configured — wxpay requests will be rejected")
	}

	// Setup Gin router
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	// Register routes
	server := NewEpayServer(cfg, store, alipayClient, wechatClient)
	server.RegisterRoutes(router)

	// Graceful shutdown
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		sig := <-sigCh
		log.Printf("Received signal %v, shutting down...", sig)
		os.Exit(0)
	}()

	log.Printf("Epay Gateway starting on %s", cfg.ListenAddr())
	if err := router.Run(cfg.ListenAddr()); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func isWechatConfigured(cfg *WechatConfig) bool {
	return strings.TrimSpace(cfg.MchID) != "" &&
		strings.TrimSpace(cfg.APIv3Key) != "" &&
		strings.TrimSpace(cfg.SerialNo) != "" &&
		strings.TrimSpace(cfg.PrivateKey) != ""
}

func maskString(s string) string {
	if len(s) <= 6 {
		return "****"
	}
	return s[:3] + "****" + s[len(s)-3:]
}
