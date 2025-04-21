package db

//dbへの接続

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
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
	//err := godotenv.Load(fmt.Sprintf("C:/Users/hatot/.vscode/go_backend_hackathon/backend/.env.dev"))
	if _, err := os.Stat(".env.dev"); err == nil {
		if err := godotenv.Load(".env.dev"); err != nil {
			log.Fatalln(err)
		}
	} else {
		log.Println(".env.dev not found, skip loading local env file")
	}

	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PW"), os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"), os.Getenv("POSTGRES_DB"))

	// GORMの設定を追加
	gormConfig := &gorm.Config{
		PrepareStmt: false, // プリペアドステートメントを無効化
	}

	db, err := gorm.Open(postgres.Open(url), gormConfig)
	if err != nil {
		log.Fatalln(err)
	}

	// 接続プールの設定
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalln(err)
	}

	// 接続プールの最大数を設定
	sqlDB.SetMaxOpenConns(10)  // 同時に開くことができる接続の最大数
	sqlDB.SetMaxIdleConns(5)   // アイドル状態で保持する接続の最大数
	sqlDB.SetConnMaxLifetime(time.Hour) // 接続の最大寿命

	fmt.Println("Connected")
	return db
}

func CloseDB(db *gorm.DB) {
	sqlDB, _ := db.DB()
	if err := sqlDB.Close(); err != nil {
		log.Fatalln(err)
	}
}
