package handler

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"ogm-file/config"
	"ogm-file/engine"
	"ogm-file/model"
	"strings"

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
	if nil != err {
		_rsp.Status.Code = -1
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

func (this *Bucket) GenerateManifest(_ctx context.Context, _req *proto.BucketGenerateManifestRequest, _rsp *proto.BucketGenerateManifestResponse) error {
	logger.Infof("Received Bucket.GenerateManifest, req is %v", _req)
	_rsp.Status = &proto.Status{}

	if "" == _req.Uuid {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "uuid is required"
		return nil
	}

	if nil == _req.Field || len(_req.Field) == 0 {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "field is required"
		return nil
	}

	if "json" != _req.Format {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "format only support json"
		return nil
	}

	// 构建where查询语句列表
	like_sql := make([]string, len(_req.Include))
	for i, e := range _req.Include {
		like_sql[i] = strings.ReplaceAll(e, "*", "%")
	}
	notlike_sql := make([]string, len(_req.Exclude))
	for i, e := range _req.Exclude {
		notlike_sql[i] = strings.ReplaceAll(e, "*", "%")
	}

	dao := model.NewObjectDAO(nil)

	objects, err := dao.WhereFilepath(_req.Uuid, like_sql, notlike_sql)
	if nil != err {
		_rsp.Status.Code = -1
		_rsp.Status.Message = err.Error()
		return nil
	}

	hasFilepath := false
	hasUname := false
	hasUrl := false
	hasMd5 := false
	hasSize := false
	for _, e := range _req.Field {
		if e == "filepath" {
			hasFilepath = true
		} else if e == "uname" {
			hasUname = true
		} else if e == "url" {
			hasUrl = true
		} else if e == "md5" {
			hasMd5 = true
		} else if e == "size" {
			hasSize = true
		}
	}

	// 过滤字段
	content := make([]map[string]interface{}, len(objects))
	for i, e := range objects {
		content[i] = make(map[string]interface{})
		if hasFilepath {
			content[i]["filepath"] = e.Filepath
		}
		if hasUname {
			content[i]["uname"] = e.UName
		}
		if hasUrl {
			content[i]["url"] = e.URL
		}
		if hasMd5 {
			content[i]["md5"] = e.MD5
		}
		if hasSize {
			content[i]["size"] = e.Size
		}
	}

	// 生成json
	data, err := json.Marshal(content)
	if nil != err {
		_rsp.Status.Code = -1
		_rsp.Status.Message = err.Error()
		return nil
	}

	result := ""
	if "" != _req.Template {
		result = strings.ReplaceAll(_req.Template, "$content$", string(data))
	} else {
		result = string(data)
	}

	if "" != _req.SaveAs {
		daoBucket := model.NewBucketDAO(nil)
		bucket, err := daoBucket.Get(_req.Uuid)
		if nil != err {
			_rsp.Status.Code = -1
			_rsp.Status.Message = err.Error()
			return nil
		}
		size := int64(len([]byte(result)))
		length := int64(len(result))
		md5sum := md5.New()
		md5sum.Write([]byte(result))
		md5str := hex.EncodeToString(md5sum.Sum(nil))
		//保存进存储引擎
		reader := strings.NewReader(result)
		err = engine.Save(bucket.Engine, bucket.Address, bucket.Scope, _req.SaveAs, reader, length, bucket.AccessKey, bucket.AccessSecret)
		if nil != err {
			_rsp.Status.Code = 9
			_rsp.Status.Message = err.Error()
			return nil
		}

		// 写入数据库
		object := &model.Object{
			UUID:     model.ToUUID(_req.Uuid + _req.SaveAs),
			Filepath: _req.SaveAs,
			Bucket:   _req.Uuid,
			MD5:      strings.ToUpper(md5str),
			UName:    _req.SaveAs,
			Size:     uint64(size),
		}

		err = dao.Upsert(object)
		if nil != err {
			_rsp.Status.Code = 9
			_rsp.Status.Message = err.Error()
			return nil
		}
		_rsp.Result = ""
	} else {
		_rsp.Result = result
	}

	return nil
}
