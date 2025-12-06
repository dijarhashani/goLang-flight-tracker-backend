package controllers

import (
	"airbook-backend/database"
	"airbook-backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func BookFlight(c *gin.Context) {
	userID := c.GetUint("user_id")

	var body struct {
		PassangerName string `json:"passenger_name"`
		FlightID      string `json:"flight_id"`
		FlightName    string `json:"flight_name"`
		Departure     string `json:"flight_departure"`
		Arrival       string `json:"flight_arrival"`
		SeatNumber    string `json:"seat_number"`
		Date          string `json:"date"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON body"})
		return
	}

	booking := models.Booking{
		UserID:        userID,
		PassangerName: body.PassangerName,
		FlightID:      body.FlightID,
		FlightName:    body.FlightName,
		Departure:     body.Departure,
		Arrival:       body.Arrival,
		SeatNumber:    body.SeatNumber,
		Date:          body.Date,
	}

	database.DB.Create(&booking)

	c.JSON(http.StatusOK, gin.H{"message": "Flight booked"})
}

func GetMyBookings(c *gin.Context) {
	userID := c.GetUint("user_id")

	var bookings []models.Booking
	database.DB.Where("user_id = ?", userID).Find(&bookings)

	c.JSON(http.StatusOK, bookings)
}

func DeleteBooking(c *gin.Context) {
	bookingID := c.Param("booking_id")
	userID := c.GetUint("user_id")
	var booking models.Booking

	if err := database.DB.First(&booking, "flight_id = ? AND user_id = ?", bookingID, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
		return
	}
	database.DB.Delete(&booking)
	c.JSON(http.StatusOK, gin.H{"message": "Booking deleted"})
}
