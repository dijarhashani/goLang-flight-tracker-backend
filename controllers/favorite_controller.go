package controllers

import (
	"airbook-backend/database"
	"airbook-backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func MakeFavorite(c *gin.Context) {
	userID := c.GetUint("user_id")

	var body struct {
		AirportIATA string `json:"airport_iata"`
		AirportName string `json:"airport_name"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON body"})
		return
	}

	var existing models.Favorites
	result := database.DB.Where("user_id = ? AND airport_iata = ?", userID, body.AirportIATA).Take(&existing)

	if result.RowsAffected > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Airport is already in favorites."})
		return
	}

	favorites := models.Favorites{
		UserID:      userID,
		AirportIATA: body.AirportIATA,
		AirportName: body.AirportName,
	}

	database.DB.Create(&favorites)

	c.JSON(http.StatusOK, gin.H{"message": "Airport is added to favorites."})
}

func GetMyFavorites(c *gin.Context) {
	userID := c.GetUint("user_id")

	var favorites []models.Favorites
	database.DB.Where("user_id = ?", userID).Find(&favorites)

	type FavoriteResponse struct {
		AirportIATA string `json:"airportIATA"`
		AirportName string `json:"airportName"`
	}

	var result []FavoriteResponse
	for _, fav := range favorites {
		result = append(result, FavoriteResponse{
			AirportIATA: fav.AirportIATA,
			AirportName: fav.AirportName,
		})
	}

	c.JSON(http.StatusOK, result)
}

func RemoveFavorite(c *gin.Context) {
	userID := c.GetUint("user_id")
	airportIATA := c.Param("airport_iata")

	result := database.DB.Where("user_id = ? AND airport_iata = ?", userID, airportIATA).Delete(&models.Favorites{})

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Favorite not found."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Airport removed from favorites."})

}
