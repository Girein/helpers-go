# Go Functions Helper

## Installation
`go get github.com/Girein/helpers-go`

## Functions
```
// ToDateTimeString converts DateTime into string with Y-m-d H:i:s format
func ToDateTimeString(dateTime time.Time) string {}

// GormOpen returns database connection
func GormOpen(driver string) *gorm.DB {}

// LogIfError logs the error with message
func LogIfError(err error, message string) {}

// RandomString generates random string with custom length
func RandomString(length int) string {}

// RandomInteger returns random integer between parameters
func RandomInteger(min int, max int) int {}

// JSONEncode converts data into JSON string
func JSONEncode(data interface{}) string {}

// InArray checks if a value exists in an array
func InArray(needle string, haystack []interface{}) bool {}

// AESEncrypt encrypts text using cipher AES/ECB/PKCS5PADDING
func AESEncrypt(text string, key []byte) string {}

// RSAVerifySignature verifies RSA PKCS #1 v1.5 signature with SHA256 hashing
func RSAVerifySignature(publicKey string, signature string, message string) bool {}
```