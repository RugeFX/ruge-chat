package userHandlers

import (
	"github.com/RugeFX/ruge-chat-app/database"
	"github.com/RugeFX/ruge-chat-app/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type ReqCreateUser struct {
	Username string `json:"username"`
	Password string `json:"-"`
	Email    string `json:"email"`
}

// Gets all users
func GetAllUsers(c *fiber.Ctx) error {
	db := database.DB
	var users []models.User

	db.Find(&users)

	if len(users) == 0 {
		return c.Status(404).JSON(fiber.Map{
			"status": "failed",
			"error":  "Users not found",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status": "success",
		"data":   users,
	})
}

func GetUserByUsername(c *fiber.Ctx) error {
	db := database.DB
	var findUser models.User

	db.Where("LOWER(username) = LOWER(?)", c.Params("username")).Find(&findUser)

	if findUser.Username == "" {
		return c.Status(404).JSON(fiber.Map{
			"status": "failed",
			"error":  "User not found",
		})
	}

	user := models.ReqUser{
		Username:       findUser.Username,
		Email:          findUser.Email,
		ProfilePicture: findUser.ProfilePicture,
	}

	return c.Status(200).JSON(fiber.Map{
		"status": "success",
		"data":   user,
	})
}

func CreateUser(c *fiber.Ctx) error {
	db := database.DB

	var user ReqCreateUser

	// TODO : Validate user

	if err := c.BodyParser(&user); err != nil {
		return err
	}

	newPass := hashPassword([]byte(user.Password))

	newUser := models.User{
		Username: user.Username,
		Password: newPass,
		Email:    user.Email,
	}

	 err := db.Create(&newUser).Error

	if err != nil {
		if err == gorm.ErrDuplicatedKey {
			return c.Status(400).JSON(fiber.Map{
				"status": "failed",
				"error":  "User already exists",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"status": "failed",
			"error":  err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status": "success",
		"data":   user,
	})
}

func DeleteUserByID(c *fiber.Ctx) error {
	db := database.DB
	var user models.User

	userId, err := uuid.Parse(c.Params("id"))

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status": "failed",
			"error":  "Invalid user ID format",
		})
	}
	// fmt.Println("User ID: ", id)
	result := db.Find(&user, userId)

	if user.Username == "" {
		return c.Status(404).JSON(fiber.Map{
			"status": "failed",
			"error":  "Record not found",
		})
	}

	result = result.Delete(&user)

	if err := result.Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(404).JSON(fiber.Map{
				"status": "failed",
				"error":  "Record not found",
			})
		} else {
			return c.Status(500).JSON(fiber.Map{
				"status": "failed",
				"error":  err.Error(),
			})
		}
	}

	return c.Status(200).JSON(fiber.Map{
		"status": "success",
		"data":   user,
	})
}

func UpdateUser(c *fiber.Ctx) error {
	db := database.DB
	var user models.User

	userId, err := uuid.Parse(c.Params("id"))

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status": "failed",
			"error":  "Invalid user ID format",
		})
	}

	result := db.Find(&user, userId)

	if user.Username == "" {
		return c.Status(404).JSON(fiber.Map{
			"status": "failed",
			"error":  "Record not found",
		})
	}

	if err := c.BodyParser(&user); err != nil {
		return err
	}

	result = result.Save(&user)
	if err := result.Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(404).JSON(fiber.Map{
				"status": "failed",
				"error":  "Record not found",
			})
		} else {
			return c.Status(500).JSON(fiber.Map{
				"status": "failed",
				"error":  err.Error(),
			})
		}
	}

	return c.Status(200).JSON(fiber.Map{
		"status": "success",
		"data":   user,
	})
}

func hashPassword(p []byte) string {
	hash, _ := bcrypt.GenerateFromPassword(p, 8)
	return string(hash)
}
