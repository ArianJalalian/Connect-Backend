package authService

import (
	"bitbucket.org/dyfrag-internal/mass-media-core/pkg/database"
	"bitbucket.org/dyfrag-internal/mass-media-core/pkg/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func User(c *fiber.Ctx) (*models.User, error) {
	user := new(models.User)
	result := database.DB.Find(user, "id = ?", Id(c))
	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func Id(c *fiber.Ctx) string {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	return claims["user_id"].(string)
}
