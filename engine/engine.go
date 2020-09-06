package engine

import (
	proto "github.com/xtech-cloud/omo-msp-file/proto/file"
)

func Auth(_engine int, _address string, _scope string, _accessKey string, _accessSecret string) (string, error) {
	if proto.Engine(_engine) == proto.Engine_ENGINE_QINIU {
		return qiniuAuth(_scope, _accessKey, _accessSecret)
	}
	return "", nil
}
