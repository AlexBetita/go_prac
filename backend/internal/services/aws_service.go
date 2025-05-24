package services

import (
    "context"
    "fmt"
    "io"
    "os"
    "time"

    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/credentials"
    "github.com/aws/aws-sdk-go-v2/service/s3"
    "github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type AWSService struct {
    client *s3.Client
    bucket string
}

func NewAWSService(ctx context.Context) (*AWSService, error) {
    // Load AWS config from env
    awsCfg, err := config.LoadDefaultConfig(ctx,
        config.WithRegion(os.Getenv("AWS_S3_REGION")),
        config.WithCredentialsProvider(
            credentials.NewStaticCredentialsProvider(
                os.Getenv("AWS_S3_ACCESS_KEY"),
                os.Getenv("AWS_S3_SECRET_ACCESS_KEY"),
                ""),
        ),
    )
    if err != nil {
        return nil, fmt.Errorf("loading aws config: %w", err)
    }

    return &AWSService{
        client: s3.NewFromConfig(awsCfg),
        bucket: os.Getenv("AWS_S3_BUCKET"),
    }, nil
}

// UploadFile uploads an io.Reader to S3 and returns the public URL.
func (s *AWSService) UploadFile(ctx context.Context, filename string, body io.Reader, contentType string) (string, error) {
    key := fmt.Sprintf("uploads/%d-%s", time.Now().UnixNano(), filename)
    _, err := s.client.PutObject(ctx, &s3.PutObjectInput{
        Bucket:      aws.String(s.bucket),
        Key:         aws.String(key),
        Body:        body,
        ContentType: aws.String(contentType),
        ACL:         types.ObjectCannedACLPublicRead,
    })
    if err != nil {
        return "", fmt.Errorf("put object: %w", err)
    }

    // Construct URL (works for most public buckets)
    url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s",
        s.bucket,
        os.Getenv("AWS_S3_REGION"),
        key,
    )
    return url, nil
}
