package yorha_qq_auth_key

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"

	_ "embed"
)

//go:embed private.key
var QQAuthKeyBytes []byte
var QQAuthKey *rsa.PrivateKey

func init() {
	var err error
	keyBlock, _ := pem.Decode(QQAuthKeyBytes)
	QQAuthKey, err = x509.ParsePKCS1PrivateKey(keyBlock.Bytes)
	if err != nil {
		panic(err)
	}
}
