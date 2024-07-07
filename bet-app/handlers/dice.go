package handlers

import (
	"bet-app/config"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func RollDice(c *gin.Context) {
	var input struct {
		BetAmount  float64 `json:"bet_amount" binding:"required,gt=0"`
		BetType    string  `json:"bet_type" binding:"required"`
		Prediction int     `json:"prediction,omitempty"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var balance float64
	err := config.DB.QueryRow("SELECT balance FROM users WHERE id = ?", userID).Scan(&balance)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user balance"})
		return
	}

	if input.BetAmount > balance {
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient balance to place bet"})
		return
	}

	newBalance := balance - input.BetAmount

	rand.Seed(time.Now().UnixNano())
	result := rand.Intn(6) + 1

	_, err = config.DB.Exec("UPDATE users SET balance = ? WHERE id = ?", newBalance, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update balance"})
		return
	}

	switch input.BetType {
	case "number":
		if input.Prediction < 1 || input.Prediction > 6 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Prediction must be between 1 and 6 for number bet"})
			return
		}

		if result == input.Prediction {
			newBalance += input.BetAmount * 6
			_, err = config.DB.Exec("UPDATE users SET balance = ? WHERE id = ?", newBalance, userID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update balance"})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"message":     "Congratulations! You guessed the correct number.",
				"result":      result,
				"new_balance": newBalance,
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"message":     "Sorry, you didn't guess the correct number.",
				"result":      result,
				"new_balance": newBalance,
			})
		}

	case "odd":
		if result%2 != 0 {
			newBalance += input.BetAmount * 2
			_, err = config.DB.Exec("UPDATE users SET balance = ? WHERE id = ?", newBalance, userID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update balance"})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"message":     "Congratulations! You guessed correctly that the number is odd.",
				"result":      result,
				"new_balance": newBalance,
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"message":     "Sorry, you guessed incorrectly that the number is odd.",
				"result":      result,
				"new_balance": newBalance,
			})
		}

	case "even":
		if result%2 == 0 {

			newBalance += input.BetAmount * 2
			_, err = config.DB.Exec("UPDATE users SET balance = ? WHERE id = ?", newBalance, userID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update balance"})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"message":     "Congratulations! You guessed correctly that the number is even.",
				"result":      result,
				"new_balance": newBalance,
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"message":     "Sorry, you guessed incorrectly that the number is even.",
				"result":      result,
				"new_balance": newBalance,
			})
		}

	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bet type"})
	}
}
