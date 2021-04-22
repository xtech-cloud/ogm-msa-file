package handler

import (
	"context"
	"errors"
	"ogm-msa-file/engine"
	"ogm-msa-file/model"
	"path"
	"strings"

	"github.com/micro/go-micro/v2/logger"
	proto "github.com/xtech-cloud/ogm-msp-file/proto/file"
)

type Object struct{}

func (this *Object) Prepare(_ctx context.Context, _req *proto.ObjectPrepareRequest, _rsp *proto.ObjectPrepareResponse) error {
	logger.Infof("Received Object.Prepare, req is %v", _req)
	_rsp.Status = &proto.Status{}

	if "" == _req.Bucket {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "bucket is required"
		return nil
	}

	if "" == _req.Uname {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "uname is required"
		return nil
	}

	if 0 == _req.Size {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "size is required"
		return nil
	}

	daoBucket := model.NewBucketDAO(nil)
	bucket, err := daoBucket.Get(_req.Bucket)
	if errors.Is(err, model.ErrBucketNotFound) {
		_rsp.Status.Code = 2
		_rsp.Status.Message = "bucket not found"
		return nil
	}

	if bucket.UsedSize+_req.Size > bucket.TotalSize {
		_rsp.Status.Code = 3
		_rsp.Status.Message = "out of capacity"
		return nil
	}

	accessToken, err := engine.Prepare(bucket.Engine, bucket.Address, bucket.Scope, _req.Uname, bucket.AccessKey, bucket.AccessSecret)
	if nil != err {
		return err
	}

	if "" == accessToken {
		_rsp.Status.Code = 200
		_rsp.Status.Message = "object is exists"
		return nil
	}

	_rsp.Address = bucket.Address
	_rsp.Engine = proto.Engine(bucket.Engine)
	_rsp.AccessToken = accessToken
	return nil
}

func (this *Object) Flush(_ctx context.Context, _req *proto.ObjectFlushRequest, _rsp *proto.BlankResponse) error {
	logger.Infof("Received Object.Flush, req is %v", _req)
	_rsp.Status = &proto.Status{}

	if "" == _req.Bucket {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "bucket is required"
		return nil
	}

	if "" == _req.Uname {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "uname is required"
		return nil
	}

	if "" == _req.Path {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "path is required"
		return nil
	}

	daoBucket := model.NewBucketDAO(nil)
	// 获取存储桶
	bucket, err := daoBucket.Get(_req.Bucket)
	if nil != err {
		if errors.Is(err, model.ErrBucketNotFound) {
			_rsp.Status.Code = 2
			_rsp.Status.Message = err.Error()
			return nil
		} else {
			return err
		}

	}

	// 从存储引擎中获取文件的实际大小
	fsize, err := engine.Flush(bucket.Engine, bucket.Address, bucket.Scope, _req.Uname, bucket.AccessKey, bucket.AccessSecret)
	if nil != err {
		return err
	}

	daoObject := model.NewObjectDAO(nil)
	object := &model.Object{
		UUID:     model.ToUUID(_req.Bucket + _req.Path),
		Filepath: _req.Path,
		Bucket:   _req.Bucket,
		MD5:      strings.TrimSuffix(_req.Uname, path.Ext(_req.Uname)),
		Size:     uint64(fsize),
	}

	err = daoObject.Insert(object)
	if errors.Is(err, model.ErrObjectExists) {
		_rsp.Status.Code = 3
		_rsp.Status.Message = err.Error()
		return nil
	}
	if nil != err {
		return err
	}

	count, err := daoObject.CountOfMD5(_req.Bucket, object.MD5)
	if nil != err {
		return err
	}

	if 1 == count {
		// 更新已用空间
		bucket.UsedSize = bucket.UsedSize + uint64(fsize)
		logger.Debugf("the used of size is %d", bucket.UsedSize)
		err = daoBucket.Update(bucket)
		if nil != err {
			return err
		}
	}

	return nil
}

/*
func (this *Object) Upload(_ctx context.Context, _stream proto.Object_UploadStream) error {
    //logger.Infof("Received Object.Upload, req is %v", _req)

    return nil
}

func (this *Object) Download(_ctx context.Context, _req *proto.ObjectDownloadRequest, _stream proto.Object_DownloadStream) error {
    logger.Infof("Received Object.Download, req is %v", _req)

    return nil
}

func (this *Object) Link(_ctx context.Context, _req *proto.ObjectLinkRequest, _rsp *proto.BlankResponse) error {
    logger.Infof("Received Object.Link, req is %v", _req)
    _rsp.Status = &proto.Status{}

    return nil
}
*/

func (this *Object) Get(_ctx context.Context, _req *proto.ObjectGetRequest, _rsp *proto.ObjectGetResponse) error {
	logger.Infof("Received Object.Get, req is %v", _req)
	_rsp.Status = &proto.Status{}

	if "" == _req.Uuid{
		_rsp.Status.Code = 1
		_rsp.Status.Message = "uuid is required"
		return nil
	}

	dao := model.NewObjectDAO(nil)
    object, err := dao.Get(_req.Uuid)
    if nil != err {
        if errors.Is(err, model.ErrObjectNotFound) {
            _rsp.Status.Code = 2
            _rsp.Status.Message = err.Error()
            return nil
        } else {
            return err
        }
    }

    _rsp.Entity = &proto.ObjectEntity{
        Uuid: object.UUID,
        Filepath: object.Filepath,
        Md5: object.MD5,
        Url: object.URL,
        Size: object.Size,
    }
	return nil
}

func (this *Object) Find(_ctx context.Context, _req *proto.ObjectFindRequest, _rsp *proto.ObjectFindResponse) error {
	logger.Infof("Received Object.Find, req is %v", _req)
	_rsp.Status = &proto.Status{}

	dao := model.NewObjectDAO(nil)
    object, err := dao.QueryOne(&model.ObjectQuery{
        Bucket: _req.Bucket,
        Filepath: _req.Filepath,
    })
    if nil != err {
        if errors.Is(err, model.ErrObjectNotFound) {
            _rsp.Status.Code = 2
            _rsp.Status.Message = err.Error()
            return nil
        } else {
            return err
        }
    }

    _rsp.Entity = &proto.ObjectEntity{
        Uuid: object.UUID,
        Filepath: object.Filepath,
        Md5: object.MD5,
        Url: object.URL,
        Size: object.Size,
    }
	return nil
}

func (this *Object) Remove(_ctx context.Context, _req *proto.ObjectRemoveRequest, _rsp *proto.BlankResponse) error {
	logger.Infof("Received Object.Remove, req is %v", _req)
	_rsp.Status = &proto.Status{}

	return nil
}

func (this *Object) List(_ctx context.Context, _req *proto.ObjectListRequest, _rsp *proto.ObjectListResponse) error {
	logger.Infof("Received Object.List, req is %v", _req)
	_rsp.Status = &proto.Status{}

	offset := int64(0)
	count := int64(100)

	if _req.Offset > 0 {
		offset = _req.Offset
	}

	if _req.Count > 0 {
		count = _req.Count
	}

	dao := model.NewObjectDAO(nil)

	total, objects, err := dao.List(offset, count, _req.Bucket)
	if nil != err {
		return nil
	}

	_rsp.Total = uint64(total)
	_rsp.Entity = make([]*proto.ObjectEntity, len(objects))
	for i, object := range objects {
		_rsp.Entity[i] = &proto.ObjectEntity{
			Uuid:     object.UUID,
			Filepath: object.Filepath,
			Md5:      object.MD5,
			Size:     object.Size,
			Url:      object.URL,
		}
	}
	return nil
}

func (this *Object) Search(_ctx context.Context, _req *proto.ObjectSearchRequest, _rsp *proto.ObjectSearchResponse) error {
	logger.Infof("Received Object.Search, req is %v", _req)
	_rsp.Status = &proto.Status{}

	if "" == _req.Bucket {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "bucket is required"
		return nil
	}

	offset := int64(0)
	count := int64(100)

	if _req.Offset > 0 {
		offset = _req.Offset
	}

	if _req.Count > 0 {
		count = _req.Count
	}

	dao := model.NewJoinsDAO(nil)

	query := model.JoinsQuery{
		Bucket:   _req.Bucket,
		Filepath: _req.Prefix,
	}
	total, objects, err := dao.SearchObject(offset, count, &query)
	if nil != err {
		return nil
	}

	_rsp.Total = uint64(total)
	_rsp.Entity = make([]*proto.ObjectEntity, len(objects))
	for i, object := range objects {
		_rsp.Entity[i] = &proto.ObjectEntity{
			Uuid:     object.UUID,
			Filepath: object.Filepath,
			Md5:      object.MD5,
			Size:     object.Size,
			Url:      object.URL,
		}
	}
	return nil
}

func (this *Object) Publish(_ctx context.Context, _req *proto.ObjectPublishRequest, _rsp *proto.ObjectPublishResponse) error {
	logger.Infof("Received Object.Publish, req is %v", _req)
	_rsp.Status = &proto.Status{}

	if "" == _req.Uuid {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "uuid is required"
		return nil
	}

	dao := model.NewObjectDAO(nil)
	object, err := dao.Get(_req.Uuid)
	if nil != err {
		if errors.Is(err, model.ErrObjectNotFound) {
			_rsp.Status.Code = 1
			_rsp.Status.Message = err.Error()
			return nil
		} else {
			return err
		}
	}

	daoBucket := model.NewBucketDAO(nil)
	bucket, err := daoBucket.Get(object.Bucket)
	if nil != err {
		if errors.Is(err, model.ErrBucketNotFound) {
			_rsp.Status.Code = 1
			_rsp.Status.Message = err.Error()
			return nil
		} else {
			return err
		}
	}

    uname := object.MD5 + path.Ext(object.Filepath)
    // 有效期100年
	url, err := engine.Publish(bucket.Engine, bucket.Address, bucket.Scope, uname, 60*60*24*365*100, bucket.AccessKey, bucket.AccessSecret)
	if nil != err {
		return err
	}
    object.URL = url
    err = dao.Update(object)
    if nil != err {
        return nil
    }
    // 赋值对象访问地址
	_rsp.Url = url
	return nil
}

func (this *Object) Preview(_ctx context.Context, _req *proto.ObjectPreviewRequest, _rsp *proto.ObjectPreviewResponse) error {
	logger.Infof("Received Object.Preview, req is %v", _req)
	_rsp.Status = &proto.Status{}

	if "" == _req.Uuid {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "uuid is required"
		return nil
	}

	dao := model.NewObjectDAO(nil)
	object, err := dao.Get(_req.Uuid)
	if nil != err {
		if errors.Is(err, model.ErrObjectNotFound) {
			_rsp.Status.Code = 1
			_rsp.Status.Message = err.Error()
			return nil
		} else {
			return err
		}
	}

    // 如果对象有公开访问地址，返回公开访问地址
    if "" != object.URL {
        _rsp.Url = object.URL
        return nil
    }

    // 如果对象没有公开访问地址，返回一个有效期5分钟的临时访问地址

	daoBucket := model.NewBucketDAO(nil)
	bucket, err := daoBucket.Get(object.Bucket)
	if nil != err {
		if errors.Is(err, model.ErrBucketNotFound) {
			_rsp.Status.Code = 1
			_rsp.Status.Message = err.Error()
			return nil
		} else {
			return err
		}
	}

    uname := object.MD5 + path.Ext(object.Filepath)
    //有效期5分钟
	url, err := engine.Publish(bucket.Engine, bucket.Address, bucket.Scope, uname, 300, bucket.AccessKey, bucket.AccessSecret)
	if nil != err {
		return err
	}
    //!注意： 临时的访问地址不能赋值给Object.URL
	_rsp.Url = url
	return nil
}

func (this *Object) Retract(_ctx context.Context, _req *proto.ObjectRetractRequest, _rsp *proto.BlankResponse) error {
	logger.Infof("Received Object.Retract, req is %v", _req)
	_rsp.Status = &proto.Status{}

	if "" == _req.Uuid {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "uuid is required"
		return nil
	}

	dao := model.NewObjectDAO(nil)
	object, err := dao.Get(_req.Uuid)
	if nil != err {
		if errors.Is(err, model.ErrObjectNotFound) {
			_rsp.Status.Code = 1
			_rsp.Status.Message = err.Error()
			return nil
		} else {
			return err
		}
	}

	daoBucket := model.NewBucketDAO(nil)
	bucket, err := daoBucket.Get(object.Bucket)
	if nil != err {
		if errors.Is(err, model.ErrBucketNotFound) {
			_rsp.Status.Code = 1
			_rsp.Status.Message = err.Error()
			return nil
		} else {
			return err
		}
	}

    uname := object.MD5 + path.Ext(object.Filepath)
    // 有效期60秒
	_, err = engine.Publish(bucket.Engine, bucket.Address, bucket.Scope, uname, 60, bucket.AccessKey, bucket.AccessSecret)
	if nil != err {
		return err
	}
    // 置空对象访问地址
    object.URL = ""
    err = dao.Update(object)
    if nil != err {
        return nil
    }
	return nil
}
