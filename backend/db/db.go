package db

//dbへの接続

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDB() *gorm.DB {
	// if os.Getenv("GO_ENV") == "dev" {
	// 	err := godotenv.Load()
	// 	if err != nil {
	// 		log.Fatalln(err)
	// 	}
	// }
	// if os.Getenv("GO_ENV") == "dev" {
	// 	err := godotenv.Load(fmt.Sprintf(".env.%s", os.Getenv("GO_ENV")))
	// 	if err != nil {
	// 		log.Fatalln(err)
	// 	}
	// }

	//ローカルの場合は以下のコメントアウトを外す
	// err := godotenv.Load(fmt.Sprintf(".env.dev"))
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PW"), os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"), os.Getenv("POSTGRES_DB"))

	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Connceted")
	return db
}

func CloseDB(db *gorm.DB) {
	sqlDB, _ := db.DB()
	if err := sqlDB.Close(); err != nil {
		log.Fatalln(err)
	}
}
