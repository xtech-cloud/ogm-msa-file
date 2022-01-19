package engine

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/minio/minio-go/v7/pkg/s3utils"
	"github.com/minio/minio-go/v7/pkg/signer"
)

func prepareMinio(_address string, _url string, _scope string, _uname string, _accessKey string, _accessSecret string, _expiry int64, _override bool) (string, error) {
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
	if errBucketExists != nil {
		return "", errBucketExists
	}

	if !exists {
		return "", errors.New("bucket not found")
	}

	_, err = minioClient.StatObject(context.Background(), _scope, _uname, minio.StatObjectOptions{})
	if nil != err {
		if err.Error() != "The specified key does not exist." {
			return "", err
		}
	} else {
		if !_override {
			// 文件存在
			return "", nil
		}
	}

	// 方式一：使MINIO SDK
	//token := fmt.Sprintf("%s %s", _accessKey, _accessSecret)

	// 方式二： 使用HTTP PUT
	// 签名不需要使用minio.Client，使用minio.Client将在签名前访问minio服务，
	// 如果minio地址在外网不可访问时，生成的链接外网可不可用，
	// 直接使用签名函数签名，但不保证文件有效
	queryUrl := make(url.Values)
	queryUrl.Set("response-content-disposition", fmt.Sprintf("attachment; filename=\"%s\"", _uname))
	targetURL := fmt.Sprintf("%s/%s/%s?%s", _url, _scope, s3utils.EncodePath(_uname), s3utils.QueryEncode(queryUrl))
	req, err := http.NewRequestWithContext(context.TODO(), "PUT", targetURL, nil)
	if err != nil {
		return "", err
	}
	// 1day
	expiry := int64(60 * 60 * 24)
	if _expiry > 0 {
		expiry = _expiry
	}
	// TODO us-east-1 为默认的location，需要实现可配置
	req = signer.PreSignV4(*req, _accessKey, _accessSecret, "", "us-east-1", expiry)
	return req.URL.String(), nil
}

func flushMinio(_address string, _scope string, _uname string, _accessKey string, _accessSecret string) (int64, error) {
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

func publishMinio(_address string, _url string, _scope string, _uname string, _filename string, _accessKey string, _accessSecret string) (string, error) {
	/*
		useSSL := false
		minioClient, err := minio.New(_address, &minio.Options{
			Creds:  credentials.NewStaticV4(_accessKey, _accessSecret, ""),
			Secure: useSSL,
		})
		if err != nil {
			return "", err
		}

		//TODO public是公有的返回公开链接，私有返回一个带有效期的链接
		//presignedURL, err := minioClient.PresignedGetObject(context.Background(), _scope, _uname, time.Second*24*60*60*7, nil)
		//if err != nil {
		//	return "", err
		//}
		//return presignedURL.String(), nil
	*/
	queryUrl := make(url.Values)
	queryUrl.Set("response-content-disposition", fmt.Sprintf("attachment; filename=\"%s\"", _filename))
	targetURL := fmt.Sprintf("%s/%s/%s?%s", _url, _scope, s3utils.EncodePath(_uname), s3utils.QueryEncode(queryUrl))
	return targetURL, nil
}

func previewMinio(_address string, _url string, _scope string, _uname string, _filename string, _expiry int64, _accessKey string, _accessSecret string) (string, error) {
	/*
		useSSL := false
		minioClient, err := minio.New(_address, &minio.Options{
			Creds:  credentials.NewStaticV4(_accessKey, _accessSecret, ""),
			Secure: useSSL,
		})
		if err != nil {
			return "", err
		}

		//reqParams := make(url.Values)
		//reqParams.Set("response-content-disposition", fmt.Sprintf("attachment; filename=\"%s\"", _filename))

		//TODO public是公有的返回公开链接，私有返回一个带有效期的链接
		presignedURL, err := minioClient.PresignedGetObject(context.Background(), _scope, _uname, time.Second*time.Duration(_expiry), nil)
		if err != nil {
			return "", err
		}
		return presignedURL.String(), nil
	*/

	// 签名不需要使用minio.Client，使用minio.Client将在签名前访问minio服务，
	// 如果minio地址在外网不可访问时，生成的链接外网可不可用，
	// 直接使用签名函数签名，但不保证文件有效
	queryUrl := make(url.Values)
	queryUrl.Set("response-content-disposition", fmt.Sprintf("attachment; filename=\"%s\"", _filename))
	targetURL := fmt.Sprintf("%s/%s/%s?%s", _url, _scope, s3utils.EncodePath(_uname), s3utils.QueryEncode(queryUrl))
	req, err := http.NewRequestWithContext(context.TODO(), "GET", targetURL, nil)
	if err != nil {
		return "", err
	}
	// 1 hour
	expiry := int64(60 * 60)
	if _expiry > 0 {
		expiry = _expiry
	}
	// TODO us-east-1 为默认的location，需要实现可配置
	req = signer.PreSignV4(*req, _accessKey, _accessSecret, "", "us-east-1", expiry)
	return req.URL.String(), nil
}

func saveMinio(_address string, _scope string, _uname string, _reader io.Reader, _size int64, _accessKey string, _accessSecret string) error {
	useSSL := false
	minioClient, err := minio.New(_address, &minio.Options{
		Creds:  credentials.NewStaticV4(_accessKey, _accessSecret, ""),
		Secure: useSSL,
	})
	if err != nil {
		return err
	}

	ctx := context.TODO()
	_, err = minioClient.PutObject(ctx, _scope, _uname, _reader, _size, minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		return err
	}
	return nil
}
