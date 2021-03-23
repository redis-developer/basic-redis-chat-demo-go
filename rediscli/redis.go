package rediscli

import (
	"github.com/go-redis/redis"
	"log"
	"sync"
)

type Redis struct {
	client             *redis.Client
	channelsPubSub     map[string]*ChannelPubSub
	channelsPubSubSync *sync.RWMutex
}

type ChannelPubSub struct {
	close  chan struct{}
	closed chan struct{}
	pubSub *redis.PubSub
}

func (channel *ChannelPubSub) Channel() <-chan *redis.Message {
	return channel.pubSub.Channel()
}

func (channel *ChannelPubSub) Close() <-chan struct{} {
	return channel.close
}

func (channel *ChannelPubSub) Closed() <-chan struct{} {
	return channel.closed
}

func NewRedis(addr, passwd string) *Redis {

	log.Println("Initialized redis client", addr,passwd)

	opt := &redis.Options{
		Addr: addr,
	}

	if passwd != "" {
		opt.Password = passwd
	}

	c := redis.NewClient(opt)
	if err := c.Ping().Err(); err != nil {
		panic(err)
	}

	r := &Redis{
		client:             c,
		channelsPubSub:     make(map[string]*ChannelPubSub, 0),
		channelsPubSubSync: &sync.RWMutex{},
	}
	return r
}
