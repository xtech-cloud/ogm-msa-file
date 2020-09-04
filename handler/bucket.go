package handler

import (
	"context"
	"errors"
	"omo-msa-file/model"

	"github.com/micro/go-micro/v2/logger"
	proto "github.com/xtech-cloud/omo-msp-file/proto/file"
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

	domain := "local"
	if "" != _req.Domain {
		domain = _req.Domain
	}

	bucket := &model.Bucket{
		Name:         _req.Name,
		Domain:       domain,
		AccessKey:    _req.AccessKey,
		AccessSecret: _req.AccessSecret,
	}

	dao := model.NewBucketDAO(nil)
	err := dao.Insert(bucket)
	if errors.Is(err, model.ErrBucketExists) {
		_rsp.Status.Code = 2
		_rsp.Status.Message = "bucket exists"
		return nil
	}
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

	_rsp.Total = total
	_rsp.Entity = make([]*proto.BucketEntity, len(buckets))
	for i, bucket := range buckets {
		_rsp.Entity[i] = &proto.BucketEntity{
			Name:         bucket.Name,
			Domain:       bucket.Domain,
			AccessKey:    bucket.AccessKey,
			AccessSecret: bucket.AccessSecret,
			TotalSize:    bucket.TotalSize,
			FreeSize:     bucket.FreeSize,
		}
	}

	return nil
}

func (this *Bucket) Update(_ctx context.Context, _req *proto.BucketUpdateRequest, _rsp *proto.BlankResponse) error {
	logger.Infof("Received Bucket.Update, req is %v", _req)
	_rsp.Status = &proto.Status{}

	if "" == _req.Name {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "name is required"
		return nil
	}

	bucket := &model.Bucket{
		Name:         _req.Name,
		Domain:       _req.Domain,
		AccessKey:    _req.AccessKey,
		AccessSecret: _req.AccessSecret,
	}

	dao := model.NewBucketDAO(nil)
	err := dao.Update(bucket)
	if errors.Is(err, model.ErrBucketNotFound) {
		_rsp.Status.Code = 2
		_rsp.Status.Message = "bucket not found"
		return nil
	}
	return err
}

func (this *Bucket) Remove(_ctx context.Context, _req *proto.BucketRemoveRequest, _rsp *proto.BlankResponse) error {
	logger.Infof("Received Bucket.Remove, req is %v", _req)
	_rsp.Status = &proto.Status{}

	if "" == _req.Name {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "name is required"
		return nil
	}

	dao := model.NewBucketDAO(nil)
	err := dao.Delete(_req.Name)
	if errors.Is(err, model.ErrBucketNotFound) {
		_rsp.Status.Code = 2
		_rsp.Status.Message = "bucket not found"
		return nil
	}
	return err
}

func (this *Bucket) Get(_ctx context.Context, _req *proto.BucketGetRequest, _rsp *proto.BucketGetResponse) error {
	logger.Infof("Received Bucket.Get, req is %v", _req)
	_rsp.Status = &proto.Status{}

	if "" == _req.Name {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "name is required"
		return nil
	}

	dao := model.NewBucketDAO(nil)
    query := model.BucketQuery {
        Name: _req.Name,
    }
	bucket, err := dao.QueryOne(&query)
	if errors.Is(err, model.ErrBucketNotFound) {
		_rsp.Status.Code = 2
		_rsp.Status.Message = "bucket not found"
		return nil
	}
	_rsp.Entity = &proto.BucketEntity{
		Name:         bucket.Name,
		Domain:       bucket.Domain,
		AccessKey:    bucket.AccessKey,
		AccessSecret: bucket.AccessSecret,
		TotalSize:    bucket.TotalSize,
		FreeSize:     bucket.FreeSize,
	}

	return nil
}
