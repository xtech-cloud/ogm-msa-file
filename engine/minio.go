package engine

import (
	"context"
	"errors"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func prepareMinio(_address, _scope string, _uname string, _accessKey string, _accessSecret string) (string, error) {
	useSSL := false
	minioClient, err := minio.New(_address, &minio.Options{
		Creds:  credentials.NewStaticV4(_accessKey, _accessSecret, ""),
		Secure: useSSL,
	})
	if err != nil {
		return "", err
	}

	ctx := context.TODO()
	exists, errBucketExists := minioClient.BucketExists(ctx, _scope)
	if errBucketExists != nil && exists {
		return "", errBucketExists
	}

	if !exists {
		return "", errors.New("bucket not found")
	}

	_, err = minioClient.StatObject(context.Background(), _scope, _uname, minio.StatObjectOptions{})
	if err != nil {
		return "", err
	}
	//TODO return the policy of bucket
	return "POLICY", nil
}

func flushMinio(_address, _scope string, _uname string, _accessKey string, _accessSecret string) (int64, error) {
	useSSL := false
	minioClient, err := minio.New(_address, &minio.Options{
		Creds:  credentials.NewStaticV4(_accessKey, _accessSecret, ""),
		Secure: useSSL,
	})
	if err != nil {
		return 0, err
	}

	// 获取已上传的文件的尺寸
    info, err := minioClient.StatObject(context.Background(), _scope, _uname, minio.StatObjectOptions{})
	if nil != err {
		return 0, err
	}
	return info.Size, nil
}

func publishMinio(_address, _scope string, _uname string, _accessKey string, _accessSecret string) (string, error) {
	useSSL := false
	minioClient, err := minio.New(_address, &minio.Options{
		Creds:  credentials.NewStaticV4(_accessKey, _accessSecret, ""),
		Secure: useSSL,
	})
	if err != nil {
		return "", err
	}

    //TODO public是公有的返回公开链接，私有返回一个带有效期的链接
	presignedURL, err := minioClient.PresignedGetObject(context.Background(), _scope, _uname, time.Second*24*60*60*7, nil)
	if err != nil {
		return "", err
	}
	return presignedURL.String(), nil
}

func previewMinio(_address, _scope string, _uname string, _expiry uint64, _accessKey string, _accessSecret string) (string, error) {
	useSSL := false
	minioClient, err := minio.New(_address, &minio.Options{
		Creds:  credentials.NewStaticV4(_accessKey, _accessSecret, ""),
		Secure: useSSL,
	})
	if err != nil {
		return "", err
	}

    //TODO public是公有的返回公开链接，私有返回一个带有效期的链接
	presignedURL, err := minioClient.PresignedGetObject(context.Background(), _scope, _uname, time.Second * time.Duration(_expiry), nil)
	if err != nil {
		return "", err
	}
	return presignedURL.String(), nil
}
