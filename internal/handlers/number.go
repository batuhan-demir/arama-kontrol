package handlers

import (
	"arama-kontrol/internal/dal"
	"arama-kontrol/pkg/database"
	"github.com/gofiber/fiber/v2"

	"arama-kontrol/pkg/validation"
)

func GetNumbers(c *fiber.Ctx) error {
	var numbers []dal.Number

	res := database.DB.Find(&numbers)

	if res.Error != nil {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": "An error occurred while fetching numbers",
			"error":   res.Error.Error(),
		})
	}

	return c.JSON(&fiber.Map{
		"success": true,
		"data":    numbers,
	})
}

func CreateNumber(c *fiber.Ctx) error {
	number := new(dal.Number)

	err, errMsg := validation.ValidateBodyData(number, c)
	if err {
		return c.Status(400).JSON(errMsg)
	}

	if err := database.DB.Create(number).Error; err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": "An error occurred while creating the number",
			"error":   err.Error(),
		})
	}

	return c.Status(201).JSON(&fiber.Map{
		"success": true,
		"message": "Number created successfully",
		"data":    number,
	})
}

func DeleteNumber(c *fiber.Ctx) error {
	number := c.Params("number")

	if err := database.DB.Delete(&dal.Number{}, "Number = ?", number).Error; err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": "An error occurred while deleting the number",
			"error":   err.Error(),
		})
	}

	return c.JSON(&fiber.Map{
		"success": true,
		"message": "Number deleted successfully",
	})
}