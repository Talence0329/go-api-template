package gcsmamnger

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"

	"cloud.google.com/go/storage"
)

// 透過gcloud進行默認登入(透過自己google帳號的權限) "$ gcloud auth application-default login"

var cfg Config
var ctx context.Context
var client *storage.Client

func Init(initCfg Config) error {
	cfg = initCfg
	ctx = context.Background()
	if c, err := storage.NewClient(ctx); err != nil {
		return fmt.Errorf("storage.NewClient: %v", err)
	} else {
		client = c
		return nil
	}
}

// WriteBucket : 將file 寫入 GCS
func WriteBucket(f io.Reader, filePath, contentType, cacheControl string) error {
	o := client.Bucket(cfg.BucketName).Object(filePath)
	wc := o.NewWriter(ctx)

	// 快取規則
	wc.CacheControl = cacheControl
	// 檔案contentType
	if contentType != "" {
		wc.ContentType = contentType
	}

	if _, err := io.Copy(wc, f); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %v", err)
	}

	return nil
}

// WriteBuckeByByte : 將[]byte 寫入 GCS
func WriteBuckeByByte(b []byte, filePath, contentType, cacheControl string) error {
	o := client.Bucket(cfg.BucketName).Object(filePath)
	wc := o.NewWriter(ctx)

	// 快取規則
	wc.CacheControl = cacheControl
	// 檔案contentType
	if contentType != "" {
		wc.ContentType = contentType
	}

	if _, err := wc.Write(b); err != nil {
		return fmt.Errorf("wc.Write: %v", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %v", err)
	}

	return nil
}

// ReadBucket : 讀取GCS
func ReadBucket(filePath string) ([]byte, string, error) {
	rc, err := client.Bucket(cfg.BucketName).Object(filePath).NewReader(ctx)
	if err != nil {
		return nil, "", fmt.Errorf("io.Copy: %v", err)
	}
	defer rc.Close()

	data, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, "", fmt.Errorf("ioutil.ReadAll: %v", err)
	}

	return data, rc.ContentType(), nil
}
