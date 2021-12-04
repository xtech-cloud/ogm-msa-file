package handler

import (
	"context"
	"errors"
	"ogm-file/config"
	"ogm-file/model"

	"github.com/asim/go-micro/v3/logger"
	proto "github.com/xtech-cloud/ogm-msp-file/proto/file"
)

type Bucket struct{}

func (this *Bucket) Make(_ctx context.Context, _req *proto.BucketMakeRequest, _rsp *proto.UuidResponse) error {
	logger.Infof("Received Bucket.Make, req is %v", _req)
	_rsp.Status = &proto.Status{}

	if "" == _req.Name {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "name is required"
		return nil
	}

	if 0 == _req.Capacity {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "capacity is required"
		return nil
	}

	// 默认存储引擎为本地
	engine := int(proto.Engine_ENGINE_LOCAL)
	if proto.Engine_ENGINE_INVALID != _req.Engine {
		engine = int(_req.Engine)
	}

	// 本地数据库使用存储桶名生成UUID，方便测试和开发
	uuid := model.NewUUID()
	if "sqlite" == config.Schema.Database.Driver {
		uuid = model.ToUUID(_req.Name)
	}

	bucket := &model.Bucket{
		UUID:         uuid,
		Name:         _req.Name,
		Engine:       engine,
		Token:        model.NewUUID(),
		TotalSize:    _req.Capacity,
		Address:      _req.Address,
		Scope:        _req.Scope,
		AccessKey:    _req.AccessKey,
		AccessSecret: _req.AccessSecret,
		Url:          _req.Url,
	}

	dao := model.NewBucketDAO(nil)
	err := dao.Insert(bucket)
	if errors.Is(err, model.ErrBucketExists) {
		_rsp.Status.Code = 2
		_rsp.Status.Message = err.Error()
		return nil
	}

	_rsp.Uuid = uuid
	return err
}

func (this *Bucket) List(_ctx context.Context, _req *proto.BucketListRequest, _rsp *proto.BucketListResponse) error {
	logger.Infof("Received Bucket.List, req is %v", _req)
	_rsp.Status = &proto.Status{}

	offset := int64(0)
	count := int64(0)

	if _req.Offset > 0 {
		offset = _req.Offset
	}

	if _req.Count > 0 {
		count = _req.Count
	}

	dao := model.NewBucketDAO(nil)

	total, err := dao.Count()
	if nil != err {
		return nil
	}
	buckets, err := dao.List(offset, count)
	if nil != err {
		return nil
	}

	_rsp.Total = uint64(total)
	_rsp.Entity = make([]*proto.BucketEntity, len(buckets))
	for i, bucket := range buckets {
		_rsp.Entity[i] = &proto.BucketEntity{
			Uuid:         bucket.UUID,
			Name:         bucket.Name,
			Engine:       proto.Engine(bucket.Engine),
			TotalSize:    bucket.TotalSize,
			UsedSize:     bucket.UsedSize,
			Token:        bucket.Token,
			Address:      bucket.Address,
			Scope:        bucket.Scope,
			AccessKey:    bucket.AccessKey,
			AccessSecret: bucket.AccessSecret,
			Url:          bucket.Url,
		}
	}
	return nil
}

func (this *Bucket) Search(_ctx context.Context, _req *proto.BucketSearchRequest, _rsp *proto.BucketSearchResponse) error {
	logger.Infof("Received Bucket.Search, req is %v", _req)
	_rsp.Status = &proto.Status{}

	offset := int64(0)
	count := int64(0)

	if _req.Offset > 0 {
		offset = _req.Offset
	}

	if _req.Count > 0 {
		count = _req.Count
	}

	dao := model.NewBucketDAO(nil)

	total, err := dao.CountByName(_req.Name)
	if nil != err {
		return nil
	}
	buckets, err := dao.Search(offset, count, _req.Name)
	if nil != err {
		return nil
	}

	_rsp.Total = uint64(total)
	_rsp.Entity = make([]*proto.BucketEntity, len(buckets))
	for i, bucket := range buckets {
		_rsp.Entity[i] = &proto.BucketEntity{
			Uuid:         bucket.UUID,
			Name:         bucket.Name,
			Engine:       proto.Engine(bucket.Engine),
			TotalSize:    bucket.TotalSize,
			UsedSize:     bucket.UsedSize,
			Token:        bucket.Token,
			Address:      bucket.Address,
			Scope:        bucket.Scope,
			AccessKey:    bucket.AccessKey,
			AccessSecret: bucket.AccessSecret,
			Url:          bucket.Url,
		}
	}
	return nil
}

func (this *Bucket) Update(_ctx context.Context, _req *proto.BucketUpdateRequest, _rsp *proto.UuidResponse) error {
	logger.Infof("Received Bucket.Update, req is %v", _req)
	_rsp.Status = &proto.Status{}

	if "" == _req.Uuid {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "uuid is required"
		return nil
	}

	bucket := &model.Bucket{
		UUID:         _req.Uuid,
		Name:         _req.Name,
		TotalSize:    _req.Capacity,
		Engine:       int(_req.Engine),
		Address:      _req.Address,
		Scope:        _req.Scope,
		AccessKey:    _req.AccessKey,
		AccessSecret: _req.AccessSecret,
		Url:          _req.Url,
	}

	dao := model.NewBucketDAO(nil)
	err := dao.Update(bucket)
	if errors.Is(err, model.ErrBucketNotFound) {
		_rsp.Status.Code = 2
		_rsp.Status.Message = err.Error()
		return nil
	}

	_rsp.Uuid = _req.Uuid
	return err
}

func (this *Bucket) ResetToken(_ctx context.Context, _req *proto.BucketResetTokenRequest, _rsp *proto.UuidResponse) error {
	logger.Infof("Received Bucket.ResetToken, req is %v", _req)
	_rsp.Status = &proto.Status{}

	if "" == _req.Uuid {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "uuid is required"
		return nil
	}

	bucket := &model.Bucket{
		UUID:  _req.Uuid,
		Token: model.NewUUID(),
	}

	dao := model.NewBucketDAO(nil)
	err := dao.Update(bucket)
	if errors.Is(err, model.ErrBucketNotFound) {
		_rsp.Status.Code = 2
		_rsp.Status.Message = err.Error()
		return nil
	}
	_rsp.Uuid = _req.Uuid
	return err
}

func (this *Bucket) Remove(_ctx context.Context, _req *proto.BucketRemoveRequest, _rsp *proto.UuidResponse) error {
	logger.Infof("Received Bucket.Remove, req is %v", _req)
	_rsp.Status = &proto.Status{}

	if "" == _req.Uuid {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "uuid is required"
		return nil
	}

	dao := model.NewBucketDAO(nil)
	err := dao.Delete(_req.Uuid)
	if errors.Is(err, model.ErrBucketNotFound) {
		_rsp.Status.Code = 2
		_rsp.Status.Message = err.Error()
		return nil
	}
	_rsp.Uuid = _req.Uuid
	return err
}

func (this *Bucket) Get(_ctx context.Context, _req *proto.BucketGetRequest, _rsp *proto.BucketGetResponse) error {
	logger.Infof("Received Bucket.Get, req is %v", _req)
	_rsp.Status = &proto.Status{}

	if "" == _req.Uuid {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "uuid is required"
		return nil
	}

	dao := model.NewBucketDAO(nil)
	bucket, err := dao.Get(_req.Uuid)
	if errors.Is(err, model.ErrBucketNotFound) {
		_rsp.Status.Code = 2
		_rsp.Status.Message = err.Error()
		return nil
	}
	_rsp.Entity = &proto.BucketEntity{
		Uuid:         bucket.UUID,
		Name:         bucket.Name,
		Engine:       proto.Engine(bucket.Engine),
		TotalSize:    bucket.TotalSize,
		UsedSize:     bucket.UsedSize,
		Token:        bucket.Token,
		Address:      bucket.Address,
		Scope:        bucket.Scope,
		AccessKey:    bucket.AccessKey,
		AccessSecret: bucket.AccessSecret,
		Url:          bucket.Url,
	}
	return nil
}

func (this *Bucket) Find(_ctx context.Context, _req *proto.BucketFindRequest, _rsp *proto.BucketFindResponse) error {
	logger.Infof("Received Bucket.Find, req is %v", _req)
	_rsp.Status = &proto.Status{}

	if "" == _req.Name {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "name is required"
		return nil
	}

	dao := model.NewBucketDAO(nil)
	query := model.BucketQuery{
		Name: _req.Name,
	}
	bucket, err := dao.QueryOne(&query)
	if errors.Is(err, model.ErrBucketNotFound) {
		_rsp.Status.Code = 2
		_rsp.Status.Message = err.Error()
		return nil
	}
	_rsp.Entity = &proto.BucketEntity{
		Uuid:         bucket.UUID,
		Name:         bucket.Name,
		Engine:       proto.Engine(bucket.Engine),
		TotalSize:    bucket.TotalSize,
		UsedSize:     bucket.UsedSize,
		Token:        bucket.Token,
		Address:      bucket.Address,
		Scope:        bucket.Scope,
		AccessKey:    bucket.AccessKey,
		AccessSecret: bucket.AccessSecret,
		Url:          bucket.Url,
	}
	return nil
}
