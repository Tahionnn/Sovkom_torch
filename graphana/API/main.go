package main

import (
	_ "HacatonSovKomBank/docs"
	"HacatonSovKomBank/handlers"
	"HacatonSovKomBank/middleware"
	"HacatonSovKomBank/pkg/postgreSQL"
	"github.com/coalaura/mistral"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"net/http"
	"time"
)

// @title           HalvaStats
// @version         1.0
// @description     API для WEB приложения анализа трат
// @host            localhost:8080
// @BasePath		/
func main() {
	router := gin.Default()                                                // Создаем роутер
	client := mistral.NewMistralClient("6c3t8lc5ske6tIraZDakla91bZZMZ6Hf") // Клиент AI
	db := postgreSQL.ConnectDB()                                           //Подключаем DB

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if err := db.AutoMigrate( // Подключаем Автомиграции БД
		&postgreSQL.User{},
		&postgreSQL.Receipt{},
		&postgreSQL.ReceiptCategories{},
	); err != nil {
		log.Fatalf("Error AutoMigrate: %s", err)
	}

	router.Use(middleware.LoggerMiddleware()) // Подключаем Middleware

	handler := handlers.NewHandler(db, client) // Подключаем Handlers
	handler.Register(router)                   // Регистрация хэндлера

	// Запускаем роутер
	start(router)

	//client := mistral.NewMistralClient("6c3t8lc5ske6tIraZDakla91bZZMZ6Hf") // Клиент AI
	//ans, _ := AI.FindCategories(client, "{\n  \"gt_parse\": {\n    \"shop\": \"Монетка\",\n    \"date\": \"25.04.2025\",\n    \"time\": [\n      \"22:04:00\"\n    ],\n    \"items\": [\n      {\n        \"name\": \"Фильтр м\",\n        \"measurement\": \"л\",\n        \"count\": 3,\n        \"price\": 87.82,\n        \"overall\": 263.46\n      },\n      {\n        \"name\": \"Электрич\",\n        \"measurement\": \"шт\",\n        \"count\": 1,\n        \"price\": 456.4,\n        \"overall\": 456.4\n      },\n      {\n        \"name\": \"Принтер\",\n        \"measurement\": \"шт\",\n        \"count\": 3,\n        \"price\": 341.27,\n        \"overall\": 1023.81\n      },\n      {\n        \"name\": \"Чай\",\n        \"measurement\": \"шт\",\n        \"count\": 3,\n        \"price\": 253.68,\n        \"overall\": 761.04\n      },\n      {\n        \"name\": \"Хлеб\",\n        \"measurement\": \"шт\",\n        \"count\": 2,\n        \"price\": 440.71,\n        \"overall\": 881.42\n      }\n    ],\n    \"overall\": \"3386.13\"\n  }\n}")
	//fmt.Println(ans)
}

func start(router *gin.Engine) {
	server := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatalln(server.ListenAndServe())
}
