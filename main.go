package main

import (
	"database/sql"
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

	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:5432/qiscus?sslmode=disable")
	if err != nil {
		log.Fatal("❌ Gagal koneksi DB:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("❌ DB tidak tersedia:", err)
	}

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
