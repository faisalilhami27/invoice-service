package gcs

import (
	"bytes"

	"cloud.google.com/go/storage"

	"context"
	"encoding/json"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"google.golang.org/api/option"

	"io"
	"time"
)

const (
	coldLineStorageClass = "COLDLINE"
)

//nolint:revive
type GCSPackage struct {
	Client                 *storage.Client
	ServiceAccountKeyJSON  ServiceAccountKeyJSON
	SignedURLTimeInMinutes uint
	BucketName             string
	TimeoutInSeconds       uint
}

type ServiceAccountKeyJSON struct {
	Type                    string `json:"type"`
	ProjectId               string `json:"project_id"`     //nolint:revive,stylecheck
	PrivateKeyId            string `json:"private_key_id"` //nolint:revive,stylecheck
	PrivateKey              string `json:"private_key"`
	ClientEmail             string `json:"client_email"`
	ClientId                string `json:"client_id"`                   //nolint:revive,stylecheck
	AuthUri                 string `json:"auth_uri"`                    //nolint:revive,stylecheck
	TokenUri                string `json:"token_uri"`                   //nolint:revive,stylecheck
	AuthProviderX509CertUrl string `json:"auth_provider_x509_cert_url"` //nolint:revive,stylecheck
	ClientX509CertUrl       string `json:"client_x509_cert_url"`        //nolint:revive,stylecheck
	UniverseDomain          string `json:"universe_domain"`             //nolint:revive,stylecheck
}

type UploadOptions struct {
	Object  string
	Timeout time.Duration
	File    io.Reader
	Public  bool
}

type Option func(*GCSPackage)

type IGCSClient interface {
	UploadFileInByte(ctx context.Context, fileName string, data []byte) (string, error)
	GetSignedURL(context.Context, string, time.Duration) (string, error)
	Upload(context.Context, *UploadOptions) error
	Delete(ctx context.Context, object string, hardDelete bool, timeout time.Duration) error
}

func withGCSClient(gcsClient *storage.Client) Option {
	return func(s *GCSPackage) {
		s.Client = gcsClient
	}
}

func WithServiceAccountKeyJSON(serviceAccountKeyJSON ServiceAccountKeyJSON) Option {
	return func(s *GCSPackage) {
		s.ServiceAccountKeyJSON = serviceAccountKeyJSON
	}
}

func WithSignedURLTimeInMinutes(signedURLTimeInMinutes uint) Option {
	return func(s *GCSPackage) {
		s.SignedURLTimeInMinutes = signedURLTimeInMinutes
	}
}

func WithBucketName(bucketName string) Option {
	return func(s *GCSPackage) {
		s.BucketName = bucketName
	}
}

func WithTimeoutInSeconds(timeoutInSeconds uint) Option {
	return func(s *GCSPackage) {
		s.TimeoutInSeconds = timeoutInSeconds
	}
}

func NewGCSClient(opts ...Option) *GCSPackage {
	s := &GCSPackage{}
	for _, opt := range opts {
		opt(s)
	}

	ctx := context.Background()
	client, err := s.createClient(ctx)
	if err != nil {
		panic(err)
	}

	optionGCS := withGCSClient(client)
	optionGCS(s)
	return s
}

func (c *GCSPackage) createClient(ctx context.Context) (*storage.Client, error) {
	reqBodyBytes := new(bytes.Buffer)
	err := json.NewEncoder(reqBodyBytes).Encode(c.ServiceAccountKeyJSON)
	if err != nil {
		log.Errorf("an error occurred when encode invoice account key json : %v", err)
		return nil, err
	}

	jsonByte := reqBodyBytes.Bytes()
	client, err := storage.NewClient(ctx, option.WithCredentialsJSON(jsonByte))
	if err != nil {
		log.Errorf("an error occurred when create gcs client : %v", err)
		return nil, err
	}

	return client, nil
}

func (c *GCSPackage) GetSignedURL(ctx context.Context, object string, timeout time.Duration) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout) //nolint:ineffassign,staticcheck
	defer cancel()

	bucket := c.Client.Bucket(c.BucketName)
	url, err := bucket.SignedURL(object, &storage.SignedURLOptions{
		GoogleAccessID: c.ServiceAccountKeyJSON.ClientEmail,
		PrivateKey:     []byte(c.ServiceAccountKeyJSON.PrivateKey),
		Method:         "GET",
		Expires:        time.Now().Add(time.Duration(c.SignedURLTimeInMinutes) * time.Minute), //nolint:gosec
	})
	if err != nil {
		log.Errorf("an error occurred when get signed url : %v", err)
		return "", err
	}

	return url, nil
}

func (c *GCSPackage) Upload(ctx context.Context, opts *UploadOptions) error {
	ctx, cancel := context.WithTimeout(ctx, opts.Timeout)
	defer cancel()

	bucket := c.Client.Bucket(c.BucketName)
	obj := bucket.Object(opts.Object)
	writer := obj.NewWriter(ctx)
	writer.ChunkSize = 0

	if _, err := io.Copy(writer, opts.File); err != nil {
		log.Errorf("an error occurred when copy file : %v", err)
		return err
	}

	if err := writer.Close(); err != nil {
		log.Errorf("an error occurred when close writer : %v", err)
		return err
	}

	if opts.Public {
		acl := c.Client.Bucket(c.BucketName).Object(opts.Object).ACL()
		err := acl.Set(ctx, storage.AllUsers, storage.RoleReader)
		if err != nil {
			log.Errorf("an error occurred when set acl : %v", err)
			return err
		}
	}

	return nil
}

func (c *GCSPackage) Delete(ctx context.Context, object string, hardDelete bool, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	o := c.Client.Bucket(c.BucketName).Object(object)

	// Optional: set a generation-match precondition to avoid potential race
	// conditions and data corruptions. The request to copy is aborted if the
	// object's generation number does not match your precondition.
	attrs, err := o.Attrs(ctx)
	if err != nil {
		log.Errorf("an error occurred when get attrs : %v", err)
		return errors.Wrap(err, "object.Attrs")
	}
	o = o.If(storage.Conditions{GenerationMatch: attrs.Generation})

	if hardDelete {
		if err := o.Delete(ctx); err != nil {
			log.Errorf("an error occurred when delete object : %v", err)
			return errors.Wrap(err, "object.HardDelete")
		}
		return nil
	}

	// You can't change an object's storage class directly, the only way is
	// to rewrite the object with the desired storage class.
	copier := o.CopierFrom(o)
	copier.StorageClass = coldLineStorageClass
	if _, err := copier.Run(ctx); err != nil {
		log.Errorf("an error occurred when copy object : %v", err)
		return errors.Wrap(err, "object.Copy.Run")
	}

	return nil
}

func (c *GCSPackage) UploadFileInByte(ctx context.Context, fileName string, data []byte) (string, error) {
	var (
		contentType      = "application/octet-stream"
		timeoutInSeconds uint
	)
	defer func(Client *storage.Client) {
		err := Client.Close()
		if err != nil {
			log.Errorf("an error occurred when close client : %v", err)
			return
		}
	}(c.Client)

	if c.TimeoutInSeconds > 0 {
		timeoutInSeconds = c.TimeoutInSeconds
	} else {
		timeoutInSeconds = 30
	}

	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeoutInSeconds)*time.Second) //nolint:gosec
	defer cancel()

	buc := c.Client.Bucket(c.BucketName)
	obj := buc.Object(fileName)
	buf := bytes.NewBuffer(data)

	writer := obj.NewWriter(ctx)
	writer.ChunkSize = 0

	if _, err := io.Copy(writer, buf); err != nil {
		log.Errorf("an error occurred when copy file : %v", err)
		return "", err
	}
	if err := writer.Close(); err != nil {
		log.Errorf("an error occurred when close writer : %v", err)
		return "", err
	}

	if _, err := obj.Update(ctx, storage.ObjectAttrsToUpdate{ContentType: contentType}); err != nil {
		log.Errorf("an error occurred when update object : %v", err)
		return "", err
	}

	url, err := buc.SignedURL(fileName, &storage.SignedURLOptions{
		GoogleAccessID: c.ServiceAccountKeyJSON.ClientEmail,
		PrivateKey:     []byte(c.ServiceAccountKeyJSON.PrivateKey),
		Method:         "GET",
		Expires:        time.Now().Add(time.Duration(c.SignedURLTimeInMinutes) * time.Minute), //nolint:gosec
	})

	if err != nil {
		log.Errorf("an error occurred when get signed url : %v", err)
		return "", err
	}

	return url, nil
}
