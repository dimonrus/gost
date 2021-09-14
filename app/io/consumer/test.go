package consumer

import (
	"github.com/dimonrus/gorabbit"
	"github.com/streadway/amqp"
	"gost/app/base"
)

func init() {
	base.App.GetRabbit().GetRegistry()["test"] = &gorabbit.Consumer{Queue: "gost.test", Server: "local", Count: 1, Callback: func(d amqp.Delivery) {
		// show the message
		base.App.GetLogger().Infof("%s", d.Body)
	}}
}
