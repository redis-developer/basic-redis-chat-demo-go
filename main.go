package main

import (
	"fmt"
	"github.com/redis-developer/basic-redis-chat-demo-go/config"
	"github.com/redis-developer/basic-redis-chat-demo-go/message"
	"github.com/redis-developer/basic-redis-chat-demo-go/rediscli"
	"github.com/redis-developer/basic-redis-chat-demo-go/websocket"
	"log"
	"net/http"
)

func main() {

	cnf := config.NewConfig()

	log.Println(fmt.Sprintf("%+v", cnf))
	redisCli := rediscli.NewRedis(cnf.RedisAddress, cnf.RedisPassword)
	messageController := message.NewController(redisCli)
	http.Handle("/ws", websocket.Handler(redisCli, messageController))
	http.HandleFunc("/links", func(writer http.ResponseWriter, request *http.Request) {
		_,_ = writer.Write([]byte(`{"github":"https://github.com/redis-developer/basic-redis-chat-demo-go"}`))
	})
	http.Handle("/", http.FileServer(http.Dir(cnf.ClientLocation)))
	log.Fatal(http.ListenAndServe(cnf.ServerAddress, nil))
}

