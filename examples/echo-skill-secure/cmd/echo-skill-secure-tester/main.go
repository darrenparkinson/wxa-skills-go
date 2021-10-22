package main

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/fernet/fernet-go"

	wxas "github.com/darrenparkinson/wxa-skills-go"
)

type payload struct {
	Signature string `json:"signature"`
	Message   string `json:"message"`
}

type config struct {
	publicKey string
	secret    string
}

func main() {
	cfg := loadEnv()
	message := generateJSONMessage()
	payloadJSON, err := preparePayload(message, cfg)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := http.Post("http://localhost:8080", "application/json", strings.NewReader(payloadJSON))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	log.Println("Status:", resp.Status)
	log.Println("Response:")
	var result wxas.WebexAssistantResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}
	jres, err := json.MarshalIndent(result, "  ", " ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(jres))

}

// preparePayload returns the json string required to send to the skill, simulating webex assistant
func preparePayload(message string, cfg config) (string, error) {
	var result payload
	token, err := generateToken(message, cfg.publicKey)
	if err != nil {
		return "", err
	}
	signature := signToken(token, cfg.secret)
	if err != nil {
		return "", err
	}
	result.Signature = signature
	result.Message = token
	resultJSON, err := json.MarshalIndent(result, "  ", " ")
	if err != nil {
		return "", err
	}
	return string(resultJSON), nil
}

// generateToken takes a plaintext message and a public key. It uses a fernet key to encrypt
// the message and encrypts the fernet key with the public key.  It returns the encrypted fernet
// key and the encrypted message each individually base64 encoded and separated by a ".".
// IKR 🙄
func generateToken(message string, publicKey string) (string, error) {
	var fernetKey fernet.Key
	err := fernetKey.Generate()
	if err != nil {
		return "", err
	}
	encryptedMessage, err := fernet.EncryptAndSign([]byte(message), &fernetKey)
	if err != nil {
		return "", err
	}
	encodedEncryptedFernetKey, err := encryptFernetKey(fernetKey.Encode(), publicKey)
	if err != nil {
		return "", err
	}
	encodedEncryptedMessage := base64.StdEncoding.EncodeToString(encryptedMessage)
	return fmt.Sprintf("%s.%s", encodedEncryptedFernetKey, encodedEncryptedMessage), nil
}

func encryptFernetKey(fernetKey string, publicKey string) (string, error) {
	block, _ := pem.Decode([]byte(publicKey))
	if block.Type != "PUBLIC KEY" {
		return "", errors.New("error decoding public key from pem")
	}
	parsedKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("error parsing public key: %s", err)
	}
	var ok bool
	var pubkey *rsa.PublicKey
	if pubkey, ok = parsedKey.(*rsa.PublicKey); !ok {
		return "", errors.New("unable to parse public key")
	}
	rng := rand.Reader
	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rng, pubkey, []byte(fernetKey), nil)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func signToken(token string, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(token))
	signature := mac.Sum(nil)
	return base64.StdEncoding.EncodeToString(signature)
}

func loadEnv() config {
	pub, err := os.ReadFile("public.pem")
	if err != nil {
		log.Fatal(err)
	}
	secret, err := os.ReadFile("secret.txt")
	if err != nil {
		log.Fatal(err)
	}
	cfg := config{
		publicKey: string(pub),
		secret:    string(secret),
	}
	if cfg.publicKey == "" || cfg.secret == "" {
		log.Fatal("missing environment variables")
	}
	return cfg
}

func generateJSONMessage() string {
	challenge, _ := generateChallenge()
	message := wxas.WebexAssistantMessage{
		// Text:      []string{"Hello World."},
		Text:      "Hello World.",
		Challenge: challenge,
	}
	j, _ := json.MarshalIndent(message, "  ", " ")
	return string(j)
}

func generateChallenge() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
