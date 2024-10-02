package storage

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"io"
	"net/url"
	"strings"
)

type S3 struct {
	session *session.Session
	bucket  string
	host    string
}

type Upload struct {
	File        io.ReadSeeker
	Filename    string
	Size        int64
	ContentType string
}

func NewS3(sess *session.Session, bucket string) Storage {
	return &S3{
		session: sess,
		bucket:  bucket,
		host:    strings.Trim(*sess.Config.Endpoint, "/"),
	}
}

func NewS3FromUrl(dsn string) (Storage, error) {
	u, err := url.Parse(dsn)
	if err != nil {
		return nil, err
	}

	password, _ := u.User.Password()

	sess, err := session.NewSession(&aws.Config{
		Credentials:      credentials.NewStaticCredentials(u.User.Username(), password, ""),
		S3ForcePathStyle: aws.Bool(true),
		Endpoint:         aws.String(u.Scheme + "//" + u.Host),
		Region:           aws.String("ru-1"),
	})
	if err != nil {
		return nil, err
	}

	return NewS3(sess, strings.Trim(u.Path, "/")), nil
}

func (s *S3) UploadFile(key string, upload Upload) error {
	buf, err := io.ReadAll(upload.File)
	if err != nil {
		return err
	}

	_, err = s3.New(s.session).PutObject(&s3.PutObjectInput{
		Bucket:        aws.String(s.bucket),
		Key:           aws.String(key),
		Body:          bytes.NewReader(buf),
		ContentType:   aws.String(upload.ContentType),
		ContentLength: aws.Int64(int64(len(buf))),
	})

	return err
}

func (s *S3) URL(path string) string {
	return fmt.Sprintf("%s/%s/%s", s.host, s.bucket, strings.Trim(path, "/"))
}
