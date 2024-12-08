package storage

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"strings"
	"time"

	"github.com/adiubaidah/rfid-syafiiyah/pkg/util"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go"
)

type StorageManager struct {
	client     *s3.Client
	BucketName string
	Region     string
}

func NewStorageManager(client *s3.Client, BucketName string, Region string) *StorageManager {
	return &StorageManager{
		client:     client,
		BucketName: BucketName,
		Region:     Region,
	}
}

func (s *StorageManager) UploadFile(ctx context.Context, entity *multipart.FileHeader, fileName string) (string, error) {
	file, err := entity.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	checksum := util.CalculateSHA256(fileBytes)

	putObjectInput := &s3.PutObjectInput{
		Bucket:         aws.String(s.BucketName),
		Key:            aws.String(fileName),
		Body:           bytes.NewReader(fileBytes),
		ChecksumSHA256: aws.String(checksum),
	}

	_, err = s.client.PutObject(ctx, putObjectInput)
	if err != nil {
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) && apiErr.ErrorCode() == "EntityTooLarge" {
			return "", fmt.Errorf("file %v is too large to upload to %v:%v", fileName, s.BucketName, fileName)
		}
		return "", err
	}

	waiter := s3.NewObjectExistsWaiter(s.client)
	headObjectInput := &s3.HeadObjectInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(fileName),
	}
	err = waiter.Wait(ctx, headObjectInput, time.Minute)
	if err != nil {
		return "", err
	}

	objectURL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s.BucketName, s.Region, fileName)

	return objectURL, nil
}

func (s *StorageManager) DeleteFile(ctx context.Context, fileName string) error {

	stringObj := strings.Split(fileName, "/")[3]

	deleteObjectInput := &s3.DeleteObjectInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(stringObj),
	}

	_, err := s.client.DeleteObject(ctx, deleteObjectInput)
	if err != nil {
		return err
	}

	return nil
}
