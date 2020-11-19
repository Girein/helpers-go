package helpers

import (
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
	if err != nil {
		slack.PostMessage("Failed to connect database.")
	}

	return db
}
