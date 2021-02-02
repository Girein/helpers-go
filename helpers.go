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
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// ToDateTimeString converts DateTime into string with Y-m-d H:i:s format
func ToDateTimeString(dateTime time.Time) string {
	return dateTime.Format("2006-01-02 15:04:05")
}

// GormOpen returns database connection
func GormOpen(driver string) *gorm.DB {
	db, err := gorm.Open(driver, os.Getenv("DB_CONNECTION_URL"))
	LogIfError(err, "Failed to connect database")

	return db
}

// LogIfError logs the error with message
func LogIfError(err error, message string) {
	if err != nil {
		log.Fatalf("%s: %s", message, err)
		slack.PostMessage("[" + os.Getenv("APP_NAME") + "]\r\n" + message)
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
func JSONEncode(data interface{}) string {
	jsonResult, err := json.Marshal(data)
	LogIfError(err, "JSON encode failed")

	return string(jsonResult)
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
func AESEncrypt(text string, key []byte) string {
	block, err := aes.NewCipher(key)
	LogIfError(err, "Failed to create cipher block")

	ecb := NewECBEncrypter(block)
	content := []byte(text)
	content = PKCS5Padding(content, block.BlockSize())
	crypted := make([]byte, len(content))
	ecb.CryptBlocks(crypted, content)
	cryptedString := base64.StdEncoding.EncodeToString(crypted)

	return cryptedString
}
