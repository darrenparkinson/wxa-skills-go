package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"log"
	"os"
)

func main() {
	es := encryptString("Hello")
	log.Println(es)
	ds := decryptString(es)
	log.Println(ds)
}

func encryptString(s string) string {
	publicKey, _ := os.ReadFile("public.pem")
	block, _ := pem.Decode([]byte(publicKey))
	if block.Type != "PUBLIC KEY" {
		log.Fatal("error decoding public key from pem")
	}
	parsedKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		log.Fatal("error parsing key")
	}
	var ok bool
	var pubkey *rsa.PublicKey
	if pubkey, ok = parsedKey.(*rsa.PublicKey); !ok {
		log.Fatal("unable to parse public key")
	}
	rng := rand.Reader
	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rng, pubkey, []byte(s), nil)
	if err != nil {
		log.Fatal(err)
	}
	return base64.StdEncoding.EncodeToString(ciphertext)
}

func decryptString(s string) string {
	privateKey, _ := os.ReadFile("private.pem")
	message, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		log.Fatal(err)
	}
	block, _ := pem.Decode([]byte(privateKey))
	if block.Type != "RSA PRIVATE KEY" {
		log.Fatal("error decoding private key from pem")
	}
	parsedKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		log.Fatal(err)
	}
	rng := rand.Reader
	// Both of these work
	plaintext, err := parsedKey.Decrypt(rng, message, &rsa.OAEPOptions{Hash: crypto.SHA256})
	// plaintext, err := rsa.DecryptOAEP(sha256.New(), rng, parsedKey, message, nil)
	if err != nil {
		log.Fatal(err)
	}
	return string(plaintext)

}
