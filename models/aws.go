package models

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
	beego "github.com/beego/beego/v2/server/web"
)

var awsBucketName = "gobucket"
var awsEndpoint = "http://localhost.localstack.cloud:4566"
var awsRegion = "us-east-1"

func SetAWSInfo(bucket string, endpoint string, region string) {
	awsBucketName = bucket
	awsEndpoint = endpoint
	awsRegion = region
}

func GetAWSInfo() (string, string, string) {
	return awsBucketName, awsEndpoint, awsRegion
}

// for aws
type GenericStorage interface {
	Read(ctx context.Context, filename string) ([]byte, error)
	Write(ctx context.Context, filename string, data []byte) error
	Exists(ctx context.Context, filename string) (bool, error)
}

type S3 struct {
}

func (s S3) Read(ctx context.Context, filename string) ([]byte, error) {
	awsBucketName, awsEndpoint, awsRegion := GetAWSInfo()

	// Load the AWS configuration.
	awsCfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(awsRegion),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("dummy", "dummy", "dummy")),
	)
	if err != nil {
		return nil, fmt.Errorf("cannot load the AWS configs: %s", err)
	}

	// Create the S3 client.
	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.UsePathStyle = true
		o.BaseEndpoint = aws.String(awsEndpoint)
	})

	// Get the object from S3.
	resp, err := client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(awsBucketName), // Bucket name
		Key:    aws.String(filename),      // Key for the file in S3
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get file from S3: %s", err)
	}
	defer resp.Body.Close()

	// Read the file content into a byte slice.
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read file content: %s", err)
	}

	return buf.Bytes(), nil
}

func (s S3) Write(ctx context.Context, filename string, data []byte) error {
	awsBucketName, awsEndpoint, awsRegion := GetAWSInfo()

	// config
	awsCfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(awsRegion),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("dummy", "dummy", "dummy")),
	)
	if err != nil {
		return fmt.Errorf("cannot load the AWS configs: %s", err)
	}
	// create client
	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.UsePathStyle = true
		o.BaseEndpoint = aws.String(awsEndpoint)
	})

	fileBody := bytes.NewReader(data)

	_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(awsBucketName), // Replace with your bucket name
		Key:    aws.String(filename),      // Key for the file in S3
		Body:   fileBody,                  // File body
	})
	if err != nil {
		return fmt.Errorf("failed to upload file: %s", err)
	}

	return nil
}

func (s S3) Exists(ctx context.Context, filename string) (bool, error) {
	return true, nil
}

type Local struct {
}

func (l Local) Read(ctx context.Context, filename string) ([]byte, error) {
	filePath := fmt.Sprintf("storage/%s", filename)

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (l Local) Write(ctx context.Context, filename string, data []byte) error {
	filePath := fmt.Sprintf("storage/%s", filename)
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// write to file
	if _, err = file.Write(data); err != nil {
		return err
	}

	return nil
}

func (l Local) Exists(ctx context.Context, filename string) (bool, error) {
	return true, nil
}

func GetStorage(ctx context.Context) GenericStorage {
	storageType := beego.AppConfig.DefaultString("storagetype", "local")
	switch storageType {
	case "s3":
		return S3{}
	case "local":
		return Local{}
	}
	return nil
}
