package localhost

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"testing"
	"time"
)

// 启动 HTTPS 服务（测试监听 8443）
func StartHTTPS() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintln(writer, "Hello, HTTPS world!")
	})

	server := &http.Server{
		Addr:    ":8443",
		Handler: mux,
	}
	homePath, _ := os.UserHomeDir()

	pem := path.Join(homePath, ".tlsctl", "certificates", "localhost", "localhost.pem")
	key := path.Join(homePath, ".tlsctl", "certificates", "localhost", "localhost.key")

	return server.ListenAndServeTLS(pem, key)
}

// 测试用例
func TestHTTPS(t *testing.T) {
	// 启动 HTTPS 服务
	go func() {
		err := StartHTTPS()
		if err != nil {
			t.Log("Server stopped:", err)
		}
	}()

	// 等待服务器启动
	time.Sleep(1 * time.Second)

	// 载入 CA 根证书
	homePath, _ := os.UserHomeDir()
	caCertPath := path.Join(homePath, ".tlsctl", "certificates", "localhost", "rootCA.pem") // 假设你有根证书
	caCert, err := os.ReadFile(caCertPath)
	if err != nil {
		t.Fatalf("Failed to read CA certificate: %v", err)
	}

	// 加入系统信任池
	rootCAs := x509.NewCertPool()
	if ok := rootCAs.AppendCertsFromPEM(caCert); !ok {
		t.Fatal("Failed to append CA certificate")
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: rootCAs,
				// InsecureSkipVerify: false 是默认值，可以省略
			},
		},
	}

	resp, err := client.Get("https://localhost:8443/")
	if err != nil {
		t.Fatalf("Failed to make HTTPS request: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	got := string(body)

	expected := "Hello, HTTPS world!\n"
	if got != expected {
		t.Errorf("Unexpected response body: got %q, want %q", got, expected)
	}
}
