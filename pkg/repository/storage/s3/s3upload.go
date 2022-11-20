package s3

import (
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/workshopapps/pictureminer.api/internal/config"
)

var AccessKeyID string
var SecretAccessKey string
var MyRegion string
var MyBucket string
var filepath string

var (
	s3session *session.Session
)

func Session() (session *session.Session) {
	return s3session
}

func ConnectAws() *session.Session {
	var err error
	AccessKeyID = config.GetConfig().S3.AWSAccessKeyID
	SecretAccessKey = config.GetConfig().S3.AWSSecretAccessKey
	MyRegion = config.GetConfig().S3.AWSRegion
	s3session, err = session.NewSession(
		&aws.Config{
			Region: aws.String(MyRegion),
			Credentials: credentials.NewStaticCredentials(
				AccessKeyID,
				SecretAccessKey,
				"", // a token will be created when the session it's used.
			),
		})
	if err != nil {
		panic(err)
	}

	return s3session
}

func UploadImage(file io.ReadCloser, filename string) (string, error) {
	uploader := s3manager.NewUploader(s3session)
	MyBucket = config.GetConfig().S3.BucketName

	//upload to the s3 bucket
	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(MyBucket),
		// ACL:    aws.String("public-read"),
		Key:  aws.String(filename),
		Body: file,
	})

	if err != nil {
		return "", err
	}
	filepath = "https://" + MyBucket + ".s3.amazonaws.com/" + filename

	return filepath, nil
}
