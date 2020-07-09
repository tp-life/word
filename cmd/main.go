package main

import (
	_ "word/api/swagger-spec/api"
	"word/cmd/api"
	"word/pkg/database/mongo"
	"word/pkg/database/orm"
	"word/pkg/email"
	"word/pkg/redis"
	"word/pkg/server"
)

func main() {
	server.Register(server.NewService("api", api.Tasks, api.Service, mongo.Start, orm.Start, email.StartEmailSender, redis.Start))

	server.Run()
}
