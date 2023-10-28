package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

const imagesDirectory = "/home/pi/images"

func main() {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-1"), // Update this to your preferred region
	})
	if err != nil {
		log.Fatalf("Failed to create AWS session: %v", err)
	}

	svc := s3.New(sess)

	files, err := ioutil.ReadDir(imagesDirectory)
	if err != nil {
		log.Fatalf("Failed to read directory: %v", err)
	}

	for _, file := range files {
		if file.IsDir() {
			folder := file.Name()
			imageFiles, _ := ioutil.ReadDir(filepath.Join(imagesDirectory, folder))
			for _, img := range imageFiles {
				uploadToS3(svc, folder, filepath.Join(imagesDirectory, folder, img.Name()))
			}
		}
	}
}

func uploadToS3(svc *s3.S3, folder, filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("Failed to open file %s: %v", filePath, err)
		return
	}
	defer file.Close()

	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String("ratcam"),
		Key:    aws.String(filepath.Join(folder, filepath.Base(filePath))),
		Body:   file,
	})
	if err != nil {
		log.Printf("Failed to upload %s to S3: %v", filePath, err)
		return
	}

	// Delete the image from the local disk after uploading
	err = os.Remove(filePath)
	if err != nil {
		log.Printf("Failed to delete %s: %v", filePath, err)
		return
	}
	fmt.Printf("Successfully uploaded %s to S3 and deleted it locally\n", filePath)
}

