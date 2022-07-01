package controllers

import (
	"ambassador/src/database"
	"ambassador/src/middlewares"
	"ambassador/src/models"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

func Register(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	if data["password"] != data["password_confirm"] {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "password do not match",
		})
	}

	user := models.User{
		FirstName:    data["first_name"],
		LastName:     data["last_name"],
		Email:        data["email"],
		IsAmbassador: false,
	}
	user.SetPassword(data["password"])
	if err := database.DB.Create(&user); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"error": err.Error,
		})
	}
	return c.JSON(user)
}
func Login(c *fiber.Ctx) error {
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return err
	}

	var user models.User
	database.DB.Where("email = ?", data["email"]).First(&user)
	if user.Id == 0 {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Invalid Credentials 1",
		})
	}

	if err := user.ComparePassword(data["password"]); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Invalid Credentials 2",
		})
	}

	payload := jwt.StandardClaims{
		Subject:   strconv.Itoa(int(user.Id)),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, payload).SignedString([]byte("secret"))
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Invalid Credentials 3",
			"error":   err.Error(),
		})
	}
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)
	// return c.JSON(token)
	return c.JSON(fiber.Map{
		"message": "success",
	})

}
func User(c *fiber.Ctx) error {

	var user models.User
	id, _ := middlewares.GetUserId(c)
	database.DB.Where("id = ?", id).First(&user)
	return c.JSON(user)
}
func Logout(c *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(time.Hour),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)
	return c.JSON(fiber.Map{
		"message": "success",
	})
}
func UpdateInfo(c *fiber.Ctx) error {
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return err
	}
	id, _ := middlewares.GetUserId(c)
	user := models.User{
		Id:        id,
		FirstName: data["first_name"],
		LastName:  data["last_name"],
		Email:     data["email"],
	}

	database.DB.Model(&user).Updates(&user)
	return c.JSON(user)
}
func UpdatePassword(c *fiber.Ctx) error {
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return err
	}
	if data["password"] != data["password_confirm"] {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "password do not match",
		})
	}


	id, _ := middlewares.GetUserId(c)
	user := models.User{
		Id:        id,
	}
	user.SetPassword(data["password"])

	database.DB.Model(&user).Updates(&user)
	return c.JSON(user)
}
