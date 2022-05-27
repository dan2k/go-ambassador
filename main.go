package main

import (
    "log"

    "github.com/gofiber/fiber/v2"
	"gorm.io/driver/mysql"
  	"gorm.io/gorm"
)

func main() {
	dsn := "root:root@tcp(db:3306)/ambassador?charset=utf8mb4&parseTime=True&loc=Local"
  	_, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err !=nil {
		panic("Could not connect database")
	}
    app := fiber.New()

    app.Get("/", func (c *fiber.Ctx) error {
        return c.SendString("Hello, World!")
    })

    log.Fatal(app.Listen(":3000"))
}