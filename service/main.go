package service

import (
	"log"
	"mime"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

//svc is s3 service which would be used for signing stuff
var svc *s3.S3

//bucketName const is used as a bucket name for s3
const bucketName = "challengrPost"

func init() {
	config := aws.NewConfig().WithRegion("eu-west-1")
	config.WithCredentials(credentials.NewStaticCredentials("AKIAJZXILDFBECZXALOQ", "YiwmeitY+Orp4Mx8v1M+0rRHByV5ENLN2NX4QVvU", ""))
	svc = s3.New(session.New(config))
}

func PreSignS3(c *gin.Context) {
	fileName := c.Query("file-name")
	if fileName == "" {
		c.JSON(http.StatusBadRequest, []string{"Invalid file-name"})
		return
	}
	fileName = fileName + "-" + uuid.NewV4().String()

	contentType := c.Query("content-type")
	if contentType == "" {
		c.JSON(http.StatusBadRequest, []string{"Invalid content-type"})
		return
	}
	ext, err := mime.ExtensionsByType(contentType)
	if err != nil {
		log.Printf("unknown mimetype err :%v", err)
		c.JSON(http.StatusBadRequest, []string{"Invalid content-type"})
		return
	}
	log.Printf("ext: %v", ext)
	fileName = fileName + ext[0]
	//id := uuid.NewV4().String()
	req, _ := svc.PutObjectRequest(&s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(fileName),
		ContentType: aws.String(contentType),
		ACL:         aws.String("public-read"),
	})
	signedURL, headers, err := req.PresignRequest(1 * time.Minute)
	if err != nil {
		log.Printf("error: %v", err)
		c.JSON(http.StatusInternalServerError, []string{"error"})
		return
	}
	m := make(map[string]interface{})
	heads := make(map[string]string)
	for k, values := range headers {
		for _, value := range values {
			heads[k] = value
		}
	}
	m["headers"] = heads
	log.Printf("url: %v", signedURL)
	url := "https://" + bucketName + ".s3.amazonaws.com/" + fileName

	m["signedRequest"] = signedURL
	m["url"] = url
	c.JSON(http.StatusOK, m)
}
