package main

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"ogm-file/config"
	"ogm-file/handler"
	"ogm-file/model"
	"os"
	"path/filepath"
	"time"

	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/logger"
	"github.com/asim/go-micro/plugins/server/grpc/v3"
	proto "github.com/xtech-cloud/ogm-msp-file/proto/file"
)

func main() {
	config.Setup()
	model.Setup()
    defer model.Cancel()
	model.AutoMigrateDatabase()

	// New Service
	service := micro.NewService(
        micro.Server(grpc.NewServer()),
		micro.Name(config.Schema.Service.Name),
		micro.Version(BuildVersion),
		micro.RegisterTTL(time.Second*time.Duration(config.Schema.Service.TTL)),
		micro.RegisterInterval(time.Second*time.Duration(config.Schema.Service.Interval)),
		micro.Address(config.Schema.Service.Address),
	)

	// Initialise service
	service.Init()

	// Register Handler
	proto.RegisterHealthyHandler(service.Server(), new(handler.Healthy))
	proto.RegisterBucketHandler(service.Server(), new(handler.Bucket))
	proto.RegisterObjectHandler(service.Server(), new(handler.Object))

	app, _ := filepath.Abs(os.Args[0])

	logger.Info("-------------------------------------------------------------")
	logger.Info("- Micro Service Agent -> Run")
	logger.Info("-------------------------------------------------------------")
	logger.Infof("- version      : %s", BuildVersion)
	logger.Infof("- application  : %s", app)
	logger.Infof("- md5          : %s", md5hex(app))
	logger.Infof("- build        : %s", BuildTime)
	logger.Infof("- commit       : %s", CommitID)
	logger.Info("-------------------------------------------------------------")
	// Run service
	if err := service.Run(); err != nil {
		logger.Error(err)
	}
}

func md5hex(_file string) string {
	h := md5.New()

	f, err := os.Open(_file)
	if err != nil {
		return ""
	}
	defer f.Close()

	io.Copy(h, f)

	return hex.EncodeToString(h.Sum(nil))
}
