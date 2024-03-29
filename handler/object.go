package handler

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"ogm-file/engine"
	"ogm-file/model"
	"path/filepath"

	"github.com/asim/go-micro/v3/logger"
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

	override := false
	uname := ""
	if bucket.Mode == "hash" {
		if "" == _req.Hash {
			_rsp.Status.Code = 1
			_rsp.Status.Message = "hash is required"
			return nil
		}
		override = false
		uname = _req.Hash
	} else if bucket.Mode == "path" {
		if "" == _req.Path {
			_rsp.Status.Code = 1
			_rsp.Status.Message = "path is required"
			return nil
		}
		override = _req.Override
		uname = _req.Path
	} else {
		_rsp.Status.Code = -1
		_rsp.Status.Message = "the mode of bucket is invalid"
		return nil
	}

	if bucket.UsedSize+_req.Size > bucket.TotalSize {
		_rsp.Status.Code = 3
		_rsp.Status.Message = "out of capacity"
		return nil
	}

	accessToken, err := engine.Prepare(bucket.Engine, bucket.Address, bucket.Url, bucket.Scope, uname, bucket.AccessKey, bucket.AccessSecret, _req.Expiry, override)
	if nil != err {
		_rsp.Status.Code = 9
		_rsp.Status.Message = err.Error()
		return nil
	}

	if "" == accessToken {
		_rsp.Status.Code = 200
		_rsp.Status.Message = "object is exists"
		return nil
	}

	_rsp.Url = bucket.Url
	_rsp.Engine = proto.Engine(bucket.Engine)
	_rsp.AccessToken = accessToken
	return nil
}

func (this *Object) Flush(_ctx context.Context, _req *proto.ObjectFlushRequest, _rsp *proto.UuidResponse) error {
	logger.Infof("Received Object.Flush, req is %v", _req)
	_rsp.Status = &proto.Status{}

	if "" == _req.Bucket {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "bucket is required"
		return nil
	}

	if "" == _req.Hash {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "hash is required"
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

	uname := ""
	if bucket.Mode == "hash" {
		uname = _req.Hash
	} else if bucket.Mode == "path" {
		uname = _req.Path
	} else {
		_rsp.Status.Code = -1
		_rsp.Status.Message = "the mode of bucket is invalid"
		return nil
	}

	fsize, err := engine.Flush(bucket.Engine, bucket.Address, bucket.Scope, uname, bucket.AccessKey, bucket.AccessSecret)
	if nil != err {
		return err
	}

	daoObject := model.NewObjectDAO(nil)
	object := &model.Object{
		UUID:   model.ToUUID(_req.Bucket + _req.Path),
		Path:   _req.Path,
		Bucket: _req.Bucket,
		Hash:   _req.Hash,
		Size:   uint64(fsize),
	}

	err = daoObject.Upsert(object)
	if errors.Is(err, model.ErrObjectExists) {
		_rsp.Status.Code = 3
		_rsp.Status.Message = err.Error()
		return nil
	}
	if nil != err {
		_rsp.Status.Code = -1
		_rsp.Status.Message = err.Error()
		return nil
	}

	sum, err := daoObject.SumOfBucket(_req.Bucket)
	if nil != err {
		_rsp.Status.Code = -1
		_rsp.Status.Message = err.Error()
		return nil
	}

	// 更新已用空间
	bucket.UsedSize = sum
	err = daoBucket.Update(bucket)
	if nil != err {
		_rsp.Status.Code = -1
		_rsp.Status.Message = err.Error()
		return nil
	}

	_rsp.Uuid = object.UUID
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

	if "" == _req.Uuid {
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
		Path: object.Path,
		Hash: object.Hash,
		Url:  object.URL,
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
		Path:   _req.Path,
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
		Path: object.Path,
		Hash: object.Hash,
		Url:  object.URL,
		Size: object.Size,
	}
	return nil
}

func (this *Object) Remove(_ctx context.Context, _req *proto.ObjectRemoveRequest, _rsp *proto.UuidResponse) error {
	logger.Infof("Received Object.Remove, req is %v", _req)
	_rsp.Status = &proto.Status{}

	if "" == _req.Uuid {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "uuid is required"
		return nil
	}

	dao := model.NewObjectDAO(nil)
	err := dao.Delete(_req.Uuid)
	if nil != err {
		_rsp.Status.Code = -1
		_rsp.Status.Message = err.Error()
		return nil
	}
	_rsp.Uuid = _req.Uuid
	return nil
}

func (this *Object) List(_ctx context.Context, _req *proto.ObjectListRequest, _rsp *proto.ObjectListResponse) error {
	logger.Infof("Received Object.List, req is %v", _req)
	_rsp.Status = &proto.Status{}

	if _req.Bucket == "" {
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

	dao := model.NewObjectDAO(nil)

	total, objects, err := dao.List(offset, count, _req.Bucket)
	if nil != err {
		_rsp.Status.Code = -1
		_rsp.Status.Message = err.Error()
		return nil
	}

	_rsp.Total = uint64(total)
	_rsp.Entity = make([]*proto.ObjectEntity, len(objects))
	for i, object := range objects {
		_rsp.Entity[i] = &proto.ObjectEntity{
			Uuid: object.UUID,
			Path: object.Path,
			Hash: object.Hash,
			Size: object.Size,
			Url:  object.URL,
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

	dao := model.NewObjectDAO(nil)

	total, objects, err := dao.Search(offset, count, _req.Bucket, _req.Prefix, _req.Name)
	if nil != err {
		return nil
	}

	_rsp.Total = uint64(total)
	_rsp.Entity = make([]*proto.ObjectEntity, len(objects))
	for i, object := range objects {
		_rsp.Entity[i] = &proto.ObjectEntity{
			Uuid: object.UUID,
			Path: object.Path,
			Hash: object.Hash,
			Size: object.Size,
			Url:  object.URL,
		}
	}
	return nil
}

/// 发布
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
	uname := ""
	if bucket.Mode == "hash" {
		uname = object.Hash
	} else if bucket.Mode == "path" {
		uname = object.Path
	} else {
		_rsp.Status.Code = -1
		_rsp.Status.Message = "the mode of bucket is invalid"
		return nil
	}

	filename := filepath.Base(object.Path)
	url, err := engine.Publish(bucket.Engine, bucket.Address, bucket.Url, bucket.Scope, uname, filename, bucket.AccessKey, bucket.AccessSecret)
	if nil != err {
		return err
	}
	// 将永久链接赋值给文件对象
	object.URL = url
	err = dao.Update(object)
	if nil != err {
		return nil
	}
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
			_rsp.Status.Code = -1
			_rsp.Status.Message = err.Error()
			return nil
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
			_rsp.Status.Code = -1
			_rsp.Status.Message = err.Error()
			return nil
		}
	}

	uname := ""
	if bucket.Mode == "hash" {
		uname = object.Hash
	} else if bucket.Mode == "path" {
		uname = object.Path
	} else {
		_rsp.Status.Code = -1
		_rsp.Status.Message = "the mode of bucket is invalid"
		return nil
	}
	filename := filepath.Base(object.Path)
	url, err := engine.Preview(bucket.Engine, bucket.Address, bucket.Url, bucket.Scope, uname, filename, _req.Expiry, bucket.AccessKey, bucket.AccessSecret)
	if nil != err {
		_rsp.Status.Code = -1
		_rsp.Status.Message = err.Error()
		return nil
	}
	//!注意： 临时的访问地址不能赋值给Object.URL
	_rsp.Url = url
	return nil
}

func (this *Object) Retract(_ctx context.Context, _req *proto.ObjectRetractRequest, _rsp *proto.UuidResponse) error {
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

	uname := ""
	if bucket.Mode == "hash" {
		uname = object.Hash
	} else if bucket.Mode == "path" {
		uname = object.Path
	} else {
		_rsp.Status.Code = -1
		_rsp.Status.Message = "the mode of bucket is invalid"
		return nil
	}

	// 有效期60秒
	filename := filepath.Base(object.Path)
	_, err = engine.Preview(bucket.Engine, bucket.Address, bucket.Url, bucket.Scope, uname, filename, 60, bucket.AccessKey, bucket.AccessSecret)
	if nil != err {
		return err
	}
	// 置空对象访问地址
	object.URL = ""
	err = dao.Update(object)
	if nil != err {
		return nil
	}

	_rsp.Uuid = _req.Uuid
	return nil
}

func (this *Object) ConvertFromBase64(_ctx context.Context, _req *proto.ObjectConvertFromBase64Request, _rsp *proto.ObjectConvertFromBase64Response) error {
	logger.Infof("Received Object.ConvertFromBase64Request, req is bucket:%v source:%v", _req.Bucket, len(_req.Source))
	_rsp.Status = &proto.Status{}

	daoBucket := model.NewBucketDAO(nil)
	bucket, err := daoBucket.Get(_req.Bucket)
	if nil != err {
		_rsp.Status.Code = -1
		_rsp.Status.Message = err.Error()
		return nil
	}

	dao := model.NewObjectDAO(nil)

	failure := make([]string, 0)
	for _, source := range _req.Source {
		if "" == source.Path {
			failure = append(failure, source.Path)
			continue
		}
		data, err := base64.StdEncoding.DecodeString(source.Content)
		if nil != err {
			failure = append(failure, source.Path)
			continue
		}

		size := int64(len(data))
		//保存进存储引擎
		reader := bytes.NewReader(data)
		err = engine.Save(bucket.Engine, bucket.Address, bucket.Scope, source.Path, reader, size, bucket.AccessKey, bucket.AccessSecret)
		if nil != err {
			failure = append(failure, source.Path)
			continue
		}

		// 写入数据库
		object := &model.Object{
			UUID:   model.ToUUID(_req.Bucket + source.Path),
			Path:   source.Path,
			Bucket: _req.Bucket,
			Hash:   model.Md5FromBytes(data),
			Size:   uint64(size),
		}
		err = dao.Upsert(object)
		if nil != err {
			failure = append(failure, source.Path)
			continue
		}
	}

	_rsp.Failure = failure

	return nil
}

func (this *Object) ConvertFromUrl(_ctx context.Context, _req *proto.ObjectConvertFromUrlRequest, _rsp *proto.ObjectConvertFromUrlResponse) error {
	logger.Infof("Received Object.ConvertFromUrlRequest, req is bucket:%v source:%v", _req.Bucket, len(_req.Source))
	_rsp.Status = &proto.Status{}

	daoBucket := model.NewBucketDAO(nil)
	bucket, err := daoBucket.Get(_req.Bucket)
	if nil != err {
		_rsp.Status.Code = -1
		_rsp.Status.Message = err.Error()
		return nil
	}

	dao := model.NewObjectDAO(nil)

	failure := make([]string, 0)
	for _, source := range _req.Source {
		if "" == source.Path {
			failure = append(failure, source.Path)
			continue
		}
		// 写入数据库
		object := &model.Object{
			UUID:   model.ToUUID(bucket.UUID + source.Path),
			Path:   source.Path,
			Bucket: _req.Bucket,
			Hash:   source.Hash,
			URL:    source.Content,
			Size:   uint64(source.Size),
		}
		err = dao.Upsert(object)
		if nil != err {
			failure = append(failure, source.Path)
			continue
		}
	}

	_rsp.Failure = failure
	return nil
}
