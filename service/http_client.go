package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/QuantumNous/fast-api/common"
	"github.com/QuantumNous/fast-api/setting/system_setting"
)

var (
	httpClient      *http.Client
	proxyClientLock sync.Mutex
	proxyClients    = make(map[string]*http.Client)
)

func checkRedirect(req *http.Request, via []*http.Request) error {
	fetchSetting := system_setting.GetFetchSetting()
	urlStr := req.URL.String()
	if err := common.ValidateURLWithFetchSetting(urlStr, fetchSetting.EnableSSRFProtection, fetchSetting.AllowPrivateIp, fetchSetting.DomainFilterMode, fetchSetting.IpFilterMode, fetchSetting.DomainList, fetchSetting.IpList, fetchSetting.AllowedPorts, fetchSetting.ApplyIPFilterForDomain); err != nil {
		return fmt.Errorf("redirect to %s blocked: %v", urlStr, err)
	}
	if len(via) >= 10 {
		return fmt.Errorf("stopped after 10 redirects")
	}
	return nil
}

func InitHttpClient() {
	transport := &http.Transport{
		MaxIdleConns:        common.RelayMaxIdleConns,
		MaxIdleConnsPerHost: common.RelayMaxIdleConnsPerHost,
		ForceAttemptHTTP2:   true,
		Proxy:               http.ProxyFromEnvironment, // Support HTTP_PROXY, HTTPS_PROXY, NO_PROXY env vars
	}
	if common.TLSInsecureSkipVerify {
		transport.TLSClientConfig = common.InsecureTLSConfig
	}

	if common.RelayTimeout == 0 {
		httpClient = &http.Client{
			Transport:     transport,
			CheckRedirect: checkRedirect,
		}
	} else {
		httpClient = &http.Client{
			Transport:     transport,
			Timeout:       time.Duration(common.RelayTimeout) * time.Second,
			CheckRedirect: checkRedirect,
		}
	}
}

func GetHttpClient() *http.Client {
	return httpClient
}

// GetHttpClientWithProxy returns the default client or a proxy-enabled one when proxyURL is provided.
func GetHttpClientWithProxy(proxyURL string) (*http.Client, error) {
	if proxyURL == "" {
		return GetHttpClient(), nil
	}
	return NewProxyHttpClient(proxyURL)
}

// ResetProxyClientCache 清空代理客户端缓存，确保下次使用时重新初始化
func ResetProxyClientCache() {
	proxyClientLock.Lock()
	defer proxyClientLock.Unlock()
	for _, client := range proxyClients {
		if transport, ok := client.Transport.(*http.Transport); ok && transport != nil {
			transport.CloseIdleConnections()
		}
	}
	proxyClients = make(map[string]*http.Client)
}

// NewProxyHttpClient 创建支持代理的 HTTP 客户端
func NewProxyHttpClient(proxyURL string) (*http.Client, error) {
	if proxyURL == "" {
		if client := GetHttpClient(); client != nil {
			return client, nil
		}
		return http.DefaultClient, nil
	}

	proxyClientLock.Lock()
	if client, ok := proxyClients[proxyURL]; ok {
		proxyClientLock.Unlock()
		return client, nil
	}
	proxyClientLock.Unlock()

	parsedURL, err := url.Parse(proxyURL)
	if err != nil {
		return nil, err
	}

	switch parsedURL.Scheme {
	case "http", "https":
		transport := &http.Transport{
			MaxIdleConns:        common.RelayMaxIdleConns,
			MaxIdleConnsPerHost: common.RelayMaxIdleConnsPerHost,
			ForceAttemptHTTP2:   true,
			Proxy:               http.ProxyURL(parsedURL),
		}
		if common.TLSInsecureSkipVerify {
			transport.TLSClientConfig = common.InsecureTLSConfig
		}
		client := &http.Client{
			Transport:     transport,
			CheckRedirect: checkRedirect,
		}
		client.Timeout = time.Duration(common.RelayTimeout) * time.Second
		proxyClientLock.Lock()
		proxyClients[proxyURL] = client
		proxyClientLock.Unlock()
		return client, nil

	case "socks5", "socks5h":
		// Build a SOCKS5 dialer that sends domain names through the proxy (SOCKS5h behavior).
		// We do NOT use golang.org/x/net/proxy.SOCKS5 because it resolves DNS locally,
		// which breaks connectivity in environments with DNS pollution (e.g. China GFW).
		dialer, err := newSOCKS5hDialer(parsedURL.Host, parsedURL.User)
		if err != nil {
			return nil, err
		}

		transport := &http.Transport{
			MaxIdleConns:        common.RelayMaxIdleConns,
			MaxIdleConnsPerHost: common.RelayMaxIdleConnsPerHost,
			ForceAttemptHTTP2:   true,
			DialContext:         dialer,
		}
		if common.TLSInsecureSkipVerify {
			transport.TLSClientConfig = common.InsecureTLSConfig
		}

		client := &http.Client{Transport: transport, CheckRedirect: checkRedirect}
		client.Timeout = time.Duration(common.RelayTimeout) * time.Second
		proxyClientLock.Lock()
		proxyClients[proxyURL] = client
		proxyClientLock.Unlock()
		return client, nil

	default:
		return nil, fmt.Errorf("unsupported proxy scheme: %s, must be http, https, socks5 or socks5h", parsedURL.Scheme)
	}
}

// newSOCKS5hDialer creates a DialContext function that connects through a SOCKS5 proxy
// and sends domain names through the proxy for DNS resolution (SOCKS5h behavior).
// This avoids local DNS resolution which is critical in environments with DNS pollution.
func newSOCKS5hDialer(proxyAddr string, userinfo *url.Userinfo) (func(ctx context.Context, network, addr string) (net.Conn, error), error) {
	var username, password string
	if userinfo != nil {
		username = userinfo.Username()
		password, _ = userinfo.Password()
	}

	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		if network != "tcp" && network != "tcp4" && network != "tcp6" {
			return nil, fmt.Errorf("socks5h: unsupported network %s", network)
		}

		var d net.Dialer
		conn, err := d.DialContext(ctx, "tcp", proxyAddr)
		if err != nil {
			return nil, fmt.Errorf("socks5h: connect to proxy %s: %w", proxyAddr, err)
		}

		// Handshake
		if username != "" {
			// Offer no-auth and username/password methods
			if _, err := conn.Write([]byte{0x05, 0x02, 0x00, 0x02}); err != nil {
				conn.Close()
				return nil, err
			}
		} else {
			if _, err := conn.Write([]byte{0x05, 0x01, 0x00}); err != nil {
				conn.Close()
				return nil, err
			}
		}

		resp := make([]byte, 2)
		if _, err := io.ReadFull(conn, resp); err != nil {
			conn.Close()
			return nil, fmt.Errorf("socks5h: handshake read: %w", err)
		}
		if resp[0] != 0x05 {
			conn.Close()
			return nil, fmt.Errorf("socks5h: unexpected version %d", resp[0])
		}

		switch resp[1] {
		case 0x00:
			// No authentication required
		case 0x02:
			// Username/password authentication
			authMsg := make([]byte, 0, 3+len(username)+len(password))
			authMsg = append(authMsg, 0x01, byte(len(username)))
			authMsg = append(authMsg, []byte(username)...)
			authMsg = append(authMsg, byte(len(password)))
			authMsg = append(authMsg, []byte(password)...)
			if _, err := conn.Write(authMsg); err != nil {
				conn.Close()
				return nil, err
			}
			authResp := make([]byte, 2)
			if _, err := io.ReadFull(conn, authResp); err != nil {
				conn.Close()
				return nil, fmt.Errorf("socks5h: auth read: %w", err)
			}
			if authResp[1] != 0x00 {
				conn.Close()
				return nil, errors.New("socks5h: authentication failed")
			}
		default:
			conn.Close()
			return nil, fmt.Errorf("socks5h: unsupported auth method %d", resp[1])
		}

		// CONNECT with domain name (address type 0x03)
		host, portStr, err := net.SplitHostPort(addr)
		if err != nil {
			conn.Close()
			return nil, fmt.Errorf("socks5h: split host: %w", err)
		}
		port, err := strconv.Atoi(portStr)
		if err != nil {
			conn.Close()
			return nil, fmt.Errorf("socks5h: invalid port: %w", err)
		}
		if len(host) > 255 {
			conn.Close()
			return nil, fmt.Errorf("socks5h: hostname too long: %d", len(host))
		}

		req := make([]byte, 0, 7+len(host))
		req = append(req, 0x05, 0x01, 0x00, 0x03, byte(len(host)))
		req = append(req, []byte(host)...)
		req = append(req, byte(port>>8), byte(port&0xff))
		if _, err := conn.Write(req); err != nil {
			conn.Close()
			return nil, err
		}

		// Read CONNECT response: first 4 bytes (ver, rep, rsv, atyp) + variable
		// For domain name response: atyp=0x03, then 1 byte length + domain + 2 bytes port
		// For IPv4 response: atyp=0x01, then 4 bytes IP + 2 bytes port
		respHeader := make([]byte, 4)
		if _, err := io.ReadFull(conn, respHeader); err != nil {
			conn.Close()
			return nil, fmt.Errorf("socks5h: connect response header: %w", err)
		}
		if respHeader[1] != 0x00 {
			conn.Close()
			return nil, fmt.Errorf("socks5h: connect failed, reply code %d", respHeader[1])
		}

		// Skip the remaining bind address bytes
		var skipLen int
		switch respHeader[3] {
		case 0x01: // IPv4
			skipLen = 4 + 2
		case 0x03: // Domain name
			lenByte := make([]byte, 1)
			if _, err := io.ReadFull(conn, lenByte); err != nil {
				conn.Close()
				return nil, err
			}
			skipLen = int(lenByte[0]) + 2
		case 0x04: // IPv6
			skipLen = 16 + 2
		default:
			conn.Close()
			return nil, fmt.Errorf("socks5h: unknown address type %d", respHeader[3])
		}
		if skipLen > 0 {
			if _, err := io.CopyN(io.Discard, conn, int64(skipLen)); err != nil {
				conn.Close()
				return nil, err
			}
		}

		return conn, nil
	}, nil
}
