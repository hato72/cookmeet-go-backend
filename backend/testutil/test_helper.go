package testutil

import (
	"bytes"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"testing"
)

// SetupTestDirectories テストに必要なディレクトリを作成
func SetupTestDirectories(t *testing.T) {
	dirs := []string{
		"user_images/icons",
		"cuisine_images/cuisine_icons",
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("Failed to create directory %s: %v", dir, err)
		}
	}

	// テスト終了時にクリーンアップ
	t.Cleanup(func() {
		for _, dir := range dirs {
			baseDir := filepath.Dir(dir)
			os.RemoveAll(baseDir)
		}
	})
}

// CreateTestIconFile テスト用のマルチパートファイルを作成するヘルパー関数
func CreateTestIconFile(t *testing.T) *multipart.FileHeader {
	// テスト用の画像ファイルを作成
	tmpFile := CreateTempImage(t, []byte("fake image content"))

	// ファイルを開く
	file, err := os.Open(tmpFile)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	// マルチパートファイルの作成
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("icon", filepath.Base(tmpFile))
	if err != nil {
		t.Fatal(err)
	}

	// ファイルの内容を書き込み
	content := []byte("fake image content")
	if _, err := io.Copy(part, bytes.NewReader(content)); err != nil {
		t.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}

	// FileHeaderの作成
	header := &multipart.FileHeader{
		Filename: tmpFile,
		Size:     int64(len(content)),
		Header:   make(map[string][]string),
	}

	// Content-Typeの設定
	header.Header.Set("Content-Type", "image/png")

	return header
}

// CreateTempImage テスト用の一時的な画像ファイルを作成
func CreateTempImage(t *testing.T, content []byte) string {
	tmpFile, err := os.CreateTemp("", "test-*.png")
	if err != nil {
		t.Fatal(err)
	}
	defer tmpFile.Close()

	if _, err := tmpFile.Write(content); err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		os.Remove(tmpFile.Name())
	})

	return tmpFile.Name()
}

// EnsureTestDir 指定されたディレクトリが存在することを確認し、存在しない場合は作成
func EnsureTestDir(t *testing.T, dir string) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatalf("Failed to create directory %s: %v", dir, err)
	}
}

// CleanupTestDir テストディレクトリを削除
func CleanupTestDir(dir string) error {
	return os.RemoveAll(dir)
}
