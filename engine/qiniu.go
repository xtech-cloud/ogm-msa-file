package engine

import (
	"github.com/qiniu/api.v7/v7/auth/qbox"
	"github.com/qiniu/api.v7/v7/storage"
	"time"
)

func prepareQiniu(_address string, _scope string, _uname string, _accessKey string, _accessSecret string) (string, error) {

	mac := qbox.NewMac(_accessKey, _accessSecret)
	cfg := storage.Config{
		UseHTTPS: false,
	}
	bucketManager := storage.NewBucketManager(mac, &cfg)
	_, err := bucketManager.Stat(_scope, _uname)
	// 文件存在
	if nil == err {
		return "", nil
	}
	// 内部错误
	// error 是 no such file or directory
	// ErrNoSuchFile 是 No such file or directory
	// 应该是七牛的BUG
	//if !errors.Is(err, storage.ErrNoSuchFile) {
	if err.Error() != "no such file or directory" {
		return "", err
	}
	putPolicy := storage.PutPolicy{
		Scope: _scope,
	}
	// 获取上传的令牌
	upToken := putPolicy.UploadToken(mac)
	return upToken, nil
}

func flushQiniu(_address string, _scope string, _uname string, _accessKey string, _accessSecret string) (int64, error) {

	mac := qbox.NewMac(_accessKey, _accessSecret)
	cfg := storage.Config{
		UseHTTPS: false,
	}
	bucketManager := storage.NewBucketManager(mac, &cfg)
	// 获取已上传的文件的尺寸
	fileInfo, err := bucketManager.Stat(_scope, _uname)
	if nil != err {
		return 0, err
	}
	return fileInfo.Fsize, nil
}

func publishQiniu(_address string, _url string, _scope string, _uname string, _filename string, _accessKey string, _accessSecret string) (string, error) {
	//TODO public是公有的返回公开链接，私有返回一个带有效期的链接
	// 七牛的私有桶没有永久有效期API，使用100年代替
	expiry := 60 * 60 * 24 * 365 * 100
	mac := qbox.NewMac(_accessKey, _accessSecret)
	deadline := time.Now().Add(time.Second * time.Duration(expiry)).Unix()
	url := storage.MakePrivateURL(mac, _address, _uname, deadline)
	return url, nil
}

func previewQiniu(_address string, _url string, _scope string, _uname string, _filename string, _expiry uint64, _accessKey string, _accessSecret string) (string, error) {
	//TODO public是公有的返回公开链接，私有返回一个带有效期的链接
	mac := qbox.NewMac(_accessKey, _accessSecret)
	deadline := time.Now().Add(time.Second * time.Duration(_expiry)).Unix()
	url := storage.MakePrivateURL(mac, _address, _uname, deadline)
	return url, nil
}
