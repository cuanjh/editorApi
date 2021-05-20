package initNats

import (
	"editorApi/config"
	"log"

	"github.com/nats-io/nats.go"
)

var NatsConn *nats.EncodedConn

func InitNats() *nats.EncodedConn {
	natsConfig := config.GinVueAdminconfig.NatsConfig
	opts := []nats.Option{nats.Name("TalkMate")}
	var err error
	nc, err := nats.Connect(
		natsConfig.Hosts,
		opts...,
	)

	if err != nil {
		log.Fatalf("nats连接错误：%s", err)
		return nil
	}
	NatsConn, err = nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		log.Fatalf("nats连接错误：%s", err)
		return nil
	}

	return NatsConn
}
