package account

import (
	"crypto"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/chihqiang/tlsctl/pkg/certificates"
	"github.com/chihqiang/tlsctl/pkg/fp"
	"github.com/go-acme/lego/v4/certcrypto"
)

const (
	// privateKeyFileName 是账号私钥文件名（无扩展名）。
	// 保持文件名 "keys" 以兼容历史版本已生成的账号数据。
	privateKeyFileName = "keys"
	accountFileName    = "account.json"
)

type Cache struct {
	email           string
	server          string
	accountFilePath string
	privateKeyPath  string
}

func NewCache(path, email string, server string) (*Cache, error) {
	var err error
	if !slices.Contains(AllowServers, server) {
		return nil, fmt.Errorf("unsupported ACME server: %s", server)
	}
	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}
	serverPath := strings.NewReplacer(":", "_", "/", string(os.PathSeparator)).Replace(serverURL.Host)
	accountsPath := filepath.Join(path, serverPath)
	rootUserPath := filepath.Join(accountsPath, email)
	if err := fp.CreateNonExistingFolder(rootUserPath); err != nil {
		return nil, err
	}
	return &Cache{
		email:           email,
		server:          server,
		privateKeyPath:  filepath.Join(rootUserPath, privateKeyFileName),
		accountFilePath: filepath.Join(rootUserPath, accountFileName),
	}, nil
}

func (c *Cache) GetEmail() string {
	return c.email
}

func (c *Cache) GetServer() string {
	return c.server
}

func (c *Cache) Save(account *Account) error {
	jsonBytes, err := json.MarshalIndent(account, "", "\t")
	if err != nil {
		return err
	}
	if err := os.WriteFile(c.accountFilePath, jsonBytes, 0o600); err != nil {
		return err
	}
	// 写入 PEM 格式私钥，与 LoadAccount/certificates.ParseX509 兼容，
	// 并支持 ECDSA/RSA 两种密钥类型，避免对具体类型做硬断言。
	privateKeyBlock := certcrypto.PEMBlock(account.Key)
	return os.WriteFile(c.privateKeyPath, pem.EncodeToMemory(privateKeyBlock), 0o600)
}
func (c *Cache) Remove() {
	_ = os.Remove(c.accountFilePath)
	_ = os.Remove(c.privateKeyPath)
}
func (c *Cache) LoadAccount() (*Account, error) {
	fileBytes, err := os.ReadFile(c.accountFilePath)
	if err != nil {
		return nil, err
	}
	var account Account
	err = json.Unmarshal(fileBytes, &account)
	if err != nil {
		return nil, err
	}
	privateKey, err := LoadPrivateKey(c.privateKeyPath)
	if err != nil {
		return nil, err
	}
	account.Key = privateKey
	return &account, nil
}

func LoadPrivateKey(file string) (crypto.PrivateKey, error) {
	keyBytes, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return certificates.ParseX509(keyBytes)
}

func GeneratePrivateKey(file string, keyType certcrypto.KeyType) (crypto.PrivateKey, error) {
	privateKey, err := certcrypto.GeneratePrivateKey(keyType)
	if err != nil {
		return nil, err
	}
	certOut, err := os.Create(file)
	if err != nil {
		return nil, err
	}
	defer certOut.Close()
	pemKey := certcrypto.PEMBlock(privateKey)
	err = pem.Encode(certOut, pemKey)
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}
