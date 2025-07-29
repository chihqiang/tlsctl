package account

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/x509"
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
	"golang.org/x/sync/errgroup"
)

const (
	baseKeysFolderName = "keys"
	accountFileName    = "account.json"
	userKey            = "key"
)

type Cache struct {
	email           string
	server          string
	accountFilePath string
	keysPath        string
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
		keysPath:        filepath.Join(rootUserPath, baseKeysFolderName),
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
	var (
		wg errgroup.Group
	)
	wg.Go(func() error {
		jsonBytes, err := json.MarshalIndent(account, "", "\t")
		if err != nil {
			return err
		}
		return os.WriteFile(c.accountFilePath, jsonBytes, 0o600)
	})
	wg.Go(func() error {
		privateKeyBytes, err := x509.MarshalECPrivateKey(account.Key.(*ecdsa.PrivateKey))
		if err != nil {
			return err
		}
		return os.WriteFile(c.keysPath, privateKeyBytes, 0o600)
	})
	return wg.Wait()
}
func (c *Cache) Remove() {
	_ = os.Remove(c.accountFilePath)
	_ = os.Remove(c.keysPath)
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
	privateKey, err := LoadPrivateKey(c.keysPath)
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
