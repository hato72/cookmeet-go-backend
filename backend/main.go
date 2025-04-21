package main

import (
	"fmt"
	"log"
	"os"

	"backend/controller"
	"backend/model"
	"backend/repository"
	"backend/router"
	"backend/usecase"
	"backend/validator"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDB() (*gorm.DB, error) {
	// 環境変数から接続情報を取得
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PW")
	dbname := os.Getenv("POSTGRES_DB")

	// DSN を組み立て
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable "+
			"TimeZone=Asia/Tokyo prefer_simple_protocol=true",
		host, user, password, dbname, port,
	)

	db, err := gorm.Open(
		postgres.New(postgres.Config{
			DSN:                  dsn,
			PreferSimpleProtocol: true,
		}),
		&gorm.Config{PrepareStmt: false},
	)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func main() {
	// アプリケーション設定を環境変数から取得
	port := os.Getenv("PORT")
	goEnv := os.Getenv("GO_ENV")
	secret := os.Getenv("SECRET")
	feURL := os.Getenv("FE_URL")
	if port == "" {
		port = "8081"
	}

	log.Printf("env=%s, fe_url=%s, listening on :%s", goEnv, feURL, port)

	db, err := NewDB()
	if err != nil {
		log.Fatal(err)
	}
	defer log.Println("Successfully Migrated")

	// マイグレーション
	db.AutoMigrate(&model.User{}, &model.Cuisine{})

	// 以下、従来どおりの初期化…
	userValidator := validator.NewUserValidator()
	cuisineValidator := validator.NewCuisineValidator()

	userRepo := repository.NewUserRepository(db)
	cuisineRepo := repository.NewCuisineRepository(db)

	userUC := usecase.NewUserUsecase(userRepo, userValidator)
	cuisineUC := usecase.NewCuisineUsecase(cuisineRepo, cuisineValidator)

	userCtrl := controller.NewUserController(userUC)
	cuisineCtrl := controller.NewCuisineController(cuisineUC)

	e := router.NewRouter(userCtrl, cuisineCtrl)

	if err := e.Start(":" + port); err != nil {
		log.Panicf("error: %s", err)
	}
}
