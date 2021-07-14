package awss3

import (
	"clean_code_demo/domain"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"golang.org/x/net/context"
	"io"
	"time"
)

var _ domain.StorageClient = S3StorageClient{}

type S3StorageClient struct {
	channel    chan io.Reader
	address    string
	bucket     string
	timeout    time.Duration
	client     *s3.S3
	uploader   *s3manager.Uploader
	downloader *s3manager.Downloader
}

func NewStorageClient(config AWSConfig, chanFile chan io.Reader, timeout time.Duration) (*S3StorageClient, error) {
	session, err := NewSession(config)
	if err != nil {
		return nil, err
	}

	return &S3StorageClient{
		channel:    chanFile,
		address:    config.Address,
		bucket:     config.Bucket,
		timeout:    timeout,
		client:     s3.New(session),
		uploader:   s3manager.NewUploader(session),
		downloader: s3manager.NewDownloader(session),
	}, nil
}

func (s S3StorageClient) create(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	if _, err := s.client.CreateBucketWithContext(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(s.bucket),
	}); err != nil {
		return fmt.Errorf("create: %w", err)
	}

	if err := s.client.WaitUntilBucketExists(&s3.HeadBucketInput{
		Bucket: aws.String(s.bucket),
	}); err != nil {
		return fmt.Errorf("wait: %w", err)
	}

	return nil
}

func (s S3StorageClient) Save(ctx context.Context, fileName string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	fmt.Printf("StorageClient, Upload, bucket: %s, fileName: %s\n", s.bucket, fileName)

	// check if the bucket exist, if not, create it
	exist, err := s.isBucketExist(ctx)
	if err != nil {
		return "", fmt.Errorf("save: %w", err)
	}
	if !exist {
		fmt.Printf("%s is not exist, create it\n", s.bucket)
		err = s.create(ctx)
		if err != nil {
			return "", fmt.Errorf("upload: %w", err)
		}
	}

	// get file from channel
	file, ok := <-s.channel
	if ok != true {
		return "", fmt.Errorf("save: fail to get file from channel")
	}

	res, err := s.uploader.UploadWithContext(ctx, &s3manager.UploadInput{
		Body:   file,
		Bucket: aws.String(s.bucket),
		Key:    aws.String(fileName),
	})
	if err != nil {
		return "", fmt.Errorf("upload: %w", err)
	}

	return res.Location, nil
}

func (s S3StorageClient) Download(ctx context.Context, fileName string, body io.WriterAt) error {
	if _, err := s.downloader.DownloadWithContext(ctx, body, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(fileName),
	}); err != nil {
		return fmt.Errorf("download: %w", err)
	}

	return nil
}

func (s S3StorageClient) Delete(ctx context.Context, fileName string) error {
	if _, err := s.client.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(fileName),
	}); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	if err := s.client.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(fileName),
	}); err != nil {
		return fmt.Errorf("wait: %w", err)
	}

	return nil
}

func (s S3StorageClient) List(ctx context.Context) ([]string, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	res, err := s.client.ListObjectsV2WithContext(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(s.bucket),
	})
	if err != nil {
		return nil, fmt.Errorf("list: %w", err)
	}

	result := make([]string, len(res.Contents))

	for _, object := range res.Contents {
		result = append(result, *object.Key)
	}

	return result, nil
}

// IsBucketExist check if bucket exist, if not, create it
func (s S3StorageClient) isBucketExist(ctx context.Context) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	fmt.Printf("StorageClient, isBucketExist, bucket: %s\n", s.bucket)

	result, err := s.client.ListBuckets(nil)
	if err != nil {
		fmt.Println("StorageClient, isBucketExist, err:", err)
		return false, err
	}

	for _, b := range result.Buckets {
		fmt.Printf("* %s create on %s\n", aws.StringValue(b.Name), aws.TimeValue(b.CreationDate))

		if *b.Name == s.bucket {
			fmt.Printf("%s found!\n", s.bucket)
			return true, nil
		}
	}

	return false, nil
}
