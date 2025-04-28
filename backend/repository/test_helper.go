package repository

import (
	"fmt"
	"log"
	"os"

	"backend/model"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func init() {
	// テスト用の環境変数を設定
	os.Setenv("POSTGRES_USER", "hato")
	os.Setenv("POSTGRES_PW", "hato72")
	os.Setenv("POSTGRES_DB", "hato_test")
	os.Setenv("POSTGRES_PORT", "5433")
	os.Setenv("POSTGRES_HOST", "localhost")

	// .env.testファイルの読み込みを試みる
	if err := godotenv.Load("../.env.test"); err != nil {
		log.Printf("Warning: .env.test file not found: %v", err)
	}
}

// SetupTestDB initializes and returns a test database connection
func SetupTestDB() *gorm.DB {
	// テスト用のDB接続情報
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PW"), os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"), os.Getenv("POSTGRES_DB"))

	// テスト用のログ設定
	config := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	db, err := gorm.Open(postgres.Open(dsn), config)
	if err != nil {
		panic(fmt.Sprintf("failed to connect database: %v", err))
	}
	log.Println("Successfully connected to test database") // ログ追加

	// テスト用のテーブルを作成
	err = db.AutoMigrate(&model.User{}, &model.Cuisine{})
	if err != nil {
		panic(fmt.Sprintf("failed to migrate database: %v", err))
	}
	log.Println("Successfully migrated database schema") // ログ追加

	return db
}

// CleanupTestDB cleans up the test database
func CleanupTestDB(db *gorm.DB) {
	// テスト用のテーブルをクリーンアップ
	err := db.Migrator().DropTable(&model.User{}, &model.Cuisine{})
	if err != nil {
		log.Printf("Warning: failed to cleanup test database: %v", err)
	}
	log.Println("Successfully cleaned up test database") // ログ追加
}

// CreateTestUser creates a test user for testing purposes
func CreateTestUser(db *gorm.DB) *model.User {
	user := &model.User{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}
	if err := db.Create(user).Error; err != nil {
		panic(fmt.Sprintf("failed to create test user: %v", err))
	}
	return user
}
