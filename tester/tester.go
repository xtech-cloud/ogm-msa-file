package main

import (
	"context"
	"fmt"
	"omo-msa-account/config"
	"time"

	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/metadata"
	_ "github.com/micro/go-plugins/registry/consul/v2"
	_ "github.com/micro/go-plugins/registry/etcdv3/v2"
	proto "github.com/xtech-cloud/omo-msp-account/proto/account"
	pn "github.com/xtech-cloud/omo-msp-notification/proto/notification"
)

type Notification struct {
}

func (this *Notification) Handle(_ctx context.Context, _message *pn.SimpleMessage) error {
	md, ok := metadata.FromContext(_ctx)
	if ok {
		fmt.Println(fmt.Sprintf("[omo.msa.account.notification] Received message %+v with metadata %+v", _message, md))
	} else {
		fmt.Println(fmt.Sprintf("[omo.msa.account.notification] Received message %+v without metadata", _message))
	}
	return nil
}

func main() {
	config.Setup()
	service := micro.NewService(
		micro.Name("omo.msa.account.tester"),
	)
	service.Init()

	micro.RegisterSubscriber("omo.msa.account.notification", service.Server(), new(Notification))

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

	auth := proto.NewAuthService("omo.msa.account", cli)
	profile := proto.NewProfileService("omo.msa.account", cli)

	go test(auth, profile)
	service.Run()
}

func test(_auth proto.AuthService, _profile proto.ProfileService) {
	for range time.Tick(4 * time.Second) {
		fmt.Println("----------------------------------------------------------")
		accessToken := ""

		{
			fmt.Println("> Signup")
			// Make request
			rsp, err := _auth.Signup(context.Background(), &proto.SignupRequest{
				Username: "user001",
				Password: "11112222",
			})
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(rsp)
			}
		}

		//错误账号
		{
			fmt.Println("> Signin")
			// Make request
			rsp, err := _auth.Signin(context.Background(), &proto.SigninRequest{
				Strategy: proto.Strategy_STRATEGY_JWT,
				Username: "user",
				Password: "11112222",
			})
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(rsp)
			}
		}

		//错误密码
		{
			fmt.Println("> Signin")
			// Make request
			rsp, err := _auth.Signin(context.Background(), &proto.SigninRequest{
				Strategy: proto.Strategy_STRATEGY_JWT,
				Username: "user001",
				Password: "11223344",
			})
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(rsp)
			}
		}

		//正确账号
		{
			fmt.Println("> Signin")
			// Make request
			rsp, err := _auth.Signin(context.Background(), &proto.SigninRequest{
				Strategy: proto.Strategy_STRATEGY_JWT,
				Username: "user001",
				Password: "11112222",
			})
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(rsp)
				accessToken = rsp.AccessToken
			}
		}

		//重置密码
		{
			fmt.Println("> ResetPasswd")
			// Make request
			rsp, err := _auth.ResetPasswd(context.Background(), &proto.ResetPasswdRequest{
				Strategy:    proto.Strategy_STRATEGY_JWT,
				AccessToken: accessToken,
				Password:    "abcdefg",
			})
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(rsp)
			}
		}

		//重置密码
		{
			fmt.Println("> ResetPasswd")
			// Make request
			rsp, err := _auth.ResetPasswd(context.Background(), &proto.ResetPasswdRequest{
				Strategy:    proto.Strategy_STRATEGY_JWT,
				AccessToken: accessToken,
				Password:    "11112222",
			})
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(rsp)
			}
		}

		//更新Profile
		{
			fmt.Println("> Update")
			// Make request
			rsp, err := _profile.Update(context.Background(), &proto.UpdateProfileRequest{
				Strategy:    proto.Strategy_STRATEGY_JWT,
				AccessToken: accessToken,
				Profile:     "myprofile:" + time.Now().String(),
			})
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(rsp)
			}
		}

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
	}
}
