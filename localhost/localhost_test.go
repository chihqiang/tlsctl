package localhost

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// startHTTPS 在随机端口上启动一个 HTTPS 服务，返回监听地址。
func startHTTPS(t *testing.T, certFile, keyFile string) string {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, HTTPS world!")
	})
	server := &http.Server{Handler: mux}
	go func() {
		_ = server.ServeTLS(ln, certFile, keyFile)
	}()
	t.Cleanup(func() { _ = server.Close() })
	return ln.Addr().String()
}

func TestHTTPS(t *testing.T) {
	// 在临时目录中生成 CA 与证书，避免依赖用户的 ~/.tlsctl
	dir := t.TempDir()
	hostSSL, err := NewLocalHostSSL(dir)
	if err != nil {
		t.Fatalf("NewLocalHostSSL: %v", err)
	}
	if err := hostSSL.LoadCA(); err != nil {
		t.Fatalf("LoadCA: %v", err)
	}
	res, err := hostSSL.BuildResource([]string{"localhost", "127.0.0.1"})
	if err != nil {
		t.Fatalf("BuildResource: %v", err)
	}
	certFile := filepath.Join(dir, "localhost.pem")
	keyFile := filepath.Join(dir, "localhost.key")
	if err := os.WriteFile(certFile, res.Certificate, 0o600); err != nil {
		t.Fatalf("write certificate: %v", err)
	}
	if err := os.WriteFile(keyFile, res.PrivateKey, 0o600); err != nil {
		t.Fatalf("write key: %v", err)
	}

	addr := startHTTPS(t, certFile, keyFile)

	// 用生成的 CA 验证服务端证书
	rootCAs := x509.NewCertPool()
	if !rootCAs.AppendCertsFromPEM(res.IssuerCertificate) {
		t.Fatal("failed to append CA certificate")
	}
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:    rootCAs,
				ServerName: "localhost",
			},
		},
		Timeout: 5 * time.Second,
	}

	// 轮询等待服务就绪，避免固定 sleep
	url := fmt.Sprintf("https://%s/", addr)
	deadline := time.Now().Add(5 * time.Second)
	var body []byte
	for {
		resp, err := client.Get(url)
		if err == nil {
			body, _ = io.ReadAll(resp.Body)
			resp.Body.Close()
			break
		}
		if time.Now().After(deadline) {
			t.Fatalf("HTTPS request failed: %v", err)
		}
		time.Sleep(50 * time.Millisecond)
	}

	expected := "Hello, HTTPS world!\n"
	if got := string(body); got != expected {
		t.Errorf("Unexpected response body: got %q, want %q", got, expected)
	}
}
