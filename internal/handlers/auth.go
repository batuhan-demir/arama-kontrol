package handlers

import (
	"arama-kontrol/internal/dal"
	"arama-kontrol/pkg/database"
	"arama-kontrol/pkg/hash"
	"arama-kontrol/pkg/validation"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// Helper function
func getSameSite() string {
	if os.Getenv("ENV") == "production" {
		return "None"
	}
	return "Lax"
}

func Register(c *fiber.Ctx) error {
	u := new(dal.UserCreate)

	err, errMsg := validation.ValidateBodyData(u, c)

	if err {
		return c.Status(400).JSON(errMsg)
	}

	// Check if user already exists
	var existingUser dal.User
	database.DB.First(&existingUser, "Email = ?", u.Email)
	if existingUser != (dal.User{}) {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": "User already exists",
		})
	}

	u.Password, _ = hash.HashPassword(u.Password)

	newUser := dal.User{
		Name:     u.Name,
		Surname:  u.Surname,
		Email:    u.Email,
		Phone:    u.Phone,
		Password: u.Password,
	}

	res := database.DB.Create(&newUser)

	if res.Error != nil {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": "An error occured in server. Please try again later",
			"error":   res.Error.Error(),
		})
	}

	// save the number to numbers list if it doesn't exist
	var existingNumber dal.Number
	database.DB.First(&existingNumber, "Number = ?", newUser.Phone)
	
	if existingNumber == (dal.Number{}) {
		newNumber := dal.Number{
			Number: newUser.Phone,
			Name:   newUser.Name + " " + newUser.Surname,
		}
		database.DB.Create(&newNumber)
	}

	return c.Status(201).JSON(&fiber.Map{
		"success": true,
		"message": "User Created Successfully",
		"user":    newUser,
	})
}

func Login(c *fiber.Ctx) error {

	loginData := new(dal.UserLogin)

	err, errMsg := validation.ValidateBodyData(loginData, c)

	if err {
		return c.Status(400).JSON(errMsg)
	}

	var user dal.User
	result := database.DB.First(&user, "Email = ?", loginData.Email)

	if result.Error != nil || !hash.CheckPasswordHash(loginData.Password, user.Password) {
		return c.Status(403).JSON(&fiber.Map{
			"success": false,
			"message": "Invalid credentials",
		})
	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = user.Id
	claims["email"] = user.Email
	claims["name"] = user.Name
	claims["surname"] = user.Surname
	claims["phone"] = user.Phone
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	t, _err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))

	if _err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    t,
		HTTPOnly: true,
		Secure:   os.Getenv("ENV") == "production",
		SameSite: getSameSite(),
	})

	return c.JSON(&fiber.Map{
		"success": true,
		"message": "Login Successful",
		"data":    user,
	})

}

func CheckAuth(c *fiber.Ctx) error {

	id := c.Locals("user").(jwt.MapClaims)["id"]

	var user dal.User

	database.DB.First(&user, "Id = ?", id)

	if user == (dal.User{}) {
		return c.Status(404).JSON(&fiber.Map{
			"success": false,
			"message": "User not found",
		})
	}

	return c.JSON(&fiber.Map{
		"success": true,
		"data":    user,
	})

}

func Logout(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    "",
		MaxAge:   -1, // Cookie'yi silmek için
		HTTPOnly: true,
		Secure:   os.Getenv("ENV") == "production",
		SameSite: func() string {
			if os.Getenv("ENV") == "production" {
				return "None"
			}
			return "Lax"
		}(),
	})

	return c.JSON(&fiber.Map{
		"success": true,
		"message": "Logout successful",
	})
}
