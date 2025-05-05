package routes

import (
	"dashboard-starter/controllers"
	"dashboard-starter/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/list", controllers.ListUsers)
	r.POST("/insert", controllers.CreateUser)
	r.PUT("/update/:id", controllers.UpdateUser)
	r.DELETE("/delete/:id", controllers.DeleteUser)

	auth := r.Group("/admin")

	auth.POST("/login", controllers.Login)

	auth.Use(middleware.AuthMiddlewear())
	{
		auth.GET("/dashboard", func(c *gin.Context) {
			adminID := c.GetUint("admin_id")
			c.JSON(200, gin.H{"message": "Welcome Admin", "id": adminID})
		})
	}

	return r
}
