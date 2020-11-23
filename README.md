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
```