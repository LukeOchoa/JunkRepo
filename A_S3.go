package main

//import (
//	"fmt"
//	"os"
//	"strings"
//	"io/ioutil"
//
//	"github.com/aws/aws-sdk-go/aws"
//	"github.com/aws/aws-sdk-go/aws/session"
//	"github.com/aws/aws-sdk-go/service/s3"
//)
//
//const (
//	BUCKET_NAME = "borkcraftbucket"
//	REGION = "us-west-1"
//)
//
//var s3session *s3.S3
//
//func init() {
//	s3session = s3.New(session.Must(session.NewSession(&aws.Config{
//		Region: aws.String(REGION)
//	})))
//}
//
//func exitErrorf(msg string, args ...interface{}) {
//    fmt.Fprintf(os.Stderr, msg+"\n", args...)
//    os.Exit(1)
//}
//
////func uploadObject(filename string) (resp *s3.PutObjectOutput) {
////	f, err := os.Open(filename)
////}
//
//func getObject(filename string) () {
//	fmt.Println("Downloading: ", filename)
//  
//	resp, err := s3session.GetObject(&s3.GetObjectInput{
//	  Bucket:aws.String(BUCKET_NAME),
//	  Key: aws.String(filename),
//	})
//  
//	if err != nil {
//	  panic(err)
//	}
//  
//	body, err := ioutil.ReadAll(resp.Body)
//	err = ioutil.WriteFile(filename, body, 0644)
//	if err != nil {
//	  panic(err)
//	}
//}





