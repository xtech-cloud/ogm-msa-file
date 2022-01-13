package engine

import (
	"errors"
	proto "github.com/xtech-cloud/ogm-msp-file/proto/file"
)

func Prepare(_engine int, _address string, _url string, _scope string, _uname string, _accessKey string, _accessSecret string) (string, error) {
	switch proto.Engine(_engine) {
	case proto.Engine_ENGINE_QINIU:
		return prepareQiniu(_address, _scope, _uname, _accessKey, _accessSecret)
	case proto.Engine_ENGINE_MINIO:
		return prepareMinio(_address, _url, _scope, _uname, _accessKey, _accessSecret)
	}
	return "", errors.New("unsupported engine")
}

func Flush(_engine int, _address string, _scope string, _uname string, _accessKey string, _accessSecret string) (int64, error) {
	switch proto.Engine(_engine) {
	case proto.Engine_ENGINE_QINIU:
		return flushQiniu(_address, _scope, _uname, _accessKey, _accessSecret)
	case proto.Engine_ENGINE_MINIO:
		return flushMinio(_address, _scope, _uname, _accessKey, _accessSecret)
	}
	return 0, errors.New("unsupported engine")
}

// 生成一个永久有效链接
func Publish(_engine int, _address string, _url string, _scope string, _uname string, _filename string, _accessKey string, _accessSecret string) (string, error) {
	switch proto.Engine(_engine) {
	case proto.Engine_ENGINE_QINIU:
		return publishQiniu(_address, _url, _scope, _uname, _filename, _accessKey, _accessSecret)
	case proto.Engine_ENGINE_MINIO:
		return publishMinio(_address, _url, _scope, _uname, _filename, _accessKey, _accessSecret)
	}
	return "", errors.New("unsupported engine")
}

// 生成一个有效期链接
func Preview(_engine int, _address string, _url string, _scope string, _uname string, _filename string, _expiry uint64, _accessKey string, _accessSecret string) (string, error) {
	switch proto.Engine(_engine) {
	case proto.Engine_ENGINE_QINIU:
		return previewQiniu(_address, _url, _scope, _uname, _filename, _expiry, _accessKey, _accessSecret)
	case proto.Engine_ENGINE_MINIO:
		return previewMinio(_address, _url, _scope, _uname, _filename, _expiry, _accessKey, _accessSecret)
	}
	return "", errors.New("unsupported engine")
}
