package serve

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"bitbucket.org/dyfrag-internal/mass-media-core/pkg/cli/serve/controller/dto"
	serve "bitbucket.org/dyfrag-internal/mass-media-core/pkg/cli/serve/service"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type AdminController struct{}

func (c *AdminController) RegisterRoutes(group fiber.Router) {
	group.Get("/users", c.GetUsers)
	group.Post("/sport", c.AddSport)
	group.Get("/reports", c.GetReports)
	group.Put("/report-block", c.BlockUser)
	group.Get("/sports", c.GetSports)
}

// GetUsers
// @Summary Get users
// @Description get all users
// @Tags admin
// @Accept json
// @Produce json
// @Param Authorization header string true "JWT token"
// @Success 200 {object} []dto.User "User information"
// @Failure 400 {object} string "Invalid request payload"
// @Failure 500 {object} string "Internal Server Error"
// @Router /admin/users [get]
func (c *AdminController) GetUsers(ctx *fiber.Ctx) error {
	tokenHeader := ctx.Get("Authorization")
	if tokenHeader == "" {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "Authorization header missing or invalid"})
	}

	token, err := jwt.Parse(tokenHeader, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})
	if err != nil {
		return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid JWT token"})
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid JWT token"})
	}
	userID := uint(claims["user_id"].(float64))
	fmt.Println(userID)
	users := serve.GetUsers()
	return ctx.JSON(users)
}

// AddSport
// @Summary Add sport
// @Description add sport
// @Tags admin
// @Accept json
// @Produce json
// @Param Authorization header string true "JWT token"
// @Param Sport body dto.Sport true "Sport information"
// @Success 200 {object} dto.SportResponse "Sport information"
// @Failure 400 {object} string "Invalid request payload"
// @Failure 500 {object} string "Internal Server Error"
// @Router /admin/sport [post]
func (c *AdminController) AddSport(ctx *fiber.Ctx) error {
	tokenHeader := ctx.Get("Authorization")
	if tokenHeader == "" {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "Authorization header missing or invalid"})
	}

	token, err := jwt.Parse(tokenHeader, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})
	if err != nil {
		return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid JWT token"})
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid JWT token"})
	}
	userID := uint(claims["user_id"].(float64))
	fmt.Println(userID)
	sport := dto.Sport{}
	if err := ctx.BodyParser(&sport); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request payload"})
	}
	s, err := serve.AddSport(sport)
	if err != nil {
		return err
	}
	res := dto.SportResponse{
		Title:       s.Title,
		Description: s.Description,
		VideoID:     s.VideoID,
		Path:        s.Video.Path,
		Name:        s.Video.Name,
		Type:        s.Video.Type,
	}
	return ctx.JSON(res)
}

// GetReports
// @Summary Get reports
// @Description get reports
// @Tags admin
// @Accept json
// @Produce json
// @Param Authorization header string true "JWT token"
// @Success 200 {object} []dto.ReportResponse "Report information"
// @Failure 400 {object} string "Invalid request payload"
// @Failure 500 {object} string "Internal Server Error"
// @Router /admin/reports [get]
func (c *AdminController) GetReports(ctx *fiber.Ctx) error {
	tokenHeader := ctx.Get("Authorization")
	if tokenHeader == "" {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "Authorization header missing or invalid"})
	}

	token, err := jwt.Parse(tokenHeader, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})
	if err != nil {
		return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid JWT token"})
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid JWT token"})
	}
	userID := uint(claims["user_id"].(float64))
	fmt.Println(userID)
	reports := serve.GetReports()
	var reportRes []dto.ReportResponse
	for _, r := range reports {
		report := dto.ReportResponse{
			Description: r.Description,
			UserID:      r.UserID,
			ReportID:    r.ID,
		}
		reportRes = append(reportRes, report)
	}
	return ctx.JSON(reportRes)
}

// BlockUser
// @Summary Block user
// @Description block user
// @Tags admin
// @Accept json
// @Produce json
// @Param Authorization header string true "JWT Token"
// @Param report-id query uint true "Report ID"
// @Success 200 {object} dto.Response "Report information"
// @Failure 400 {object} string "Invalid request payload"
// @Failure 500 {object} string "Internal Server Error"
// @Router /admin/report-block [put]
func (c *AdminController) BlockUser(ctx *fiber.Ctx) error {
	tokenHeader := ctx.Get("Authorization")
	if tokenHeader == "" {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "Authorization header missing or invalid"})
	}

	token, err := jwt.Parse(tokenHeader, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})
	if err != nil {
		return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid JWT token"})
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid JWT token"})
	}
	userID := uint(claims["user_id"].(float64))
	fmt.Println(userID)
	reportID := ctx.Query("report-id")
	report, err := strconv.ParseUint(reportID, 10, 32)
	if err != nil {
		return err
	}
	blocked, err := serve.BlockUser(uint(report))
	if err != nil {
		return err
	}
	fmt.Println(blocked)
	res := dto.Response{
		Message: "True",
		Success: true,
		ID:      uint(report),
	}
	return ctx.JSON(res)
}

// GetSports
// @Summary Get sports
// @Description get sports
// @Tags admin
// @Accept json
// @Produce json
// @Param Authorization header string true "JWT Token"
// @Success 200 {object} []dto.SportResponse "Sports information"
// @Failure 400 {object} string "Invalid request payload"
// @Failure 500 {object} string "Internal Server Error"
// @Router /admin/sports [get]
func (c *AdminController) GetSports(ctx *fiber.Ctx) error {
	tokenHeader := ctx.Get("Authorization")
	if tokenHeader == "" {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "Authorization header missing or invalid"})
	}

	token, err := jwt.Parse(tokenHeader, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})
	if err != nil {
		return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid JWT token"})
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid JWT token"})
	}
	userID := uint(claims["user_id"].(float64))
	fmt.Println(userID)
	sports, err := serve.GetSports()
	if err != nil {
		return err
	}
	var sportRes []dto.SportResponse
	for _, s := range sports {
		sport := dto.SportResponse{
			Title:       s.Title,
			Description: s.Description,
			VideoID:     s.VideoID,
			Path:        s.Video.Path,
			Name:        s.Video.Name,
			Type:        s.Video.Type,
		}
		sportRes = append(sportRes, sport)
	}
	return ctx.JSON(sportRes)
}

func (c *AdminController) HelloWorld(ctx *fiber.Ctx) error {
	tokenHeader := ctx.Get("Authorization")
	if tokenHeader == "" {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "Authorization header missing or invalid"})
	}

	token, err := jwt.Parse(tokenHeader, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})
	if err != nil {
		return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid JWT token"})
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid JWT token"})
	}
	userID := uint(claims["user_id"].(float64))
	return ctx.JSON(userID)
}
