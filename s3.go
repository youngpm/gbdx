package gbdx

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// s3Info stores temporary credentials, bucket and customer prefix for
// data stored in AWS S3 by GBDX.  The values should be obtained via
// the GBDX s3creds API.
type s3Info struct {
	S3SecretKey    string `json:"S3_secret_key"`
	Prefix         string `json:"prefix"`
	Bucket         string `json:"bucket"`
	S3AccessKey    string `json:"S3_access_key"`
	S3SessionToken string `json:"S3_session_token"`
}

// getS3Info returns an s3Info struct that defines temporary S3
// credentials, customer bucket & customer prefix via the GBDX s3Creds
// api. https://gbdxdocs.digitalglobe.com/docs/s3-storage-service-course
func getS3Info(client *http.Client) (tmpCreds s3Info, err error) {
	url := "https://geobigdata.io/s3creds/v1/prefix"
	response, err := client.Get(url)
	if err != nil {
		return tmpCreds, fmt.Errorf("HTTP GET %s: %v", url, err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		var byteSlice []byte
		response.Request.Body.Read(byteSlice)
		return tmpCreds, fmt.Errorf("HTTP POST %s;  returned status %s, expected status %s; request body %q", url, response.Status, http.StatusOK, byteSlice)
	}

	if err = json.NewDecoder(response.Body).Decode(&tmpCreds); err != nil {
		return tmpCreds, fmt.Errorf("Decoding search response %q: %v", response.Body, err)
	}

	return tmpCreds, nil
}

// getAwsS3Client returns an S3 client with the credentials stored in s3Info
func getAwsS3Client(tmpCreds s3Info, awsRegion string) *s3.S3 {
	creds := credentials.NewStaticCredentials(tmpCreds.S3AccessKey, tmpCreds.S3SecretKey, tmpCreds.S3SessionToken)

	// Return an AWS S3 client with the provided credentials
	awsSession := session.New(&aws.Config{
		Region:      &awsRegion,
		Credentials: creds,
	})

	s3Client := s3.New(awsSession)

	return s3Client
}

// ListBucket writes a listing of contents of the GBDX customer bucket
// to w.
func (a *Api) ListBucket(requestedPrefix string, w io.Writer) (err error) {
	tmpCreds, err := getS3Info(a.client)
	if err != nil {
		return fmt.Errorf("getS3Info(%v): %v", a.client, err)
	}
	s3Client := getAwsS3Client(tmpCreds, "us-east-1")

	var rootPrefix string
	// Note we dont use path.Join as it will gobble up any
	// trailing "/"
	if utf8.RuneCountInString(requestedPrefix) == 0 {
		rootPrefix = fmt.Sprintf("%s/", tmpCreds.Prefix)
	} else if strings.Contains(requestedPrefix, tmpCreds.Prefix) {
		rootPrefix = requestedPrefix
	} else if requestedPrefix[len(requestedPrefix)-1] == '/' {
		rootPrefix = fmt.Sprintf("%s/%s", tmpCreds.Prefix, requestedPrefix)
	} else {
		rootPrefix = fmt.Sprintf("%s/%s/", tmpCreds.Prefix, requestedPrefix)
	}

	inputParams := &s3.ListObjectsInput{
		Bucket:    &tmpCreds.Bucket,
		Delimiter: aws.String("/"),
		Prefix:    &rootPrefix,
	}

	// ListObjects returns a ListObjectsOutput struct
	listObjectsOutput, err := s3Client.ListObjects(inputParams)
	if err != nil {
		return fmt.Errorf("s3Client.ListObjects(%q)b: %v", inputParams, err)
	}

	//_, err = io.WriteString(w, fmt.Sprintf("%10v %30v %v\n", "size", "LastModified", "Key"))
	//if err != nil {
	//	return fmt.Errorf("Writing header: %v", err)
	//}

	// CommonPrefixes is a slice of "directories" found in the bucket.
	for _, obj := range listObjectsOutput.CommonPrefixes {
		_, err = io.WriteString(w, fmt.Sprintf("%-20v %10v PRE %v\n", "", "", (*obj.Prefix)[len(rootPrefix):]))
		if err != nil {
			return fmt.Errorf("Writing listObjectsOutput.CommonPrefixes: %v", err)
		}
	}

	// Contents is a slice of files found in the bucket.
	for _, obj := range listObjectsOutput.Contents {
		_, err = io.WriteString(w, fmt.Sprintf("%-20v %10v %v\n", (*obj.LastModified).Format("2006-01-02 15:04:05"), *obj.Size, *obj.Key))
		if err != nil {
			return fmt.Errorf("Writing listObjectsOutput.Contents: %v", err)
		}
	}
	return nil
}
