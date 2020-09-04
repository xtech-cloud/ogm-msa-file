package publisher

import (
	"context"
	"omo-msa-file/config"

	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/logger"
	proto "github.com/xtech-cloud/omo-msp-notification/proto/notification"
)

var (
	DefaultPublisher micro.Event
	filter           map[string]bool
)

func init() {
	filter = make(map[string]bool)
}

func Publish(_ctx context.Context, _action string, _head string, _body string) {

	if _, ok := filter[_action]; !ok {
		found := false
		for _, action := range config.Schema.Publisher {
			if action == _action {
				found = true
				break
			}
		}
		filter[_action] = !found
	}

	if filter[_action] {
		return
	}

	err := DefaultPublisher.Publish(_ctx, &proto.SimpleMessage{
		Action: _action,
		Head:   _head,
		Body:   _body,
	})
	if nil != err {
		logger.Error(err)
	}
}
