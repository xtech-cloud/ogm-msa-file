package engine

import (
	"github.com/qiniu/api.v7/v7/auth/qbox"
	"github.com/qiniu/api.v7/v7/storage"
)

func qiniuAuth(_scope string, _accessKey string, _accessSecret string) (string, error) {

	putPolicy := storage.PutPolicy{
		Scope: _scope,
	}
	mac := qbox.NewMac(_accessKey, _accessSecret)
	upToken := putPolicy.UploadToken(mac)
    return upToken, nil
}
