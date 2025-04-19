package main

import (
	"backend/controller"
	"backend/db"
	"backend/model"
	"backend/repository"
	"backend/router"
	"backend/usecase"
	"backend/validator"
	"fmt"
	"log"
	"os"
)

func main() {
	//migrate

	//ローカル
	dbConn := db.NewDB()
	defer fmt.Println("Successfully Migrated")
	defer db.CloseDB(dbConn)
	dbConn.AutoMigrate(&model.User{}, &model.Cuisine{})

	db := db.NewDB()

	//本番環境
	// dbConn := db.NewPrdDB()
	// defer fmt.Println("Successfully Migrated")
	// defer db.ClosePrdDB(dbConn)
	// dbConn.AutoMigrate(&model.User{}, &model.Cuisine{})

	// db := db.NewPrdDB()

	userValidator := validator.NewUserValidator()
	cuisineValidator := validator.NewCuisineValidator()
	//taskValidator := validator.NewTaskValidator()

	userRepository := repository.NewUserRepository(db) //各コンストラクタ起動
	cuisineRepository := repository.NewCuisineRepository(db)
	//taskRepository := repository.NewTaskRepository(db)

	userUsecase := usecase.NewUserUsecase(userRepository, userValidator)
	cuisineUsecase := usecase.NewCuisineUsecase(cuisineRepository, cuisineValidator)
	//taskUsecase := usecase.NewTaskUsecase(taskRepository, taskValidator)

	userController := controller.NewUserController(userUsecase)
	cuisineController := controller.NewCuisineController(cuisineUsecase)
	//taskController := controller.NewTaskController(taskUsecase)

	//e := router.NewRouter(userController, taskController)
	e := router.NewRouter(userController, cuisineController)
	//e := router.NewRouter(userController)

	// e.Logger.Fatal(e.Start(":8080")) //サーバー起動
	// e.GET("/", func(c echo.Context) error {
	// 	return c.JSON(http.StatusOK, "hello")
	// })

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := e.Start(":" + port); err != nil {
		log.Panicf("error: %s", err)
	}

	//dockerを起動->docker compose build -> docker compose up // -> docker compose run --rm backend sh=> go run src/main.go
	// docker-compose down -v --rmi local
}
