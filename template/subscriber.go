package template

var (
	SubscriberSRV = `package subscriber

import (
	"context"

	log "github.com/micro/go-micro/v2/logger"

	"{{.Dir}}/proto/{{.Alias}}"
)

type {{title .Alias}} struct{}

func (e *{{title .Alias}}) Handle(ctx context.Context, msg *{{dehyphen .Alias}}.Message) error {
	log.Info("Handler Received message: ", msg.Say)
	return nil
}

func Handler(ctx context.Context, msg *{{dehyphen .Alias}}.Message) error {
	log.Info("Function Received message: ", msg.Say)
	return nil
}
`
)
