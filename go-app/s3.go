// s3.go
//refer this to implement
//https://docs.minio.io/docs/golang-client-quickstart-guide

package main

import (
	"log"
	"fmt"
	"os"
	"strings"
	"net/http"
	"time"
	"github.com/minio/minio-go"
)

type s3Handler struct{}

func (h s3Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "S3 Page\n\n")

	endpoint := os.Getenv("S3_ENDPOINT")
	accessKeyID := os.Getenv("S3_ACCESS_KEY")
	secretAccessKey := os.Getenv("S3_SECRET_KEY")

	s3 := newClient(endpoint, accessKeyID, secretAccessKey, false)
	jp,_ := time.LoadLocation("Asia/Tokyo")

	if strings.TrimPrefix(r.URL.Path, "/s3/")=="triggeraput" {
		s3.makeBucket("testbucket","us-east-1")
		s3.putFile("testbucket","s3_upload_test_file_"+ time.Now().In(jp).Format("2006-01-02 15:04:05")+".txt","s3_upload_test_file.txt","text/plain")	
		fmt.Fprintf(w, "New object added. \n\n")
	}else {
		fmt.Fprintf(w, "Access /s3/put to add s3_upload_test_file.txt as a new object. \n\n")
	}

	fmt.Fprintf(w, "Objects2 in the testbucket...\nBrowse buckets on http://localhost:9000\n\n")

	files:=s3.listobjects("testbucket")
	for c,obkey  := range files {
		fmt.Fprintln(w,c+1,". ",obkey)
    }

}

type client struct {
	host string
	s3   *minio.Client
}
func newClient(host, key, secret string, insecure bool) *client {
	if host == "" {
		host = "s3.amazonaws.com"
		insecure = false
	}

	s3Client, err := minio.New(host, key, secret, insecure)
	if err != nil {
		log.Fatalln("minio.New", err)
	}

	return &client{
		host,
		s3Client, 
	}
}

func (c *client) bucketExists(bucket string) bool {
	_, err := c.s3.BucketExists(bucket)
	return err == nil
}

func (c *client) makeBucket(bucketName string,location string) {

	err := c.s3.MakeBucket(bucketName, location)
    if err != nil {
        // Check to see if we already own this bucket (which happens if you run this twice)
        if c.bucketExists(bucketName) {
            log.Printf("We already own %s\n", bucketName)
        } else {
            log.Fatalln(err)
        }
    }
	log.Printf("Successfully created %s\n", bucketName)
}


func (c *client) putFile(bucketName string,objectName string,filePath string,contentType string) {
	
	 // Upload the file with FPutObject
	 n, err := c.s3.FPutObject(bucketName, objectName, filePath, minio.PutObjectOptions{ContentType:contentType})
	 if err != nil {
		 log.Fatalln(err)
	 }
	
	 log.Printf("Successfully uploaded %s of size %d\n", objectName, n)
}



func (c *client) listobjects(bucketName string) []string {
	
	 // Create a done channel to control 'ListObjects' go routine.
	doneCh := make(chan struct{})

	// Indicate to our routine to exit cleanly upon return.
	defer close(doneCh)

	var files []string
	// List all objects from a bucket-name with a matching prefix.
	for object := range c.s3.ListObjects(bucketName, "/", true, doneCh) {
		if object.Err != nil {
			log.Println(object.Err)
		}
		files = append(files, object.Key)
	}
	return files
}