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

func main() {
	// sessionの作成
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("input No.")
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

		bucketName := os.Getenv("BUCKET_NAME")
		objectKey := os.Getenv("OBJECT_KEY")

		// Downloaderを作成し、S3オブジェクトをダウンロード
		downloader := s3manager.NewDownloader(sess)
		n, err := downloader.Download(f, &s3.GetObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(objectKey),
		})
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("download compleated: %d byte", n)
	default:
		log.Printf("comman is not found")
	}

}
