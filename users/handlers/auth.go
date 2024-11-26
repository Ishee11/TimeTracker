package handlers

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"users/models"
)

type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// поиск пользователя в базе данных
func findUserByUsername(db *pgxpool.Pool, username string) (*models.User, error) {
	var user models.User
	query := `SELECT id, username, email, password_hash, created_at FROM users WHERE username = $1`
	err := db.QueryRow(context.Background(), query, username).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// проверка пароля
func checkPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

// обработчик для создания пользователя
func LoginHandler(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input LoginInput
		if err := c.ShouldBindJSON(&input); err != nil {
			fmt.Println("Ошибка валидации:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка валидации"})
			return
		}

		// Проверяем, существует ли пользователь и совпадают ли пароли
		user, err := findUserByUsername(db, input.Username)
		if err != nil || user == nil {
			fmt.Println("Неправильные учетные данные:", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Неправильные учетные данные"})
			return
		}

		// Проверяем хеш пароля
		err = checkPasswordHash(input.Password, user.PasswordHash)
		if err != nil {
			fmt.Println("Неправильные учетные данные:", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Неправильные учетные данные"})
			return
		}

		// Успешная авторизация
		c.JSON(http.StatusOK, gin.H{"message": "Успешная авторизация"})
	}
}
