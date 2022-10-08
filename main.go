package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/oklog/ulid/v2"
)

var (
	bucketName    = os.Getenv("BUCKET_NAME")
	upObjectKey   = os.Getenv("UPLOAD_OBJECT_KEY")
	downObjectKey = os.Getenv("DOWNLOAD_OBJECT_KEY")
)

func main() {
	// sessionの作成
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("please input (1: download 2: upload)")
	scanner.Scan()
	in := scanner.Text()

	switch in {
	case "1":
		// S3オブジェクトを書き込むファイルの作成
		fileName := ulid.Make().String()
		f, err := os.Create("./download/" + fileName + ".txt")
		if err != nil {
			log.Fatal(err)
		}
		// Downloaderを作成し、S3オブジェクトをダウンロード
		downloader := s3manager.NewDownloader(sess)
		n, err := downloader.Download(f, &s3.GetObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(downObjectKey),
		})
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("download compleated: %d byte", n)

	case "2":
		// ファイルを開く
		targetFilePath := "./upload/upload.txt"
		file, err := os.Open(targetFilePath)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		// Uploaderを作成し、ローカルファイルをアップロード
		uploader := s3manager.NewUploader(sess)
		_, err = uploader.Upload(&s3manager.UploadInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(upObjectKey),
			Body:   file,
		})
		if err != nil {
			log.Fatal(err)
		}
		log.Println("upload done")
	default:
		log.Printf("command is not found")
	}

}
