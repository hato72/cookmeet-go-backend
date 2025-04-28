package db

import (
	"os"
	"reflect"
	"testing"

	"gorm.io/gorm"
)

// func init() {
// 	// テスト環境変数の読み込み
// 	if err := godotenv.Load("../.env.test"); err != nil {
// 		panic("Error loading .env.test file")
// 	}
// }

func init() {
	// テスト用の環境変数を設定
	os.Setenv("POSTGRES_USER", "hato")
	os.Setenv("POSTGRES_PW", "hato72")
	os.Setenv("POSTGRES_DB", "hato_test")
	os.Setenv("POSTGRES_PORT", "5434")
	os.Setenv("POSTGRES_HOST", "localhost")

	// .env.testファイルが存在する場合は読み込む
	if err := godotenv.Load("../.env.test"); err != nil {
		panic("Error loading .env.test file")
	}
}
func TestNewDB(t *testing.T) {
	tests := []struct {
		name string
		want *gorm.DB
	}{
		{
			name: "Database connection successful",
			want: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewDB()

			// nilでないことを確認
			if got == nil {
				t.Errorf("NewDB() returned nil")
			}

			// *gorm.DB型であることを確認
			if reflect.TypeOf(got) != reflect.TypeOf(&gorm.DB{}) {
				t.Errorf("NewDB() did not return *gorm.DB type")
			}

			// データベース接続の確認
			sqlDB, err := got.DB()
			if err != nil {
				t.Errorf("Failed to get SQL database: %v", err)
			}

			// データベースネットワークの確認
			err = sqlDB.Ping()
			if err != nil {
				t.Errorf("Database ping failed: %v", err)
			}

			// テストデータベースを使用していることを確認
			var dbName string
			row := sqlDB.QueryRow("SELECT current_database()")
			err = row.Scan(&dbName)
			if err != nil {
				t.Errorf("Failed to get current database name: %v", err)
			}
			if dbName != os.Getenv("POSTGRES_DB") {
				t.Errorf("Expected database %s, got %s", os.Getenv("POSTGRES_DB"), dbName)
			}
		})
	}
}

func TestCloseDB(t *testing.T) {
	type args struct {
		db *gorm.DB
	}
	tests := []struct {
		name      string
		args      args
		wantPanic bool
	}{
		{
			name: "Close valid database connection",
			args: args{
				db: NewDB(), // 新しいDB接続を作成
			},
			wantPanic: false,
		},
		{
			name: "Close nil database connection",
			args: args{
				db: nil,
			},
			wantPanic: true, // パニックが発生するはず
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("CloseDB() should have panicked")
					}
				}()
			}

			CloseDB(tt.args.db)

			// 接続が閉じられていることを確認
			if tt.args.db != nil {
				sqlDB, err := tt.args.db.DB()
				if err == nil {
					err = sqlDB.Ping()
					if err == nil {
						t.Errorf("Database connection should be closed")
					}
				}
			}
		})
	}
}
