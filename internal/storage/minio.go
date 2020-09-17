package storage

import (
	"errors"
	"io"

	"github.com/minio/minio-go"
)

type MinioConfig struct {
	Endpoint     string `json:"endpoint"`
	AccessKey    string `json:"access_key"`
	AccessSecret string `json:"access_secret"`
	Location     string `json:"location"`
	Secure       bool   `json:"secure"`
}

// Minio implements the Storage interface for
// the MinIO SDK to connect to a MinIO instance,
// Amazon S3 or Google Cloud.
type Minio struct {
	client   *minio.Client
	location string
}

func (m *Minio) Init(params ...interface{}) (err error) {
	if len(params) == 0 {
		return errors.New("minio config must be passed")
	}

	c, ok := params[0].(MinioConfig)
	if !ok {
		return errors.New("invalid config type")
	}

	m.client, err = minio.New(c.Endpoint, c.AccessKey, c.AccessSecret, c.Secure)
	m.location = c.Location

	return
}

func (m *Minio) BucketExists(name string) (bool, error) {
	return m.client.BucketExists(name)
}

func (m *Minio) CreateBucket(name string, location ...string) error {
	return m.client.MakeBucket(name, m.getLocation(location))
}

func (m *Minio) CreateBucketIfNotExists(name string, location ...string) (err error) {
	ok, err := m.BucketExists(name)
	if err == nil && !ok {
		err = m.CreateBucket(name, location...)
	}

	return
}

func (m *Minio) PutObject(bucketName, objectName string, reader io.Reader, objectSize int64, mimeType string) (err error) {
	if err = m.CreateBucketIfNotExists(bucketName, m.location); err != nil {
		return
	}
	_, err = m.client.PutObject(bucketName, objectName, reader, objectSize, minio.PutObjectOptions{
		ContentType: mimeType,
	})
	return
}

func (m *Minio) GetObject(bucketName, objectName string) (io.ReadCloser, int64, error) {
	obj, err := m.client.GetObject(bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, 0, err
	}

	stat, err := obj.Stat()
	if err != nil {
		return nil, 0, err
	}

	return obj, stat.Size, err
}

func (m *Minio) DeleteObject(bucketName, objectName string) error {
	return m.client.RemoveObject(bucketName, objectName)
}

func (m *Minio) getLocation(loc []string) string {
	if len(loc) > 0 {
		return loc[0]
	}
	return m.location
}
