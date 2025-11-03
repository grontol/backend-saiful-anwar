package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"yard_plan/src/controller"
	"yard_plan/src/response"
	"yard_plan/src/service"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var (
	appHost	string
	appPort	int
	
	dbHost string
	dbUser string
	dbPass string
	dbName string
	dbPort int
)

func loadConfig() error {
	err := godotenv.Load(".env")
	if err != nil {
		return err
	}

	appHost = os.Getenv("APP_HOST")
	appPort, _ = strconv.Atoi(os.Getenv("APP_PORT"))
	
	dbHost = os.Getenv("DB_HOST")
	dbUser = os.Getenv("DB_USER")
	dbPass = os.Getenv("DB_PASS")
	dbName = os.Getenv("DB_NAME")
	dbPort, _ = strconv.Atoi(os.Getenv("DB_PORT"))
	
	return nil
}

func setupDb() (*sqlx.DB, error) {
	db, err := sqlx.Connect(
		"postgres",
		fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			dbHost,
			dbPort,
			dbUser,
			dbPass,
			dbName,
		),
	)
	if err != nil {
		return nil, err
	}
	
	return db, nil
}

func setupRoute(db *sqlx.DB) {
	app := fiber.New(fiber.Config{
		ErrorHandler: response.Error,
	})
	
	app.Use(cors.New())
	
	validate := validator.New()
	
	yardService := service.NewYardService(db)
	blockService := service.NewBlockService(db)
	placementService := service.NewPlacementService(db)
	yardPlanService := service.NewYardPlanService(db, blockService, placementService)
	
	yardController := controller.NewYardController(yardService, validate)
	blockController := controller.NewBlockController(blockService, validate)
	yardPlanController := controller.NewYardPlanController(yardPlanService, validate)
	placementController := controller.NewPlacementController(placementService, validate)
	
	yard := app.Group("/yard")
	{
		yard.Get("/", yardController.List)
		yard.Post("/", yardController.Create)
		yard.Put("/:id", yardController.Edit)
		yard.Delete("/:id", yardController.Delete)
	}
	
	block := app.Group("/block")
	{
		block.Get("/", blockController.List)
		block.Get("/by_yard/:yard_id", blockController.ListByYard)
		block.Post("/", blockController.Create)
		block.Put("/:id", blockController.Edit)
		block.Delete("/:id", blockController.Delete)
	}
	
	yardPlan := app.Group("/yard_plan")
	{
		yardPlan.Get("/", yardPlanController.List)
		yardPlan.Get("/by_yard/:yard_id", yardPlanController.ListByYard)
		yardPlan.Get("/by_block/:block_id", yardPlanController.ListByBlock)
		yardPlan.Post("/", yardPlanController.Create)
		yardPlan.Delete("/:id", yardPlanController.Delete)
	}
	
	app.Post("/suggestion", yardPlanController.Suggest)
	app.Post("/place", yardPlanController.Place)
	app.Post("/pickup", yardPlanController.Pickup)
	
	placement := app.Group("/placement")
	{
		placement.Get("/", placementController.List)
		placement.Get("/by_block/:block_id", placementController.ListByBlock)
	}
	
	app.Get("/", func (c *fiber.Ctx) error {
		return c.JSON(map[string]any{
			"status" : "Service is running",
		})
	})
	
	log.Fatal(app.Listen(fmt.Sprintf("%s:%d", appHost, appPort)))
}

func main() {
	err := loadConfig()
	
	if err != nil {
		fmt.Println("Error loading .env file.\nDetail :", err)
		return
	}
	
	db, err := setupDb()
	if err != nil {
		fmt.Println("Error connecting to database.\nDetail :", err)
		return
	}

	setupRoute(db)
}

