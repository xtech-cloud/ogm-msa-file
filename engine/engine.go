package engine

import (
	proto "github.com/xtech-cloud/omo-msp-file/proto/file"
)

func Prepare(_engine int, _address string, _scope string, _uname string, _accessKey string, _accessSecret string) (string, error) {
	if proto.Engine(_engine) == proto.Engine_ENGINE_QINIU {
		return prepareQiniu(_scope, _uname, _accessKey, _accessSecret)
	}
	return "", nil
}

func Flush(_engine int, _address string, _scope string, _uname string, _accessKey string, _accessSecret string) (int64, error) {
	if proto.Engine(_engine) == proto.Engine_ENGINE_QINIU {
		return flushQiniu(_scope, _uname, _accessKey, _accessSecret)
	}
	return 0, nil
}
