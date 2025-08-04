package validation

import (
	"fmt"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New()

func ValidateBodyData(model any, c *fiber.Ctx) (err bool, errMsg any) {

	if err := c.BodyParser(model); err != nil {
		log.Fatal(err)
		return true,
			&fiber.Map{
				"success": false,
				"message": "Bad Request",
			}
	}

	if err := validate.Struct(model); err != nil {

		valErrors := err.(validator.ValidationErrors)[0]

		errMessage := fmt.Sprintf("Field: '%s' failed on: '%s' with your value: '%s'",
			valErrors.Field(), valErrors.Tag(), valErrors.Value())

		return true,
			&fiber.Map{
				"success": "false",
				"message": errMessage,
			}
	}
	return false, ""
}
