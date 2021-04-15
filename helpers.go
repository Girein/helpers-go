package helpers

import (
	"crypto/aes"
	"encoding/base64"
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/Girein/slack-incoming-webhook-go"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// ToDateTimeString converts DateTime into string with Y-m-d H:i:s format
func ToDateTimeString(dateTime time.Time) string {
	return dateTime.Format("2006-01-02 15:04:05")
}

// GormOpen returns database connection
func GormOpen(connectionUrl string) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(connectionUrl), &gorm.Config{})

	return db, err
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
