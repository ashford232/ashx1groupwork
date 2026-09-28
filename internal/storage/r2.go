package storage

import (
	"context"
	"fmt"
	"time"

	"kolu/internal/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

type R2Storage struct {
	s3Client   *s3.Client
	bucketName string
	publicURL  string
}

func NewR2Storage(ctx context.Context, cfg *config.CloudflareConfig) (*R2Storage, error) {
	awsCfg, err := awsconfig.LoadDefaultConfig(ctx,
		awsconfig.WithRegion("auto"),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.AccessKeyID, cfg.SecretAccessKey, "")),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load aws config: %w", err)
	}
	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(fmt.Sprintf("https://%s.r2.cloudflarestorage.com", cfg.AccountID))
	})

	return &R2Storage{
		s3Client:   client,
		bucketName: cfg.BucketName,
		publicURL:  cfg.PublicURL,
	}, nil
}

func (r *R2Storage) GeneratePresignedUpload(ctx context.Context, userID string, originalFilename string, mimeType string) (uploadURL string, publicURL string, err error) {
	presignClient := s3.NewPresignClient(r.s3Client)

	objectKey := fmt.Sprintf("avatars/%s-%s-%s", userID, uuid.New().String(), originalFilename)

	req, err := presignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(r.bucketName),
		Key:         aws.String(objectKey),
		ContentType: aws.String(mimeType),
	}, s3.WithPresignExpires(15*time.Minute))
	if err != nil {
		return "", "", fmt.Errorf("failed to presign put object: %w", err)
	}

	finalPublicURL := fmt.Sprintf("%s/%s", r.publicURL, objectKey)
	return req.URL, finalPublicURL, nil
}
