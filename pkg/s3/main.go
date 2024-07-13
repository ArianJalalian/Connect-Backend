package s3

import (
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func Download(src string, dest string) {

}

func Upload(src string, path string) error {
	uploader := s3manager.NewUploader(createSession())

	f, err := os.Open(src)

	if err != nil {
		return fmt.Errorf("failed to open file %q, %v", src, err)
	}

	// Upload the file to S3.
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String("mass-media-core"),
		Key:    aws.String(path),
		Body:   f,
	})

	if err != nil {
		return fmt.Errorf("failed to upload file, %v", err)
	}

	return nil
}

func TempUrl(path string) (string, error) {
	svc := s3.New(createSession())

	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String("myBucket"),
		Key:    aws.String(path),
	})
	urlStr, err := req.Presign(15 * time.Minute)

	if err != nil {
		return "", fmt.Errorf("Failed to sign request, %v", err)
	}

	return urlStr, nil
}

func createSession() *session.Session {

	return session.Must(session.NewSession(&aws.Config{}))
}
