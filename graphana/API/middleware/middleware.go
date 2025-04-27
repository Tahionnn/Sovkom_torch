package middleware

import (
	"HacatonSovKomBank/pkg/jwtutils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("Новый запрос:", c.Request.Method, c.Request.URL.Path)
		c.Next() // Передаем управление следующему обработчику
	}
}

// AuthMiddleware проверяет наличие и валидность JWT токена
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Отсутствует заголовок Authorization"})
			return
		}

		// Обычно формат такой: "Bearer <токен>", нужно отделить токен
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Неверный формат токена"})
			return
		}

		tokenString := parts[1]

		// Проверяем токен через твою функцию из jwtutils
		username, err := jwtutils.ParseJWT(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Неверный токен"})
			return
		}

		// Сохраняем имя пользователя в контекст, чтобы потом можно было использовать в хендлерах
		c.Set("username", username)
		fmt.Println(username)

		c.Next()
	}
}
