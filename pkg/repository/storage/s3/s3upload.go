package s3

import (
	"fmt"
	"io"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/google/uuid"
	"github.com/workshopapps/pictureminer.api/internal/config"
	"github.com/workshopapps/pictureminer.api/internal/constants"
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

// This checks if an s3 bucket key exists
func keyExists(key string) (bool, error) {
	svc := s3.New(s3session)
	MyBucket = config.GetConfig().S3.BucketName

	_, err := svc.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(MyBucket),
		Key:    aws.String(key),
	})

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case "NotFound":
				return false, nil
			default:
				return false, err
			}
		}
	}
	return true, nil
}

// This sets up an s3 bucket key named "default" for the client
// it also sets up a default profile picture for new users.
func DefaultProfile() (profile_url, profile_key string) {
	//func DefaultProfile() (string, string) {
	cwd, _ := os.Getwd()
	file, _ := os.Open(fmt.Sprintf("%s/static/avatar.jpg", cwd))
	defer file.Close()

	profile_url = fmt.Sprintf("https://miner-pictures.s3.amazonaws.com/%s", constants.S3_generic_avatar_key)
	profile_key = uuid.New().String()

	exists, _ := keyExists(constants.S3_generic_avatar_key)

	// if the s3 bucket key "default" doesnt exist, create one
	if !exists {
		profile_url, _ := UploadImage(file, constants.S3_generic_avatar_key)
		return profile_url, profile_key
	}

	return profile_url, profile_key
}
