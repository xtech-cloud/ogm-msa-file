package main

import (
	"context"
	"fmt"
	"omo-msa-file/config"
	"time"

	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/metadata"
	_ "github.com/micro/go-plugins/registry/consul/v2"
	_ "github.com/micro/go-plugins/registry/etcdv3/v2"
	proto "github.com/xtech-cloud/omo-msp-file/proto/file"
	pn "github.com/xtech-cloud/omo-msp-notification/proto/notification"
)

type Notification struct {
}

func (this *Notification) Handle(_ctx context.Context, _message *pn.SimpleMessage) error {
	md, ok := metadata.FromContext(_ctx)
	if ok {
		fmt.Println(fmt.Sprintf("[omo.msa.file.notification] Received message %+v with metadata %+v", _message, md))
	} else {
		fmt.Println(fmt.Sprintf("[omo.msa.file.notification] Received message %+v without metadata", _message))
	}
	return nil
}

func main() {
	config.Setup()
	service := micro.NewService(
		micro.Name("omo.msa.file.tester"),
	)
	service.Init()

	micro.RegisterSubscriber("omo.msa.file.notification", service.Server(), new(Notification))

	cli := service.Client()
	cli.Init(
		client.Retries(3),
		client.RequestTimeout(time.Second*1),
		client.Retry(func(_ctx context.Context, _req client.Request, _retryCount int, _err error) (bool, error) {
			if nil != _err {
				fmt.Println(fmt.Sprintf("%v | [ERR] retry %d, reason is %v\n\r", time.Now().String(), _retryCount, _err))
				return true, nil
			}
			return false, nil
		}),
	)

	bucket:= proto.NewBucketService("omo.msa.file", cli)
	object:= proto.NewObjectService("omo.msa.file", cli)

	go test(bucket, object)
	service.Run()
}

func test(_bucket proto.BucketService, _object proto.ObjectService) {
	for range time.Tick(4 * time.Second) {
		fmt.Println("----------------------------------------------------------")
/*

		//查询Profile
		{
			fmt.Println("> Query")
			// Make request
			rsp, err := _profile.Query(context.Background(), &proto.QueryProfileRequest{
				Strategy:    proto.Strategy_STRATEGY_JWT,
				AccessToken: accessToken,
			})
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(rsp)
			}
		}
*/
	}
}
