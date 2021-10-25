package main

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/subtle"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/fernet/fernet-go"
	"github.com/golang/gddo/httputil/header"
)

type envelope map[string]interface{}

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) invalidRequestResponse(w http.ResponseWriter, r *http.Request) {
	message := "invalid request"
	app.errorResponse(w, r, http.StatusBadRequest, message)
}

func (app *application) invalidSignatureResponse(w http.ResponseWriter, r *http.Request) {
	message := "invalid signature"
	app.errorResponse(w, r, http.StatusUnauthorized, message)
}

func verifySignature(secret string, payload string, inboundSignature []byte) bool {
	signature := generateSignature(secret, payload)
	return subtle.ConstantTimeCompare([]byte(signature), inboundSignature) == 1
}

func generateSignature(secret string, payload string) []byte {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(payload))
	expectedMAC := mac.Sum(nil)
	return expectedMAC
	// return hex.EncodeToString(expectedMAC)
}

func decryptMessage(privateKey, message string) (string, error) {
	s := strings.Split(message, ".")
	encryptedFernetKey, fernetToken := s[0], s[1]
	decodedFernetKey, err := base64.StdEncoding.DecodeString(encryptedFernetKey)
	if err != nil {
		return "", fmt.Errorf("error decoding key: %s", err)
	}
	decodedFernetToken, err := base64.StdEncoding.DecodeString(fernetToken)
	if err != nil {
		return "", fmt.Errorf("error decoding token: %s", err)
	}
	block, _ := pem.Decode([]byte(privateKey))
	if block.Type != "RSA PRIVATE KEY" {
		return "", errors.New("error decoding private key from pem")
	}
	parsedKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("error parsing private key: %s", err)
	}
	fernetKey, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, parsedKey, decodedFernetKey, nil)
	if err != nil {
		return "", fmt.Errorf("error decrypting fernet key: %s", err)
	}
	key, err := fernet.DecodeKey(string(fernetKey))
	if err != nil {
		return "", fmt.Errorf("error decoding fernet key: %s", err)
	}
	payload := fernet.VerifyAndDecrypt(decodedFernetToken, 0, []*fernet.Key{key})
	return string(payload), nil
}

func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message interface{}) {
	env := envelope{"error": message}
	err := app.writeJSON(w, status, env, nil)
	if err != nil {
		app.errorLog.Println(err)
		w.WriteHeader(500)
	}
}

func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}
	js = append(js, '\n')
	for key, value := range headers {
		w.Header()[key] = value
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)
	return nil
}

// renderJSON renders 'v' as JSON and writes it as a response into w.
func renderJSON(w http.ResponseWriter, v interface{}) {
	js, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

// For decoding json bodies better
// https://www.alexedwards.net/blog/how-to-properly-parse-a-json-request-body

type malformedRequest struct {
	status int
	msg    string
}

func (mr *malformedRequest) Error() string {
	return mr.msg
}

func decodeJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	if r.Header.Get("Content-Type") != "" {
		value, _ := header.ParseValueAndParams(r.Header, "Content-Type")
		if value != "application/json" {
			msg := "Content-Type header is not application/json"
			return &malformedRequest{status: http.StatusUnsupportedMediaType, msg: msg}
		}
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields() // TODO: May need to comment this given we don't know what they'll be sending us?

	err := dec.Decode(&dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case errors.Is(err, io.ErrUnexpectedEOF):
			msg := "Request body contains badly-formed JSON"
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case errors.Is(err, io.EOF):
			msg := "Request body must not be empty"
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case err.Error() == "http: request body too large":
			msg := "Request body must not be larger than 1MB"
			return &malformedRequest{status: http.StatusRequestEntityTooLarge, msg: msg}

		default:
			return err
		}
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		msg := "Request body must only contain a single JSON object"
		return &malformedRequest{status: http.StatusBadRequest, msg: msg}
	}

	return nil
}
