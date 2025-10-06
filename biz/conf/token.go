package conf

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
)

type Token struct {
	PublicKey string
	SecretKey string
	Expire    int64
}

func GetSecretKey(sk string) (*rsa.PrivateKey, error) {
	b, _ := pem.Decode([]byte(sk))
	privateKey, err := x509.ParsePKCS8PrivateKey(b.Bytes)
	if err != nil {
		return nil, err
	}
	return privateKey.(*rsa.PrivateKey), nil
}

func GetPublicKey(pk string) (*rsa.PublicKey, error) {
	b, _ := pem.Decode([]byte(pk))
	pubKey, err := x509.ParsePKIXPublicKey(b.Bytes)
	if err != nil {
		return nil, err
	}
	return pubKey.(*rsa.PublicKey), nil
}
