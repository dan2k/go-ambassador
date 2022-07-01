package database

import (
	"ambassador/src/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)
var DB *gorm.DB
func Connect(){
	var err error
	dsn := "root:root@tcp(db:3306)/ambassador?charset=utf8mb4&parseTime=True&loc=Local"
  	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err !=nil {
		panic("Could not connect database")
	}
}

func AutoMigrate(){
	DB.AutoMigrate(models.User{},models.Product{},models.Link{})
}