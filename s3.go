package gbdx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// s3Credentials are the S3 credentials obtained via the GBDX s3creds API.
type s3Credentials struct {
	S3SecretKey    string `json:"S3_secret_key"`
	Prefix         string `json:"prefix"`
	Bucket         string `json:"bucket"`
	S3AccessKey    string `json:"S3_access_key"`
	S3SessionToken string `json:"S3_session_token"`
}

// getAwsS3Client returns an aws-sdk-go/service/s3 client using
// temporrary S3 credentials from GBDX API.
func getAwsS3Client(client *http.Client, awsRegion string) (*s3.S3, string, string, error) {
	tmpCreds := s3Credentials{}
	url := "https://geobigdata.io/s3creds/v1/prefix"
	response, err := client.Get(url)
	if err != nil {
		return nil, "", "", fmt.Errorf("HTTP GET %s: %v", url, err)
	}
	if response.Status != "200 OK" {
		var byteSlice []byte
		response.Request.Body.Read(byteSlice)
		return nil, "", "", fmt.Errorf("HTTP POST %s;  returned status %s; request body %q", url, response.Status, byteSlice)
	}

	if err = json.NewDecoder(response.Body).Decode(&tmpCreds); err != nil {
		return nil, "", "", fmt.Errorf("Decoding search response %q: %v", response.Body, err)
	}

	creds := credentials.NewStaticCredentials(tmpCreds.S3AccessKey, tmpCreds.S3SecretKey, tmpCreds.S3SessionToken)
	awsSession := session.New(&aws.Config{
		Region:      &awsRegion,
		Credentials: creds,
	})

	s3Client := s3.New(awsSession)

	return s3Client, tmpCreds.Bucket, tmpCreds.Prefix, nil
}

// ListBucket lists the contents of the customer bucket.
// TODO: Add optional "prefix string" argument which will drill down into subdirectories.
func (a *Api) ListBucket() (contents string, err error) {
	s3Client, bucket, gbdxPrefix, err := getAwsS3Client(a.client, "us-east-1")
	if err != nil {
		return "", fmt.Errorf("getAwsClient(client, %s): %v", gbdxPrefix, err)
	}
	params := &s3.ListObjectsInput{
		Bucket:    &bucket,
		Delimiter: aws.String("/"),
		Prefix:    aws.String(fmt.Sprintf("%s/", gbdxPrefix)),
	}

	output, err := s3Client.ListObjects(params)
	if err != nil {
		return "", fmt.Errorf("s3Client.ListObjects(%q)b: %v", params, err)
	}

	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("Owner\tSize\tLastModified\tKey\n"))
	for _, obj := range output.Contents {
		buffer.WriteString(fmt.Sprintf("%5v-t%5v\t%10v\t%v\n", *obj.Owner, *obj.Size, *obj.LastModified, *obj.Key))
	}
	for _, obj := range output.CommonPrefixes {
		buffer.WriteString(fmt.Sprintf("%5v\t%5v\t%10v\t%v\n", "--", "--", "--", *obj.Prefix))
	}
	return buffer.String(), nil
}
