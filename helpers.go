package helpers

import (
	"crypto/aes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/Girein/slack-incoming-webhook-go"
	"github.com/forgoer/openssl"
	"github.com/techoner/gophp/serialize"
)

// ToDateTimeString converts DateTime into string with Y-m-d H:i:s format
func ToDateTimeString(dateTime time.Time) string {
	return dateTime.Format("2006-01-02 15:04:05")
}

// LogIfError logs the error with message
func LogIfError(err error, message string) {
	log.Printf("%s: %s", message, err)

	if os.Getenv("APP_ENV") != "production" {
		err = slack.PostMessage("[" + os.Getenv("APP_NAME") + "]\r\n" + message)
		if err != nil {
			log.Printf("Failed to post message to Slack: %s: %s", err.Error(), err)
		}
	}
}

// RandomString generates random string with custom length
func RandomString(length int) string {
	bytes := make([]byte, length)

	for i := 0; i < length; i++ {
		bytes[i] = byte(RandomInteger(65, 90))
	}

	return string(bytes)
}

// RandomInteger returns random integer between parameters
func RandomInteger(min int, max int) int {
	return min + rand.Intn(max-min)
}

// JSONEncode converts data into JSON string
func JSONEncode(data interface{}) (string, error) {
	jsonResult, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	return string(jsonResult), nil
}

// InArray checks if a value exists in an array
func InArray(needle string, haystack []interface{}) bool {
	for _, value := range haystack {
		if needle == value {
			return true
		}
	}

	return false
}

// AESEncrypt encrypts text using cipher AES/ECB/PKCS5PADDING
func AESEncrypt(text string, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	ecb := NewECBEncrypter(block)
	content := []byte(text)
	content = PKCS5Padding(content, block.BlockSize())
	crypted := make([]byte, len(content))
	ecb.CryptBlocks(crypted, content)
	cryptedString := base64.StdEncoding.EncodeToString(crypted)

	return cryptedString, nil
}

// RSAVerifySignature verifies RSA PKCS #1 v1.5 signature with SHA256 hashing
func RSAVerifySignature(publicKey string, signature string, message string) (bool, error) {
	parser, err := parsePublicKey([]byte(publicKey))
	if err != nil {
		return false, err
	}

	decodedSignature, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return false, err
	}

	err = parser.Unsign([]byte(message), []byte(decodedSignature))
	if err != nil {
		return false, err
	}

	return true, nil
}

// RandomBytes generates random byte with custom length
func RandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// OpenSSLEncrypt encrypts given data with given key, returns base64 encoded string
func OpenSSLEncrypt(data []byte, passphrase []byte, iv []byte) (string, error) {
	res, err := openssl.AesCBCEncrypt(data, passphrase, iv, openssl.PKCS7_PADDING)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(res), nil
}

// ComputeHMACSHA256 hashes given message with given secret, returns hexadecimal encoded string
func ComputeHMACSHA256(message string, secret string) (string, error) {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	_, err := h.Write([]byte(message))
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

// LaravelEncrypt encrypts the given value using Laravel's encrypter (https://laravel.com/docs/6.x/encryption)
func LaravelEncrypt(value string) (string, error) {
	iv, err := RandomBytes(16)
	if err != nil {
		return "", err
	}

	message, err := serialize.Marshal(value)
	if err != nil {
		return "", err
	}

	key := os.Getenv("APP_KEY")

	resVal, err := OpenSSLEncrypt(message, []byte(key), iv)
	if err != nil {
		return "", err
	}

	resIv := base64.StdEncoding.EncodeToString(iv)

	data := resIv + resVal
	mac, err := ComputeHMACSHA256(data, key)
	if err != nil {
		return "", err
	}

	ticket := make(map[string]interface{})
	ticket["iv"] = resIv
	ticket["mac"] = mac
	ticket["value"] = resVal

	resTicket, err := json.Marshal(ticket)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(resTicket), nil
}
