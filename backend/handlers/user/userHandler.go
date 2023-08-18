package userHandler

import (
	"log"

	"github.com/RugeFX/ruge-chat-app/database"
	"github.com/RugeFX/ruge-chat-app/models"
	"github.com/gin-gonic/gin"
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
func GetAllUsers(c *gin.Context) {
	db := database.DB
	var users []models.User

	db.Find(&users)

	if len(users) == 0 {
		c.JSON(404, gin.H{
			"status": "failed",
			"error":  "Users not found",
		})
		return
	}

	c.JSON(200, gin.H{
		"status": "success",
		"data":   users,
	})
}

func GetUserByUsername(c *gin.Context) {
	db := database.DB
	var findUser models.User

	db.Where("LOWER(username) = LOWER(?)", c.Param("username")).Find(&findUser)

	if findUser.Username == "" {
		c.JSON(404, gin.H{
			"status": "failed",
			"error":  "User not found",
		})
		return
	}

	user := models.ReqUser{
		Username:       findUser.Username,
		Email:          findUser.Email,
		ProfilePicture: findUser.ProfilePicture,
	}

	c.JSON(200, gin.H{
		"status": "success",
		"data":   user,
	})
}

func CreateUser(c *gin.Context) {
	db := database.DB

	var user ReqCreateUser

	// TODO : Validate user

	if err := c.BindJSON(&user); err != nil {
		log.Println(err)
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
			c.JSON(400, gin.H{
				"status": "failed",
				"error":  "User already exists",
			})
			return
		}
		c.JSON(500, gin.H{
			"status": "failed",
			"error":  err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"status": "success",
		"data":   user,
	})
}

func DeleteUserByID(c *gin.Context) {
	db := database.DB
	var user models.User

	userId, err := uuid.Parse(c.Param("id"))

	if err != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"error":  "Invalid user ID format",
		})
		return
	}
	// fmt.Println("User ID: ", id)
	result := db.Find(&user, userId)

	if user.Username == "" {
		c.JSON(404, gin.H{
			"status": "failed",
			"error":  "Record not found",
		})
		return
	}

	result = result.Delete(&user)

	if err := result.Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(404, gin.H{
				"status": "failed",
				"error":  "Record not found",
			})
			return
		} else {
			c.JSON(500, gin.H{
				"status": "failed",
				"error":  err.Error(),
			})
			return
		}
	}

	c.JSON(200, gin.H{
		"status": "success",
		"data":   user,
	})
}

func UpdateUser(c *gin.Context) {
	db := database.DB
	var user models.User

	userId, err := uuid.Parse(c.Param("id"))

	if err != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"error":  "Invalid user ID format",
		})
		return
	}

	result := db.Find(&user, userId)

	if user.Username == "" {
		c.JSON(404, gin.H{
			"status": "failed",
			"error":  "Record not found",
		})
		return
	}

	if err := c.BindJSON(&user); err != nil {
		log.Println(err)
	}

	result = result.Save(&user)
	if err := result.Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(404, gin.H{
				"status": "failed",
				"error":  "Record not found",
			})
			return
		} else {
			c.JSON(500, gin.H{
				"status": "failed",
				"error":  err.Error(),
			})
			return
		}
	}

	c.JSON(200, gin.H{
		"status": "success",
		"data":   user,
	})
}

func hashPassword(p []byte) string {
	hash, _ := bcrypt.GenerateFromPassword(p, 8)
	return string(hash)
}
