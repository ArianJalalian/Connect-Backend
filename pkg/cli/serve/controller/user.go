package serve

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/golang-jwt/jwt/v4"

	"bitbucket.org/dyfrag-internal/mass-media-core/pkg/cli/serve/controller/dto"
	serve "bitbucket.org/dyfrag-internal/mass-media-core/pkg/cli/serve/service"
	"bitbucket.org/dyfrag-internal/mass-media-core/pkg/utils"
	"bitbucket.org/dyfrag-internal/mass-media-core/pkg/utils/authService"
	"github.com/gofiber/fiber/v2"
)

type UserController struct{}

func (c *UserController) RegisterRoutes(group fiber.Router) {
	group.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to the user group!")
	})
	group.Post("/login", c.LogIn)

	group.Post("/sign-up", c.SignUp)
	group.Get("/:id/profile", c.GetProfile)
	group.Put("/profile", c.EditProfile)
}

// LogIn
// @Summary Log in user
// @Description Logs in a user using email and password
// @Tags user
// @Accept json
// @Produce json
// @Param request body dto.LogIn true "Email and password"
// @Success 200 {object} dto.UserResponse "User information"
// @Failure 400 {object} string "Invalid request payload"
// @Failure 500 {object} string "Internal Server Error"
// @Router /user/login [post]
func (c *UserController) LogIn(ctx *fiber.Ctx) error {
	var loginRequest dto.LogIn
	if err := ctx.BodyParser(&loginRequest); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request payload"})
	}

	user, err := serve.GetUser(loginRequest.Email, loginRequest.Password)
	if err != nil {
		fmt.Println(err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Cannot get user", "error": err})
	}
	token, err := authService.JwtGenerator(user)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to generate JWT token", "error": err})
	}
	userR := dto.UserResponse{
		ID:          user.ID,
		Email:       user.Email,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		PhoneNumber: user.PhoneNumber,
		InfoID:      user.InfoID,
		InfoType:    user.InfoType,
		Block:       user.Block,
		Wallet:      user.Wallet,
		JWT:         token,
	}
	return ctx.JSON(userR)
}

// SignUp
// @Summary Sign up user
// @Description Signs up a new user with provided details
// @Tags user
// @Accept json
// @Produce json
// @Param user body dto.User true "User details"
// @Success 200 {object} dto.UserResponse "User information"
// @Failure 400 {object} string "Invalid request payload"
// @Failure 500 {object} string "Internal Server Error"
// @Router /user/sign-up [post]
func (c *UserController) SignUp(ctx *fiber.Ctx) error {
	var user dto.User
	if err := ctx.BodyParser(&user); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request payload"})
	}

	if err := utils.ValidateUser(user); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"message": err})
	}

	if err := user.HashPassword(); err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to hash password"})
	}
	userModel, err := serve.CreateUser(user)
	if err != nil {
		ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "cannot create user", "error": err})
	}
	token, err := authService.JwtGenerator(userModel)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to generate JWT token", "error": err})
	}
	fmt.Println(userModel.ID)
	userR := dto.UserResponse{
		ID:          userModel.ID,
		Email:       userModel.Email,
		FirstName:   userModel.FirstName,
		LastName:    userModel.LastName,
		PhoneNumber: userModel.PhoneNumber,
		InfoID:      userModel.InfoID,
		InfoType:    userModel.InfoType,
		Block:       userModel.Block,
		Wallet:      userModel.Wallet,
		JWT:         token,
	}
	return ctx.JSON(userR)
}

// GetProfile
// @Summary Get a User
// @Description get user profile by id
// @Tags user
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} dto.UserResponse "User information"
// @Failure 400 {object} string "Invalid request payload"
// @Failure 500 {object} string "Internal Server Error"
// @Router /user/:id/profile [get]
func (c *UserController) GetProfile(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)

	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "Invalid user ID"})
	}

	user, err := serve.GetUserById(id)

	if err != nil {
		return ctx.Status(http.StatusNotFound).JSON(fiber.Map{"message": "User not found"})
	}

	if user.InfoType == "trainer" {
		trainer, _ := serve.GetTrainerByUserID(uint(id))

		return ctx.JSON(fiber.Map{
			"user":    user,
			"profile": trainer,
		})

	} else {
		trainee, _ := serve.GetTraineeByUserID(uint(id))

		return ctx.JSON(fiber.Map{
			"user":    user,
			"profile": trainee,
		})
	}
}

// EditProfile
//
// @Summary Edit user profile
// @Description Edit user profile based on the user's role (trainer or trainee)
// @Tags user
// @Accept json
// @Produce json
// @Param Authorization header string true "JWT Token"
// @Param body body dto.UserEditTraineeOrTrainer true "Trainee or Trainer profile data"
// @Success 200 {object} dto.Response "Successful"
// @Failure 400 {object} string "Bad request"
// @Failure 401 {object} string "Unauthorized"
// @Router /user/profile [put]
func (c *UserController) EditProfile(ctx *fiber.Ctx) error {
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
	user, err := serve.GetUserById(uint64(userID))
	if err != nil {
		return err
	}
	if user.InfoType == "trainer" {
		var trainer dto.UserEditTraineeOrTrainer
		if err := ctx.BodyParser(&trainer); err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request payload"})
		}

		if err := utils.ValidateUser(trainer); err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"message": err})
		}
		if trainer.ActiveDays == nil {
			for i := 0; i < 7; i++ {
				trainer.ActiveDays = append(trainer.ActiveDays, true)
			}
		}
		trainer.User = dto.UserEdit{
			ID:          user.ID,
			Email:       user.Email,
			FirstName:   user.FirstName,
			LastName:    user.LastName,
			Age:         user.Age,
			Gender:      user.Gender,
			Password:    user.Password,
			PhoneNumber: user.PhoneNumber,
			InfoID:      user.InfoID,
			InfoType:    user.InfoType,
			Block:       user.Block,
			Wallet:      user.Wallet,
		}
		trainerModel, err := serve.EditTrainerProfile(uint64(userID), trainer)
		if err != nil {
			return err
		}
		return ctx.JSON(dto.Response{
			Message: "Successful",
			Success: true,
			ID:      trainerModel.ID,
		})
	}
	var trainee dto.UserEditTraineeOrTrainer
	if err := ctx.BodyParser(&trainee); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request payload"})
	}

	if err := utils.ValidateUser(trainee); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"message": err})
	}
	trainee.User = dto.UserEdit{
		ID:          user.ID,
		Email:       user.Email,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Age:         user.Age,
		Gender:      user.Gender,
		Password:    user.Password,
		PhoneNumber: user.PhoneNumber,
		InfoID:      user.InfoID,
		InfoType:    user.InfoType,
		Block:       user.Block,
		Wallet:      user.Wallet,
	}
	traineeModel, err := serve.EditTraineeProfile(uint64(userID), trainee)
	if err != nil {
		return err
	}
	return ctx.JSON(dto.Response{
		Message: "successful",
		Success: true,
		ID:      traineeModel.ID,
	})
}
