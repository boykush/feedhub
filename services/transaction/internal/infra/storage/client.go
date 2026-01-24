package storage

import (
	"context"
	"io"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// Client provides access to object storage (S3-compatible)
type Client struct {
	s3Client *s3.Client
	bucket   string
}

// NewClient creates a new storage client from environment variables
func NewClient(ctx context.Context) (*Client, error) {
	endpoint := os.Getenv("STORAGE_ENDPOINT")
	bucket := os.Getenv("STORAGE_BUCKET")
	accessKey := os.Getenv("STORAGE_ACCESS_KEY")
	secretKey := os.Getenv("STORAGE_SECRET_KEY")

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
		config.WithRegion("us-east-1"), // MinIO doesn't care about region
	)
	if err != nil {
		return nil, err
	}

	s3Client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(endpoint)
		o.UsePathStyle = true // Required for MinIO
	})

	return &Client{
		s3Client: s3Client,
		bucket:   bucket,
	}, nil
}

// GetObject retrieves an object from storage
func (c *Client) GetObject(ctx context.Context, key string) (io.ReadCloser, error) {
	output, err := c.s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}
	return output.Body, nil
}
