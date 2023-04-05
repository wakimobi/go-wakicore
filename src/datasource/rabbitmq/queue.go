package rabbitmq

import (
	"github.com/idprm/go-pass-tsel/src/config"
	"github.com/wiliehidayat87/rmqp"
)

func InitQueue(cfg *config.Secret) rmqp.AMQP {
	var rb rmqp.AMQP

	rb.SetAmqpURL(cfg.Rmq.Host, cfg.Rmq.Port, cfg.Rmq.User, cfg.Rmq.Pass)

	rb.SetUpConnectionAmqp()
	return rb
}
