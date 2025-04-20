package db

//dbへの接続

import (
	"fmt"
	"log"
	"os"
	"strings"

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

	
	// if _, err := os.Stat(".env.dev"); err == nil {
	// 	if err := godotenv.Load(".env.dev"); err != nil {
	// 		log.Fatalln(err)
	// 	}
	// } else {
	// 	log.Println(".env.dev not found, skip loading local env file")
	// }
	
	secretContent, err := os.ReadFile("/env/secret")
    if err != nil {
        log.Printf("環境変数ファイルの読み込みに失敗しました: %v", err)
    } else {
        log.Println("環境変数ファイルを読み込みました")
        // 改行で分割
        envVars := strings.Split(string(secretContent), "\n")
        for _, envVar := range envVars {
            if envVar == "" || strings.HasPrefix(envVar, "#") {
                continue // 空行やコメント行はスキップ
            }
            // 環境変数の名前と値を取得
            parts := strings.SplitN(envVar, "=", 2)
            if len(parts) == 2 {
                key := parts[0]
                value := parts[1]
                // 既に設定されていない場合のみ設定
                if os.Getenv(key) == "" {
                    os.Setenv(key, value)
                    log.Printf("環境変数を設定しました: %s", key)
                }
            }
        }
    }

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
