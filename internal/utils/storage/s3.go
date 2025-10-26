package storage

import (
	"context"
	"errors"
	"fmt"
	"grubzo/internal/models/entity"
	"grubzo/internal/utils"
	"io"
	"net/url"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type S3FileStorage struct {
	bucket  string
	client  *s3.Client
	mutexes *utils.KeyMutex
}

func NewS3FileStorage(bucket, region, endpoint, apiKey, apiSecret string, forcePathStyle bool) (*S3FileStorage, error) {
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(apiKey, apiSecret, "")),
	)
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg, func(opt *s3.Options) {
		opt.UsePathStyle = forcePathStyle
		opt.BaseEndpoint = aws.String(endpoint)
	})

	m := &S3FileStorage{
		bucket:  bucket,
		client:  client,
		mutexes: utils.NewKeyMutex(256),
	}

	return m, nil
}

func (fs *S3FileStorage) OpenFileByKey(key string, fileType entity.FileType) (reader io.ReadSeekCloser, err error) {

	input := &s3.GetObjectInput{
		Bucket: aws.String(fs.bucket),
		Key:    aws.String(key),
	}

	file, err := fs.getObject(context.Background(), input)
	if err != nil {
		var nsk *types.NoSuchKey
		var nf *types.NotFound
		if errors.As(err, &nsk) || errors.As(err, &nf) {
			return nil, ErrFileNotFound
		}
		return nil, err
	}
	return file, nil

}

func (fs *S3FileStorage) SaveByKey(src io.Reader, key, name, contentType string, fileType entity.FileType) (err error) {
	input := &s3.PutObjectInput{
		Bucket:             aws.String(fs.bucket),
		Key:                aws.String(key),
		Body:               src,
		ContentType:        aws.String(contentType),
		ContentDisposition: aws.String(fmt.Sprintf("attachment; filename*=UTF-8''%s", url.PathEscape(name))),
	}

	uploader := manager.NewUploader(fs.client)
	_, err = uploader.Upload(context.Background(), input)
	return
}

func (fs *S3FileStorage) DeleteByKey(key string, _ entity.FileType) (err error) {
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(fs.bucket),
		Key:    aws.String(key),
	}

	_, err = fs.client.DeleteObject(context.Background(), input)
	if err != nil {
		var nsk *types.NoSuchKey
		if errors.As(err, &nsk) {
			return ErrFileNotFound
		}
		return err
	}
	return nil
}

func (fs *S3FileStorage) GenerateAccessURL(key string, fileType entity.FileType) (string, error) {
	pc := s3.NewPresignClient(fs.client)
	req, _ := pc.PresignGetObject(context.Background(), &s3.GetObjectInput{
		Bucket: aws.String(fs.bucket),
		Key:    aws.String(key),
	}, func(options *s3.PresignOptions) {
		options.Expires = 20 * time.Minute
	})

	return req.URL, nil
}

func (fs *S3FileStorage) getObject(ctx context.Context, input *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3Object, error) {
	if input == nil {
		return nil, fmt.Errorf("input is nil")
	}

	attrIn := s3.HeadObjectInput{
		Bucket: input.Bucket,
		Key:    input.Key,
	}

	attrOut, err := fs.client.HeadObject(ctx, &attrIn)
	if err != nil {
		return nil, err
	}

	objOut, err := fs.client.GetObject(ctx, input, optFns...)
	if err != nil {
		return nil, err
	}

	obj := s3Object{
		client:   fs.client,
		input:    *input,
		resp:     objOut,
		length:   *attrOut.ContentLength,
		lengthOk: true,
		body:     objOut.Body,
	}

	return &obj, nil
}

type s3Object struct {
	client     *s3.Client
	input      s3.GetObjectInput
	resp       *s3.GetObjectOutput
	length     int64
	lengthOk   bool
	body       io.ReadCloser
	pos        int64
	overSought bool
}

func (o *s3Object) Close() (err error) {
	return o.body.Close()
}

func (o *s3Object) Read(p []byte) (n int, err error) {
	if o.overSought {
		return 0, io.EOF
	}

	n, err = o.body.Read(p)
	o.pos += int64(n)
	return
}

func (o *s3Object) Seek(offset int64, whence int) (newPos int64, err error) {
	o.overSought = false

	switch whence {
	case io.SeekStart:
		newPos = offset
	case io.SeekCurrent:
		newPos = o.pos + offset
	case io.SeekEnd:
		if !o.lengthOk {
			return o.pos, fmt.Errorf("length of file unknown")
		}
		newPos = o.length + offset
		if offset >= 0 {
			o.overSought = true
			return
		}
	default:
		panic("Unknown whence")
	}

	if newPos == o.pos {
		return
	}

	err = o.Close()
	if err != nil {
		return
	}

	if newPos > 0 {
		o.input.Range = aws.String(fmt.Sprintf("bytes=%d-", newPos))
	} else {
		o.input.Range = nil
	}

	output, err := o.client.GetObject(context.Background(), &o.input)

	if err != nil {
		return
	}

	o.body = output.Body
	o.resp = output
	o.pos = newPos
	return

}
