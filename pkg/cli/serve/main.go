package serve

import (
	"os"

	_ "bitbucket.org/dyfrag-internal/mass-media-core/docs"
	serve "bitbucket.org/dyfrag-internal/mass-media-core/pkg/cli/serve/controller"
	"bitbucket.org/dyfrag-internal/mass-media-core/pkg/configs"
	"bitbucket.org/dyfrag-internal/mass-media-core/pkg/database"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/spf13/cobra"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

func Main() {
	app := fiber.New()
	initialization()

	app.Use(cors.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	user := app.Group("/user")
	var userController serve.UserController
	userController.RegisterRoutes(user)

	trainee := app.Group("/trainee")
	var traineeController serve.TraineeController
	traineeController.RegisterRoutes(trainee)

	trainer := app.Group("/trainer")
	var trainerController serve.TrainerController
	trainerController.RegisterRoutes(trainer)

	admin := app.Group("/admin")
	var adminController serve.AdminController
	adminController.RegisterRoutes(admin)

	port := os.Getenv("SERVER_PORT")
	connection := ":" + port
	app.Listen(connection)
}

func initialization() {
	configs.SetUpConfigs()
	database.SetUpDB()
}

func New() *cobra.Command {

	return &cobra.Command{
		Use:   "serve",
		Short: "runs http server",
		Run: func(cmd *cobra.Command, args []string) {
			Main()
		},
	}
}
