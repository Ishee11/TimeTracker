package handlers

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
	"users/models"
)

// обработчик для создания пользователя
func CreateUserHandler(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			Username string `json:"username" binding:"required"`
			Email    string `json:"email" binding:"required,email"`
			Password string `json:"password" binding:"required,min=8"`
		}

		// Парсим и валидируем входные данные
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Неправильные входные данные"})
			return
		}

		// Хэшируем пароль
		passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			fmt.Println("Ошибка хеширования пароля:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при хэшировании пароля"})
			return
		}

		// Создаем пользователя в базе данных
		user := models.User{
			Username:     input.Username,
			Email:        input.Email,
			PasswordHash: string(passwordHash),
			CreatedAt:    time.Now(),
		}

		// SQL-запрос для вставки нового пользователя
		query := `INSERT INTO users (username, email, password_hash, created_at) 
          VALUES ($1, $2, $3, $4) RETURNING id`
		err = db.QueryRow(context.Background(), query, user.Username, user.Email, user.PasswordHash, user.CreatedAt).Scan(&user.ID)
		if err != nil {
			fmt.Println("Ошибка при создании пользователя:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании пользователя"})
			return
		}

	}
}
