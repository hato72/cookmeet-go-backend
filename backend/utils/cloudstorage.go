package utils

import (
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/storage"
)

// UploadToCloudStorage はファイルを GCS にアップロードし、公開URLを返す
func UploadToCloudStorage(bucketName, objectName string, file io.Reader) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel() // コンテキストをキャンセルして確実にリソースを解放

	// credentialFilePath := "/etc/secrets/cookmeet-backend.json"
	// client, err := storage.NewClient(ctx, option.WithCredentialsFile(credentialFilePath))

	client, err := storage.NewClient(ctx)
	if err != nil {
		fmt.Printf("storage client creation error: %v\n", err)
		return "", fmt.Errorf("failed to create storage client: %v", err)
	}
	defer client.Close()

	// GCS にファイルをアップロード
	bucket := client.Bucket(bucketName)
	obj := bucket.Object(objectName)
	fmt.Printf("bucket: %s, object: %s\n", bucketName, objectName)
	w := obj.NewWriter(ctx)
	w.ContentType = "image/jpeg" // 必要に応じて変更
	w.CacheControl = "public, max-age=86400"
	w.PredefinedACL = "" // ACLを無効化

	if _, copyErr := io.Copy(w, file); copyErr != nil {
		fmt.Printf("file upload error: %v\n", copyErr)
		return "", fmt.Errorf("failed to write file to cloud storage: %v", copyErr)
	}

	if closeErr := w.Close(); closeErr != nil {
		fmt.Printf("writer close error: %v\n", closeErr)
		return "", fmt.Errorf("failed to close writer: %v", closeErr)
	}

	// 公開URLを生成
	// publicURL := fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucketName, objectName)
	// publicURL := fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucketName, objectName)

	publicURL, err := generateSignedURL(bucket, objectName)
	if err != nil {
		return "", fmt.Errorf("failed to generate signed URL: %v", err)
	}
	return publicURL, nil

}

// 署名付きURLを生成する関数
func generateSignedURL(bucket *storage.BucketHandle, objectName string) (string, error) {
	// Cloud Storageクライアントから署名付きURLを生成
	opts := &storage.SignedURLOptions{
		Method:  "GET",
		Expires: time.Now().Add(24 * time.Hour * 7), // 1週間有効
		Scheme:  storage.SigningSchemeV4,            // V4署名を使用
	}

	// 署名付きURLを生成
	url, err := bucket.SignedURL(objectName, opts)
	if err != nil {
		fmt.Printf("signed URL generation error: %v\n", err)
		return "", fmt.Errorf("failed to generate signed URL (object: %s): %v",
			objectName, err)
	}

	return url, nil
}
