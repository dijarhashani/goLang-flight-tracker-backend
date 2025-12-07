package routes

import (
	"airbook-backend/controllers"
	"airbook-backend/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {

	auth := r.Group("/auth")
	{
		auth.POST("/register", controllers.Register)
		auth.POST("/login", controllers.Login)
		auth.GET("/me", middleware.Auth(), controllers.GetMe)
		auth.POST("/logout", middleware.Auth(), controllers.Logout)

	}

	book := r.Group("/booking").Use(middleware.Auth())
	{
		book.POST("/create", controllers.BookFlight)
		book.GET("/my-bookings", controllers.GetMyBookings)
		book.DELETE("/delete/:booking_id", controllers.DeleteBooking)
	}

	fav := r.Group("/favorites").Use(middleware.Auth())
	{
		fav.POST("/add", controllers.MakeFavorite)
		fav.GET("/my-favorites", controllers.GetMyFavorites)
		fav.DELETE("/remove/:airport_iata", controllers.RemoveFavorite)
	}

	r.GET("/flights", controllers.GetFlights)
	r.GET("/plane-image/:hex", controllers.GetPlaneImage)

	r.GET("/airport-info", controllers.GetAeroFlights)

	r.GET("/airline/name", controllers.GetAirlineName)

	r.GET("/search-flights", controllers.SearchFlights)

	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path

		if path == "/" {
			c.File("./docs/index.html")
			return
		}
		if path == "/docs" || path == "/docs/" {
			c.File("./docs/index.html")
			return
		}
		if len(path) > 6 && path[:6] == "/docs/" {
			filePath := "." + path
			c.File(filePath)
			return
		}
	})

}
