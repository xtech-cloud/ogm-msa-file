package engine

import (
	"github.com/qiniu/api.v7/v7/auth/qbox"
	"github.com/qiniu/api.v7/v7/storage"
)

func prepareQiniu(_scope string, _uname string, _accessKey string, _accessSecret string) (string, error) {

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
	upToken := putPolicy.UploadToken(mac)
	return upToken, nil
}

func flushQiniu(_scope string, _uname string, _accessKey string, _accessSecret string) (int64, error) {

	mac := qbox.NewMac(_accessKey, _accessSecret)
	cfg := storage.Config{
		UseHTTPS: false,
	}
	bucketManager := storage.NewBucketManager(mac, &cfg)
	fileInfo, err := bucketManager.Stat(_scope, _uname)
	if nil != err {
		return 0, err
	}
	return fileInfo.Fsize, nil
}
