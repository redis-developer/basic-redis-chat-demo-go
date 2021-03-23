package message

import "github.com/redis-developer/basic-redis-chat-demo-go/rediscli"

type Controller struct {
	r *rediscli.Redis
}

func NewController(r *rediscli.Redis) *Controller {
	return &Controller{
		r: r,
	}
}
