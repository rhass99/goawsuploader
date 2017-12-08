package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

var (
	REGION = os.Getenv("AWS_S3_REGION")
	BUCKET = os.Getenv("AWS_S3_BUCKET")
	// FILENAME = "test.txt"
)

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

func main() {

	sess := session.Must(session.NewSession(&aws.Config{Region: aws.String(REGION)}))
	uploader := s3manager.NewUploader(sess)

	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(BUCKET),
		Key:    aws.String(FILENAME),
		Body:   file,
	})
	if err != nil {
		// Print the error and exit.
		exitErrorf("Unable to upload %q to %q, %v", FILENAME, BUCKET, err)
	}

	fmt.Printf("Successfully uploaded %q to %q\n", FILENAME, BUCKET)
}
