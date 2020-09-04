package handler

import (
	"context"
	_ "omo-msa-file/model"

	"github.com/micro/go-micro/v2/logger"
	proto "github.com/xtech-cloud/omo-msp-file/proto/file"
)

type Object struct{}

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

func (this *Object) Get(_ctx context.Context, _req *proto.ObjectGetRequest, _rsp *proto.ObjectGetResponse) error {
	logger.Infof("Received Object.Get, req is %v", _req)
	_rsp.Status = &proto.Status{}

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

	return nil
}
