package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	_ "github.com/aws/aws-sdk-go/aws"
	_ "github.com/aws/aws-sdk-go/aws/session"
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
	// FILENAME = "test.txt"
)

type FileToUpload struct {
	FileName    string `json:"filename"`
	Author      string `json:"author"`
	Description string `json:"description"`
	SignedURL   string `json:"-"`
}

// func (f *FileToUpload) OK() error {
// 	if f.SignedURL == nil {
// 		err := errors.New("No Signed URL returned")
// 		return err
// 	}
// 	return nil
// }
func decoder(r *http.Request, v *FileToUpload) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return err
	}
	return nil
}

func handleIncoming(w http.ResponseWriter, r *http.Request) {
	var incomingFile FileToUpload
	if err := decoder(r, &incomingFile); err != nil {
		respond.With(w, r, http.StatusInternalServerError, err)
	}
	respond.With(w, r, http.StatusOK, &incomingFile)
}

func main() {

	if PORT == "" {
		PORT = "8080"
	}

	r := mux.NewRouter()
	r.HandleFunc("/api/uploadfile", handleIncoming).Methods("POST")
	log.Fatal(http.ListenAndServe(":"+PORT, r))

	//sess := session.Must(session.NewSession(&aws.Config{Region: aws.String(REGION)}))

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
