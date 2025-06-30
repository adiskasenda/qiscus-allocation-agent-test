package main

import (
	"log"
	"net/http"

	_ "github.com/lib/pq"

	"qiscus-test/config"
	controller "qiscus-test/controllers"

	"qiscus-test/repository"
	"qiscus-test/router"
	"qiscus-test/services"

	"github.com/redis/go-redis/v9"
)

func main() {
	config.Init()

	rdb := redis.NewClient((&redis.Options{
		Addr: "localhost:6379",
	}))

	redis := repository.Redis(rdb)
	apiService := services.NewAPIService(config.BaseUrl, config.SecretKey, config.AppCode)
	webhookService := services.AgentAllocationService(redis, apiService)
	webhookController := controller.WebHookController(webhookService)

	r := router.NewRouter(webhookController)
	go webhookService.ProcessQueue()

	log.Println("🚀 Server running at :8080")
	http.ListenAndServe(":8080", r)

}
