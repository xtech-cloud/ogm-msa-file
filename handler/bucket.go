package handler

import (
	"context"
	"errors"
	"ogm-msa-file/config"
	"ogm-msa-file/model"
	"ogm-msa-file/publisher"

	"github.com/micro/go-micro/v2/logger"
	proto "github.com/xtech-cloud/ogm-msp-file/proto/file"
)

type Bucket struct{}

func (this *Bucket) Make(_ctx context.Context, _req *proto.BucketMakeRequest, _rsp *proto.BlankResponse) error {
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
	if config.Schema.Database.Lite {
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
	}

	dao := model.NewBucketDAO(nil)
	err := dao.Insert(bucket)
	if errors.Is(err, model.ErrBucketExists) {
		_rsp.Status.Code = 2
		_rsp.Status.Message = err.Error()
		return nil
	}

    // 发布消息
    ctx := buildNotifyContext(_ctx, "root")
    publisher.Publish(ctx, "/bucket/make", _req, _rsp)
	return err
}

func (this *Bucket) List(_ctx context.Context, _req *proto.BucketListRequest, _rsp *proto.BucketListResponse) error {
	logger.Infof("Received Bucket.List, req is %v", _req)
	_rsp.Status = &proto.Status{}

	offset := int64(0)
	count := int64(100)

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
		}
	}
	return nil
}

func (this *Bucket) UpdateEngine(_ctx context.Context, _req *proto.BucketUpdateEngineRequest, _rsp *proto.BlankResponse) error {
	logger.Infof("Received Bucket.UpdateEngine, req is %v", _req)
	_rsp.Status = &proto.Status{}

	if "" == _req.Uuid {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "uuid is required"
		return nil
	}

	if proto.Engine_ENGINE_INVALID == _req.Engine {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "engine is required"
		return nil
	}

	bucket := &model.Bucket{
		UUID:         _req.Uuid,
		Engine:       int(_req.Engine),
		Address:      _req.Address,
		Scope:        _req.Scope,
		AccessKey:    _req.AccessKey,
		AccessSecret: _req.AccessSecret,
	}

	dao := model.NewBucketDAO(nil)
	err := dao.Update(bucket)
	if errors.Is(err, model.ErrBucketNotFound) {
		_rsp.Status.Code = 2
		_rsp.Status.Message = err.Error()
		return nil
	}
    // 发布消息
    ctx := buildNotifyContext(_ctx, "root")
    publisher.Publish(ctx, "/bucket/updateengine", _req, _rsp)
	return err
}

func (this *Bucket) UpdateCapacity(_ctx context.Context, _req *proto.BucketUpdateCapacityRequest, _rsp *proto.BlankResponse) error {
	logger.Infof("Received Bucket.UpdateCapacity, req is %v", _req)
	_rsp.Status = &proto.Status{}

	if "" == _req.Uuid {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "uuid is required"
		return nil
	}

	bucket := &model.Bucket{
		UUID:      _req.Uuid,
		TotalSize: _req.Capacity,
	}

	dao := model.NewBucketDAO(nil)
	err := dao.Update(bucket)
	if errors.Is(err, model.ErrBucketNotFound) {
		_rsp.Status.Code = 2
		_rsp.Status.Message = err.Error()
		return nil
	}
    // 发布消息
    ctx := buildNotifyContext(_ctx, "root")
    publisher.Publish(ctx, "/bucket/updatecapacity", _req, _rsp)
	return err
}

func (this *Bucket) ResetToken(_ctx context.Context, _req *proto.BucketResetTokenRequest, _rsp *proto.BlankResponse) error {
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
    // 发布消息
    ctx := buildNotifyContext(_ctx, "root")
    publisher.Publish(ctx, "/bucket/resettoken", _req, _rsp)
	return err
}

func (this *Bucket) Remove(_ctx context.Context, _req *proto.BucketRemoveRequest, _rsp *proto.BlankResponse) error {
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
    // 发布消息
    ctx := buildNotifyContext(_ctx, "root")
    publisher.Publish(ctx, "/bucket/remove", _req, _rsp)
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
	}
	return nil
}

func (this *Bucket) Find(_ctx context.Context, _req *proto.BucketFindRequest, _rsp *proto.BucketFindResponse) error {
	logger.Infof("Received Bucket.Find, req is %v", _req)
	_rsp.Status = &proto.Status{}

	if "" == _req.Name{
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
	}
	return nil
}
