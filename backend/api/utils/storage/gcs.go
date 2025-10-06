package storage

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)


type GCSStore struct {
	bucketName string
	client     *storage.Client
}


func NewGCSStore(ctx context.Context, bucketName string, credPath string)(*GCSStore, error){
	if bucketName == "" {
		return nil, fmt.Errorf("バケット名を指定してください")
	}

	if credPath == "" {
		return nil, fmt.Errorf("認証情報ファイルのパスを指定してください")
	}

	cli, err := storage.NewClient(ctx, option.WithCredentialsFile(credPath))
	if err != nil {
		return nil, err
	}

	checkCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
    defer cancel()
	if _, err := cli.Bucket(bucketName).Attrs(checkCtx); err != nil{
		cli.Close()
		return nil, fmt.Errorf("バケットにアクセスできません: %w", err)

	}

	return &GCSStore{
		bucketName: bucketName,
		client: cli,
	}, nil
}

func (s *GCSStore) Close() error {
	if s.client == nil { return nil}
	return s.client.Close()
}

func (s *GCSStore) PresignedPutURL(objectKey string, contentType string, ttl time.Duration) (string, error){
	if ttl <= 0 || ttl > 7*24*time.Hour {
        return "", fmt.Errorf("TTLは7日以内で指定してください")
    }
	opts := &storage.SignedURLOptions{
		Method: http.MethodPut,
		Expires: time.Now().Add(ttl),
		Scheme: storage.SigningSchemeV4,
		ContentType: contentType,
	}
	return s.client.Bucket(s.bucketName).SignedURL(objectKey, opts)
}

func (s *GCSStore) PresignedGetURL(objectKey string, ttl time.Duration) (string, error){
	if ttl <= 0 || ttl > 7*24*time.Hour {
        return "", fmt.Errorf("TTLは7日以内で指定してください")
    }
	opts := &storage.SignedURLOptions{
		Method: http.MethodGet,
		Expires: time.Now().Add(ttl),
		Scheme: storage.SigningSchemeV4,
	}
	return s.client.Bucket(s.bucketName).SignedURL(objectKey, opts)
}

func (s *GCSStore) Stat(ctx context.Context, objectKey string) (*storage.ObjectAttrs, error) {
	return s.client.Bucket(s.bucketName).Object(objectKey).Attrs(ctx)
}


