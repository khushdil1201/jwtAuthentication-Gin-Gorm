package initializers

import (
	"authentication/models"
)

func SyncDatabase() {
	DB.AutoMigrate(&models.User{})
}