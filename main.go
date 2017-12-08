package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gorilla/mux"
	"github.com/matryer/respond"
)

// type OK interface {
// 	OK() error
// }

var (
	REGION = os.Getenv("AWS_S3_REGION")
	BUCKET = os.Getenv("AWS_S3_BUCKET")
	PORT   = os.Getenv("PORT")
	SESS   = session.Must(session.NewSession(&aws.Config{Region: aws.String(REGION)}))
	// FILENAME = "test.txt"
)

type FileToUpload struct {
	FileName    string `json:"filename"`
	Author      string `json:"author"`
	Description string `json:"description"`
	SignedURL   string `json:"signedurl"`
}

func decoder(r *http.Request, v *FileToUpload) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return err
	}
	return nil
}

func SignFile() http.Handler {
	svc := s3.New(SESS)
	var incomingFile FileToUpload

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if err := decoder(r, &incomingFile); err != nil {
			log.Println("Wrong decoding")
			respond.With(w, r, http.StatusInternalServerError, err)
		}

		s3req, err := svc.PutObjectRequest(&s3.PutObjectInput{
			Bucket: aws.String(BUCKET),
			Key:    aws.String(incomingFile.FileName),
		})
		if err != nil {
			log.Println(err)
			respond.With(w, r, http.StatusInternalServerError, err)
		}
		urlStr, _ := s3req.Presign(15 * time.Minute)
		incomingFile.SignedURL = urlStr

		if err != nil {
			respond.With(w, r, http.StatusRequestTimeout, err)
		}
		respond.With(w, r, http.StatusOK, &incomingFile)
	})
}

func main() {

	if PORT == "" {
		PORT = "8080"
	}

	r := mux.NewRouter()
	r.Handle("/api/uploadfile", SignFile()).Methods("POST")
	log.Fatal(http.ListenAndServe(":"+PORT, r))

	//

	// uploader := s3manager.NewUploader(sess)

	// _, err := uploader.Upload(&s3manager.UploadInput{
	// 	Bucket: aws.String(BUCKET),
	// 	Key:    aws.String(FILENAME),
	// 	Body:   file,
	// })
	// if err != nil {
	// 	// Print the error and exit.
	// 	exitErrorf("Unable to upload %q to %q, %v", FILENAME, BUCKET, err)
	// }

	// fmt.Printf("Successfully uploaded %q to %q\n", FILENAME, BUCKET)
}
