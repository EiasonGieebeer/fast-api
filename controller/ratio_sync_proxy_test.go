package controller

import (
	"net/http"
	"testing"
)

// TestBuildRatioSyncHTTPClient_NoProxy 无代理时应返回带 DialContext 回退逻辑的默认客户端。
func TestBuildRatioSyncHTTPClient_NoProxy(t *testing.T) {
	for _, p := range []string{"", "   "} {
		client := buildRatioSyncHTTPClient(p)
		if client == nil {
			t.Fatalf("expected non-nil client for empty proxy %q", p)
		}
		tr, ok := client.Transport.(*http.Transport)
		if !ok {
			t.Fatalf("expected *http.Transport for proxy %q", p)
		}
		if tr.Proxy != nil {
			t.Errorf("default client should not set Transport.Proxy for proxy %q", p)
		}
		if tr.DialContext == nil {
			t.Errorf("default client should set custom DialContext for github.io fallback, proxy %q", p)
		}
	}
}

// TestBuildRatioSyncHTTPClient_HTTPProxy http/https 代理应通过 Transport.Proxy 生效。
func TestBuildRatioSyncHTTPClient_HTTPProxy(t *testing.T) {
	for _, p := range []string{"http://127.0.0.1:7890", "https://127.0.0.1:7890"} {
		client := buildRatioSyncHTTPClient(p)
		if client == nil {
			t.Fatalf("expected non-nil client for proxy %q", p)
		}
		tr := client.Transport.(*http.Transport)
		if tr.Proxy == nil {
			t.Errorf("http proxy %q should set Transport.Proxy", p)
		}
		req, _ := http.NewRequest(http.MethodGet, "https://easyrouter.io/api/pricing", nil)
		u, err := tr.Proxy(req)
		if err != nil {
			t.Errorf("proxy func returned error for %q: %v", p, err)
		}
		if u == nil || u.Host != "127.0.0.1:7890" {
			t.Errorf("proxy %q resolved to unexpected target: %v", p, u)
		}
	}
}

// TestBuildRatioSyncHTTPClient_Socks5Proxy socks5 代理应通过自定义 DialContext 生效。
func TestBuildRatioSyncHTTPClient_Socks5Proxy(t *testing.T) {
	for _, p := range []string{
		"socks5://172.17.0.1:7891",
		"socks5h://172.17.0.1:7891",
		"socks5://user:pass@172.17.0.1:7891",
	} {
		client := buildRatioSyncHTTPClient(p)
		if client == nil {
			t.Fatalf("expected non-nil client for proxy %q", p)
		}
		tr := client.Transport.(*http.Transport)
		if tr.DialContext == nil {
			t.Errorf("socks5 proxy %q should set custom DialContext", p)
		}
		if tr.Proxy != nil {
			t.Errorf("socks5 proxy %q should not use Transport.Proxy", p)
		}
	}
}

// TestBuildRatioSyncHTTPClient_Invalid 非法/不支持的代理应返回 nil，由调用方回退。
func TestBuildRatioSyncHTTPClient_Invalid(t *testing.T) {
	for _, p := range []string{"ftp://127.0.0.1:21", "://bad", "not a url at all\x7f"} {
		if client := buildRatioSyncHTTPClient(p); client != nil {
			t.Errorf("expected nil client for invalid proxy %q, got %v", p, client)
		}
	}
}
