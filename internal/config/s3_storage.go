package config

type S3StorageConfiguration struct {
	AWSRegion          string
	AWSAccessKeyID     string
	AWSSecretAccessKey string
	BucketName         string
}
